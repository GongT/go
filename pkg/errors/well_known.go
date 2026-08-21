package errors

type WellKnownError struct {}

var ErrNotImplemented = NewTemplate("Not Implemented")
