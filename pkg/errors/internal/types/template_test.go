package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gongt/go/internal/myenv"
	. "github.com/gongt/go/pkg/errors/internal/types"
)

var template = CreateTemplate("hello %s")

func TestErrorTemplate_New(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	var eTmpl error = template
	err := template.New("world")

	assert.Equal(t, "hello world", err.Error(), "模板功能异常")
	assert.Equal(t, true, errors.Is(err, eTmpl), "Is()功能异常")
}
