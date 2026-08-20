package types_test

import (
	"errors"
	"testing"

	"github.com/gongt/go/internal/myenv"
	. "github.com/gongt/go/pkg/errors/internal/types"
	"github.com/stretchr/testify/assert"
)

const multiShouldBe = `发生多个错误:
  - error object 1
  - error object 2
  - error object 3`

func TestJoin_New(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	err1 := Create(errors.New("error object 1"), 0)
	err2 := errors.New("error object 2")
	err3 := errors.New("error object 3")

	assert.Equal(t, nil, Join(0, []error{nil, nil, nil}, false), "Join(nil...)异常")
	assert.Equal(t, err1, Join(0, []error{nil, err1, nil}, false), "Join(err1)没有原样返回err1")

	assert.NotEqual(t, err1, Join(0, []error{nil, err1, nil}, true), "Join(err1, force_message)应该返回新的对象")

	assert.NotEqual(t, err2, Join(0, []error{nil, err2, nil}, false), "Join(err2)非Err类型应该返回新的对象")

	joined := Join(0, []error{err1, err2, err3}, false)
	assert.Equal(t, multiShouldBe, joined.Error(), "Join后的错误信息与预期不符")
}
