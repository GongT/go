package CSI

import (
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	var s St = 0
	assert.Equal(t, "", s.String())

	s.Set(Bold)
	assert.Equal(t, "\x1b[1m", s.String())

	s.Set(Italic)
	assert.Equal(t, "\x1b[1;3m", s.String())

	s.Unset(Bold)
	assert.Equal(t, "\x1b[3m", s.String())

	s.Set(Fore)
	assert.Equal(t, "\x1b[3m", s.String())

	s.Set(Gs13)
	assert.Equal(t, "\x1b[3;38;5;245m", s.String())

	rall := ResetAll
	assert.Equal(t, "\x1b[0m", rall.String())

	rcolor := Reset | Fore | Back
	assert.Equal(t, "\x1b[39;49m", rcolor.String())
}
