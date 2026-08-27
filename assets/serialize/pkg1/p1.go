package pkg1

import "github.com/gongt/go/assets/serialize/pkg1/internal"

type PublicType1 struct {
	Field1       internal.InternalType1
	Field2       string
	privateField uint64
}
