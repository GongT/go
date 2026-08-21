package types

import (
	"fmt"
)

type ErrorTemplate struct {
	message string
}

func CreateTemplate(message string) *ErrorTemplate {
	return &ErrorTemplate{message}
}

func (t *ErrorTemplate) New(args ...any) *ErrorObjectBase {
	msg := fmt.Sprintf(t.message, args...)

	err := Create(t, 1)
	err.OverrideMessage(msg)

	return err
}

func (t *ErrorTemplate) Wrap(e error, args ...any) *ErrorObjectBase {
	if e == nil {
		return nil
	}
	msg := fmt.Sprintf(t.message, args...)

	err := Create(t, 1)
	err.OverrideMessage(msg)

	err.AlsoBe(e)
	return err
}

func (t *ErrorTemplate) Error() string {
	return t.message
}
