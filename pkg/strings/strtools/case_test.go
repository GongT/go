package strtools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCase(t *testing.T) {
	assert.Equal(t, "HEllo", UcFirst("hEllo"))
	assert.Equal(t, "hello", LcFirst("Hello"))

	assert.Equal(t, "啊", LcFirst("啊"))
	assert.Equal(t, "啊", UcFirst("啊"))

	assert.Equal(t, "Δ啊", UcFirst("δ啊"))
	assert.Equal(t, "Ñ啊", UcFirst("ñ啊"))
}
