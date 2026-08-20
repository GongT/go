package types

type LibType interface {
	*ErrorObjectBase | *ErrorObjectJoined | *ErrorObjectWrapped
}
