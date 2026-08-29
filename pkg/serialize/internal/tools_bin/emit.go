package tools_bin

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/i18n/type_name"
	"github.com/gongt/go/pkg/serialize/kinds"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/strings/strtools"
)

// namedField 是结构体中的一个具名字段（匿名嵌入字段使用其类型名作为字段名）
type namedField struct {
	Name string
	Type ast.Expr
}

// structInfo 是待生成 Marshal/Unmarshal 的目标结构体
type structInfo struct {
	name   string
	fields []namedField
}

func collectStructFields(list []*ast.Field) []namedField {
	var out []namedField
	for _, f := range list {
		if len(f.Names) == 0 {
			name := embeddedFieldName(f.Type)
			if name == "" {
				continue
			}
			out = append(out, namedField{Name: name, Type: f.Type})
			continue
		}
		for _, n := range f.Names {
			if n.Name == "_" {
				continue
			}
			out = append(out, namedField{Name: n.Name, Type: f.Type})
		}
	}
	return out
}

func embeddedFieldName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return e.Sel.Name
	case *ast.StarExpr:
		return embeddedFieldName(e.X)
	default:
		return ""
	}
}

// genCtx 保存一次代码生成过程中共用的状态：输出缓冲区、类型检查信息、导入别名与临时变量计数器
type genCtx struct {
	buf  sourcecode.GoFileBuffer
	file sourcecode.FileInfo
	info *types.Info

	helpersPkg   string
	serializePkg string
}

func Ctx(file sourcecode.FileInfo, buff sourcecode.GoFileBuffer) *genCtx {
	if buff == nil {
		buff = sourcecode.NewGoFileBuffer()
	}
	return &genCtx{
		buf:  buff,
		file: file,
		info: file.Container().TypesInfo(),

		helpersPkg:   buff.AddImport("github.com/gongt/go/pkg/serialize/helpers"),
		serializePkg: buff.AddImport("github.com/gongt/go/pkg/serialize"),
	}
}

func (g *genCtx) reflectPkg() string {
	return g.buf.AddImport("reflect")
}

func (g *genCtx) kindsPkg() string {
	return g.buf.AddImport("github.com/gongt/go/pkg/serialize/kinds")
}

func (g *genCtx) errorsPkg() string {
	return g.buf.AddImport("github.com/gongt/go/pkg/errors")
}

func (g *genCtx) use(ch byte) string {
	return fmt.Sprintf("%s.Use%c(ctx)", g.helpersPkg, ch)
}

func (g *genCtx) emitMarshalFunc(st *structInfo) {
	fmt.Fprintf(g.buf, "func (t *%s) Marshal(ctx %s.SerializeContext) ([]byte, error) {\n", st.name, g.serializePkg)
	fmt.Fprintf(g.buf, "var err error\n")
	fmt.Fprintf(g.buf, "if ctx, err = %s.EnsureMarshalContext(ctx); err != nil {\n", g.helpersPkg)
	fmt.Fprintf(g.buf, "  return nil, err\n")
	fmt.Fprintf(g.buf, "}\n")
	for _, f := range st.fields {
		g.emitMarshalValue(f.Type, "t."+f.Name, "字段"+f.Name)
	}
	g.buf.WriteString(fmt.Sprintf("\treturn %s.Bytes(), nil\n}\n\n", g.use('S')))
}

func (g *genCtx) emitUnmarshalFunc(st *structInfo) {
	fmt.Fprintf(g.buf, "func (t *%s) Unmarshal(data []byte, ctx %s.DeserializeContext) error {\n", st.name, g.serializePkg)
	fmt.Fprintf(g.buf, "\tvar err error\n")
	fmt.Fprintf(g.buf, "\tctx, err = %s.EnsureUnmarshalContext(data, ctx)\n", g.helpersPkg)
	fmt.Fprintf(g.buf, "\tif err != nil { return err }\n")
	for _, f := range st.fields {
		varName := fmt.Sprintf("t.%s", f.Name)
		g.emitUnmarshalValue(varName, f.Type, "字段"+f.Name)
	}
	g.buf.WriteString("\treturn nil\n}\n\n")
}

// emitMarshalValue 生成序列化一个值（结构体字段、或slice/map的元素）的代码，accessor是该值在生成代码中的表达式
func (g *genCtx) emitMarshalValue(expr ast.Expr, accessor, label string) {
	handleChanFunc := func(handler string) {
		typeStr := typeText(g.buf, g.file, expr)
		fmt.Fprintf(g.buf, "if value, err := %s.%s(ctx, %s, %s.TypeFor[%s]()); err != nil {\n", g.helpersPkg, handler, accessor, g.reflectPkg(), typeStr)
		fmt.Fprintf(g.buf, "\treturn nil, %s.Extend(err, \"序列化%s失败\")\n", g.errorsPkg(), label)
		fmt.Fprintf(g.buf, "} else {\n")
		fmt.Fprintf(g.buf, "\t%s.WriteAnyId(%s.TypeIdChannel, value)\n", g.use('S'), g.kindsPkg())
		fmt.Fprintf(g.buf, "}\n")
	}

	fmt.Fprintf(g.buf, "// ==== %s ====\n{\n", label)
	t := g.info.TypeOf(expr)
	switch classifyType(t) {
	case kindScalar:
		g.emitScalarMarshal(t.Underlying().(*types.Basic), accessor)
	case kindSlice:
		arr := expr.(*ast.ArrayType)
		elemType := g.info.TypeOf(arr.Elt)
		if classifyType(elemType) == kindScalar {
			// 子元素是标量类型
			g.emitScalarSliceMarshal(elemType.Underlying().(*types.Basic), accessor)
		} else {
			fmt.Fprintf(g.buf, "%s.WriteAnyHeader(%s.TypeIdArray, len(%s))\n", g.use('S'), g.kindsPkg(), accessor)
			fmt.Fprintf(g.buf, "\tfor _, item := range %s {\n", accessor)
			g.emitMarshalValue(arr.Elt, "item", label+"的元素")
			g.buf.WriteString("\t}\n")
		}
	case kindMap:
		m := expr.(*ast.MapType)
		fmt.Fprintf(g.buf, "\t%s.WriteAnyHeader(%s.TypeIdMap, len(%s))\n", g.use('S'), g.kindsPkg(), accessor)
		fmt.Fprintf(g.buf, "\tfor k, v := range %s {\n", accessor)
		g.emitMarshalValue(m.Key, "k", label+"的键")
		g.emitMarshalValue(m.Value, "v", label+"的值")
		g.buf.WriteString("\t}\n")
	case kindChan:
		handleChanFunc("SendChannel")
	case kindFunc:
		handleChanFunc("SendFunc")
	default:
		addr := "&" + accessor
		if isInterfaceType(t) {
			addr = accessor
		}
		fmt.Fprintf(g.buf, "if err := %s.HelperMarshal(%s, ctx); err != nil {\n", g.helpersPkg, addr)
		fmt.Fprintf(g.buf, "  return nil, %s.Extend(err, \"序列化%s失败\")\n", g.errorsPkg(), label)
		g.buf.WriteString("}\n")
	}
	g.buf.WriteString("}\n")
}

func (g *genCtx) emitScalarMarshal(kind *types.Basic, accessor string) {
	goType, ok := basicWriteInfo(kind.Kind())
	if !ok {
		panic(errors.NewAnonymous("无法处理的标量类型: %s", type_name.TranslateBasicType(kind)))
	}

	suffix := strtools.UcFirst(goType)
	fmt.Fprintf(g.buf, "\t%s.Write%s(%s)\n", g.use('S'), suffix, accessor)
}

func (g *genCtx) emitScalarSliceMarshal(elemKind *types.Basic, accessor string) {
	_, ok := basicWriteInfo(elemKind.Kind())
	if !ok {
		panic(errors.NewAnonymous("无法处理的标量类型: %s", type_name.TranslateBasicType(elemKind)))
	}

	if elemKind.Kind() == types.Byte {
		fmt.Fprintf(g.buf, "\t%s.Write(%s)\n", g.use('S'), accessor)
		return
	}
	vt := kinds.IdOfType(elemKind.Kind())
	fmt.Fprintf(g.buf, "\t%s.WriteArray(%s.%s, %s)\n", g.use('S'), g.kindsPkg(), vt, accessor)
}

// emitUnmarshalValue 生成反序列化一个值的代码，返回一个已声明并持有该值的临时变量名，类型与expr完全一致
func (g *genCtx) emitUnmarshalValue(accessor string, expr ast.Expr, label string) {
	fmt.Fprintf(g.buf, "// ==== %s ====\n{\n", label)
	t := g.info.TypeOf(expr)
	switch classifyType(t) {
	case kindScalar:
		g.emitScalarUnmarshal(accessor, t.Underlying().(*types.Basic), expr, label)
	case kindSlice:
		arr := expr.(*ast.ArrayType)
		elemType := g.info.TypeOf(arr.Elt)
		if classifyType(elemType) == kindScalar {
			g.emitScalarSliceUnmarshal(elemType.Underlying().(*types.Basic), accessor)
		} else {
			g.emitSliceUnmarshal(accessor, arr, label)
		}
	case kindMap:
		g.emitMapUnmarshal(accessor, expr.(*ast.MapType), label)
	case kindChan:
		g.emitChanFuncUnmarshal(accessor, expr, label, true)
	case kindFunc:
		g.emitChanFuncUnmarshal(accessor, expr, label, false)
	default:
		g.emitFallbackUnmarshal(accessor, expr, label)
	}
	g.buf.WriteString("}\n")
}

func (g *genCtx) emitScalarUnmarshal(accessor string, kind *types.Basic, expr ast.Expr, label string) {
	goType, ok := basicWriteInfo(kind.Kind())
	if !ok {
		panic(errors.NewAnonymous("无法处理的标量类型: %s", type_name.TranslateBasicType(kind)))
	}

	suffix := strtools.UcFirst(goType)
	fmt.Fprintf(g.buf, "\t%s, err = %s.Read%s()\n", accessor, g.use('D'), suffix)
	fmt.Fprintf(g.buf, "\tif err != nil { return %s.Extend(err, \"读取%s失败\") }\n", g.errorsPkg(), label)
}

func (g *genCtx) emitSliceUnmarshal(accessor string, arr *ast.ArrayType, label string) {
	elemType := typeText(g.buf, g.file, arr.Elt)

	fmt.Fprintf(g.buf, "\tn, err := %s.ReadAnyHeader(%s.TypeIdArray)\n", g.use('D'), g.kindsPkg())
	fmt.Fprintf(g.buf, "\tif err != nil { return %s.Extend(err, \"读取%s长度失败\") }\n", g.errorsPkg(), label)
	fmt.Fprintf(g.buf, "\t%s = make([]%s, n)\n", accessor, elemType)

	fmt.Fprintf(g.buf, "\tfor i := range n {\n")
	vName := fmt.Sprintf("%s[i]", accessor)
	g.emitUnmarshalValue(vName, arr.Elt, label+"的元素")
	g.buf.WriteString("\t}\n")
}

func (g *genCtx) emitScalarSliceUnmarshal(elemKind *types.Basic, accessor string) {
	_, ok := basicWriteInfo(elemKind.Kind())
	if !ok {
		panic(errors.NewAnonymous("无法处理的标量类型: %s", type_name.TranslateBasicType(elemKind)))
	}

	if elemKind.Kind() == types.Byte {
		fmt.Fprintf(g.buf, "\t%s, err = %s.Read()\n", accessor, g.use('D'))
		fmt.Fprintf(g.buf, "\tif err != nil { return %s.Extend(err, \"读取%s失败\") }\n", g.errorsPkg(), accessor)
		return
	}
	vt := kinds.IdOfType(elemKind.Kind())
	fmt.Fprintf(g.buf, "\t%s, err = %s.ReadArray[%s](%s.%s)\n", accessor, g.use('D'), elemKind, g.kindsPkg(), vt)
	fmt.Fprintf(g.buf, "\tif err != nil { return %s.Extend(err, \"读取%s失败\") }\n", g.errorsPkg(), accessor)
}

func (g *genCtx) emitMapUnmarshal(accessor string, m *ast.MapType, label string) {
	keyType := typeText(g.buf, g.file, m.Key)
	valType := typeText(g.buf, g.file, m.Value)

	fmt.Fprintf(g.buf, "\tn, err := %s.ReadAnyHeader(%s.TypeIdMap)\n", g.use('D'), g.kindsPkg())
	fmt.Fprintf(g.buf, "\tif err != nil { return %s.Extend(err, \"读取%s头部失败\") }\n", g.errorsPkg(), label)
	fmt.Fprintf(g.buf, "\t%s = make(map[%s]%s, n)\n", accessor, keyType, valType)

	fmt.Fprintf(g.buf, "\tfor _ = range n {\n")
	fmt.Fprintf(g.buf, "\t\tvar key %s\n", keyType)
	fmt.Fprintf(g.buf, "\t\tvar val %s\n", valType)
	g.emitUnmarshalValue("key", m.Key, label+"的键")
	g.emitUnmarshalValue("val", m.Value, label+"的值")
	fmt.Fprintf(g.buf, "\t\t%s[key] = val\n", accessor)
	// Note: 'i' is the loop variable, but it's not used directly in the map assignment.
	g.buf.WriteString("\t}\n")
}

// emitChanUnmarshal 处理chan和func两类无法自动重建的字段：读取8字节ID并注册其静态类型信息，供调用者事后按ID取回
func (g *genCtx) emitChanFuncUnmarshal(accessor string, expr ast.Expr, label string, isChan bool) {
	typeStr := typeText(g.buf, g.file, expr)
	var id string

	if isChan {
		id = "TypeIdChannel"
	} else {
		id = "TypeIdFunc"
	}

	fmt.Fprintf(g.buf, "\tdata, err := %s.ReadAnyId(%s.%s)\n", g.use('D'), g.kindsPkg(), id)
	fmt.Fprintf(g.buf, "\tif err != nil { return %s.Extend(err, \"读取%sID失败\") }\n", g.errorsPkg(), label)
	var handler string
	if isChan {
		handler = "ReceiveChannel"
	} else {
		handler = "ReceiveFunc"
	}
	fmt.Fprintf(g.buf, "\tanyVal, err := %s.%s(ctx, data, %s.TypeFor[%s]())\n", g.helpersPkg, handler, g.reflectPkg(), typeStr)

	fmt.Fprintf(g.buf, "\tif err != nil { return err }\n")
	fmt.Fprintf(g.buf, "\t%s = anyVal.(%s)\n", accessor, typeStr)
}

func (g *genCtx) emitFallbackUnmarshal(accessor string, expr ast.Expr, label string) {
	t := g.info.TypeOf(expr)
	// typeStr := typeText(g.buf, g.file, expr)

	ref := "&"
	if isInterfaceType(t) {
		ref = ""
	}

	fmt.Fprintf(g.buf, "\tif err := %s.HelperUnmarshal(%s%s, ctx); err != nil { return %s.Extend(err, \"反序列化%s失败\") }\n",
		g.helpersPkg, ref, accessor, g.errorsPkg(), label)
}
