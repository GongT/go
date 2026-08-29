package writer

import (
	"bytes"
	"fmt"
	"go/types"
	"strings"

	"github.com/goforj/godump"
	"github.com/gongt/go/pkg/errors"
)

type TypeResolver = *typeResolver

type typeResolver struct {
	file GoFileBuffer
}

func NewTypeResolver(file GoFileBuffer) TypeResolver {
	return &typeResolver{
		file: file,
	}
}

var ErrNotName = fmt.Errorf("not a type name")
var ErrBasic = fmt.Errorf("symbol in basic")

// var ErrInStl = fmt.Errorf("symbol in standard library")

func (tr TypeResolver) PackageOf(typ types.Type) (*types.Package, error) {
	switch t := typ.(type) {
	case *types.Basic:
		return nil, ErrBasic
	case *types.Pointer:
		return tr.PackageOf(t.Elem())
	case *types.Slice:
		return tr.PackageOf(t.Elem())
	case *types.Array:
		return tr.PackageOf(t.Elem())
	case *types.Chan:
		return tr.PackageOf(t.Elem())
	case *types.Named:
		if t.Obj().Pkg() == nil {
			return nil, ErrBasic // comparable 之类的
		} else {
			return t.Obj().Pkg(), nil
		}
	case *types.TypeParam:
		if t.Obj().Pkg() == nil {
			return nil, ErrBasic // comparable 之类的
		} else {
			return t.Obj().Pkg(), nil
		}
	case *types.Alias:
		if t.Obj().Pkg() == nil {
			return nil, ErrBasic // comparable 之类的
		} else {
			return t.Obj().Pkg(), nil
		}
	case *types.Union, *types.Map, *types.Struct, *types.Signature, *types.Interface:
		return nil, ErrNotName
	default:
		panic(fmt.Sprintf("Unsupported type: %s", godump.DumpStr(typ)))
	}
}

func (tr TypeResolver) qtBasic(typ types.Type) string {
	t := typ.(*types.Basic)
	if t.Kind() == types.Invalid {
		return "__invalid_type__"
	}
	return t.Name()
}
func (tr TypeResolver) qtPointer(typ types.Type) string {
	return "*" + tr.QualifyType(typ.(*types.Pointer).Elem())
}
func (tr TypeResolver) qtSlice(typ types.Type) string {
	return "[]" + tr.QualifyType(typ.(*types.Slice).Elem())
}
func (tr TypeResolver) qtArray(typ types.Type) string {
	t := typ.(*types.Array)
	return fmt.Sprintf("[%d]%s", t.Len(), tr.QualifyType(t.Elem()))
}
func (tr TypeResolver) qtMap(typ types.Type) string {
	t := typ.(*types.Map)
	return fmt.Sprintf("map[%s]%s", tr.QualifyType(t.Key()), tr.QualifyType(t.Elem()))
}
func (tr TypeResolver) qtStruct(typ types.Type) string {
	t := typ.(*types.Struct)
	body := &strings.Builder{}
	body.WriteString("struct{")
	for i := 0; i < t.NumFields(); i++ {
		field := t.Field(i)
		fmt.Fprintf(body, "%s %s", field.Name(), tr.QualifyType(field.Type()))
		if tag := t.Tag(i); tag != "" {
			fmt.Fprintf(body, " `%s`", tag)
		}
		if i < t.NumFields()-1 {
			body.WriteString("; ")
		}
	}
	body.WriteString("}")
	return body.String()
}
func (tr TypeResolver) qtInterface(typ types.Type) string {
	t := typ.(*types.Interface)
	if t.NumMethods() == 0 {
		return "any"
	}
	body := &strings.Builder{}
	body.WriteString("interface{")
	for i := 0; i < t.NumMethods(); i++ {
		method := t.Method(i)
		fmt.Fprintf(body, "%s%s", method.Name(), tr.qtSignature(method.Type()))
		if i < t.NumMethods()-1 {
			body.WriteString("; ")
		}
	}
	body.WriteString("}")
	return body.String()
}
func (tr TypeResolver) qtSignature(typ types.Type) string {
	t := typ.(*types.Signature)
	fc := &bytes.Buffer{}

	if t.TypeParams() != nil {
		// 目前能想到的调用类型都在其他调用的参数中，不可能存在泛型，需要再做
		panic(errors.NewAnonymous("暂不支持泛型参数"))
	}

	fmt.Fprintf(fc, "(")
	params := t.Params()
	for i := 0; i < params.Len(); i++ {
		if i > 0 {
			fmt.Fprintf(fc, ", ")
		}
		typ := params.At(i)
		types.WriteType(fc, typ.Type(), tr.file.Qualifier)
	}
	fmt.Fprintf(fc, ")")
	results := t.Results()
	if results.Len() > 0 {
		fmt.Fprintf(fc, " (")
		for i := 0; i < results.Len(); i++ {
			if i > 0 {
				fmt.Fprintf(fc, ", ")
			}
			typ := results.At(i)
			fmt.Fprintf(fc, "%s", tr.QualifyType(typ.Type()))
		}
		fmt.Fprintf(fc, ")")
	}
	return fc.String()
}
func (tr TypeResolver) qtChan(typ types.Type) string {
	t := typ.(*types.Chan)
	dir := ""
	switch t.Dir() {
	case types.SendRecv:
		dir = "chan "
	case types.SendOnly:
		dir = "chan<- "
	case types.RecvOnly:
		dir = "<-chan "
	}
	return dir + tr.QualifyType(t.Elem())
}
func (tr TypeResolver) qtAlias(typ types.Type) string {
	t := typ.(*types.Alias)
	return tr.QualifyType(t.Underlying())
}
func (tr TypeResolver) qtNamed(typ types.Type) string {
	t := typ.(*types.Named)
	pkgPath := ""
	if t.Obj().Pkg() != nil {
		pkgPath = t.Obj().Pkg().Path()
	}
	return tr.file.QualifyTypeName(pkgPath, t.Obj().Name())
}
func (tr TypeResolver) qtUnion(typ types.Type) string {
	t := typ.(*types.Union)
	terms := &strings.Builder{}
	for i := 0; i < t.Len(); i++ {
		if i > 0 {
			terms.WriteString(" | ")
		}

		term := t.Term(i)
		if term.Tilde() {
			terms.WriteString("~")
		}

		terms.WriteString(tr.QualifyType(term.Type()))
	}
	return terms.String()
}

func (tr TypeResolver) QualifyType(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Basic:
		if t.Kind() == types.Invalid {
			return "__invalid_type__"
		}
		return t.Name()
	default:
		buff := &bytes.Buffer{}
		types.WriteType(buff, typ, tr.file.Qualifier)
		return buff.String()
	}
}
