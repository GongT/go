package strtools

import (
	"hash/fnv"
	"io"

	"github.com/gongt/go/pkg/interfaces"
)

func TinyHash[T interfaces.ByteSeq](content T) string {
	h := fnv.New64a()
	h.Write([]byte(content))
	return Base52Encode(h.Sum64())
}

func TinyHashStream(r io.Reader) (string, error) {
	h := fnv.New64a()
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return Base52Encode(h.Sum64()), nil
}

func Base52Encode(num uint64) string {
	const base = 52
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if num == 0 {
		return string(chars[0])
	}
	result := ""
	for num > 0 {
		remainder := num % base
		result = string(chars[remainder]) + result
		num /= base
	}
	return result
}
