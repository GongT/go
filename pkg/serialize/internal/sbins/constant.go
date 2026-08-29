package sbins

import (
	"bytes"

	"github.com/gongt/go/pkg/errors"
)

var startFlag [4]byte = [4]byte{0x12, 0x34, 0x56, 0x78}
var endFlag [4]byte = [4]byte{0xFE, 0xDC, 0xBA, 0x98}

// validateBuffer 验证缓冲区的起始和结束标志，并返回去掉标志后的有效数据。
func validateBuffer(buff []byte) ([]byte, error) {
	if len(buff) < len(startFlag)+len(endFlag) {
		return nil, errors.NewAnonymous("缓冲区太短")
	}

	if bytes.Equal(buff[0:4], startFlag[:]) != true {
		return nil, errors.NewAnonymous("起始标志不匹配")
	}
	if bytes.Equal(buff[len(buff)-4:], endFlag[:]) != true {
		return nil, errors.NewAnonymous("结束标志不匹配")
	}

	return buff[len(startFlag) : len(buff)-len(endFlag)], nil
}
