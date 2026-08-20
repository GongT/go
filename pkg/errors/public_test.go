package errors_test

import (
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/errors/internal"
	"github.com/stretchr/testify/assert"
)

var tmpl = errors.NewTemplate("test template: %s")

func TestNewInstance(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	err1 := errors.NewInstance(tmpl, "123")

	err2 := err1.WithDetails("hello", "world").OverrideMessage("new message")

	assert.Equal(t, err1, err2, "WithDetails返回应该是同一个地址")

	assert.Equal(t, "new message", err2.Error(), "OverrideMessage应该修改错误信息")
	assert.Equal(t, "new message", err1.Error())

	assert.Equal(t, tmpl, err2.(internal.UnWrap).Unwrap())
}
