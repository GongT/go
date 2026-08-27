package pkg2

import "github.com/gongt/go/assets/serialize/pkg2/internal"

type PublicType2 struct {
	Field1       internal.InternalType2
	Field2       func(ch <-chan int) error
	privateField uint64
}
