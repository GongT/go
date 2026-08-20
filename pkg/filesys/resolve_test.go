package filesys_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gongt/go/internal/myenv"
	. "github.com/gongt/go/pkg/filesys"
	"github.com/stretchr/testify/assert"
)

var wd, _ = os.Getwd()

func TestResolvePath(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	assert.Equal(t, "/a/b/c/d", MustResolvePath("/a/b", "c", "d"))
	assert.Equal(t, fmt.Sprintf("%s/a/b/c/d", wd), MustResolvePath("a/b", "c", "d"))

	assert.Equal(t, "/c/d", MustResolvePath("/a/b", "/c", "d"))
	assert.Equal(t, "/a/c/d", MustResolvePath("/a/b", "../c", "d"))

}
