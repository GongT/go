package sharederrors

import "errors"

var ErrBrokenPipe = errors.New("broken pipe")
var ErrEntityTooLarge = errors.New("entity too large")
var ErrDuplicateCall = errors.New("duplicate call")
var ErrDataCorrupted = errors.New("data corrupted")
