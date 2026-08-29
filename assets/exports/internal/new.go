package internal

type SomeStruct struct{}

// @exported
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{}
}
