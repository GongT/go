package tools_bin

import (
	"fmt"
	"uuid"
)

func (g *genCtx) emitHelpers(st *structInfo) {
	id := uuid.NewV7()

	fmt.Fprintf(g.buf, "var _ %s.Marshaller = (*%s)(nil)\n", g.serializePkg, st.name)

	fmt.Fprintf(g.buf, "// %s\nfunc (*%s) TypeId() [16]byte { return %#v }\n", id, st.name, [16]byte(id))

	//		fmt.Fprintf(g.buf, `func init() {
	//		%s.RegisterType((*%s)(nil).TypeId(), %s.TypeFor[*%s]())
	//	}
	//
	// `, g.serializePkg, st.name, g.reflectPkg(), st.name)
}
