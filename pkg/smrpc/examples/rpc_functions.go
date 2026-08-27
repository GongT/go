package examples

type MsgType struct {
	Msg string
}

type SubType struct {
	privateChan chan MsgType

	field int
}

type MainType struct {
	SubPtr *SubType

	funcPtr func(arg *MainType) error
}

func RPCFunction(arg *MainType) error {
	return arg.funcPtr(arg)
}

func implementFunc(arg *MainType) error {
	// Example implementation
	arg.SubPtr.field = 123

	return nil
}
