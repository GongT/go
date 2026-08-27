//go:generate go run ../../cmd/serialize/main.go

package serialize

import (
	"github.com/gongt/go/assets/serialize/pkg1"
	"github.com/gongt/go/assets/serialize/pkg2"
)

type privateType struct {
	Field1 pkg1.PublicType1
	Field2 pkg2.PublicType2
}

type SomeType struct {
	Field1 pkg1.PublicType1
	Field2 pkg2.PublicType2
	Field3 privateType
}
