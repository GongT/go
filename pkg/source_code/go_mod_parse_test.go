package sourcecode

import (
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors/errfmt"
	"github.com/stretchr/testify/assert"
)

func TestGoModFile(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	me, _ := myenv.CurrentFileLine()

	mod, err := FindGoMod(me)
	errfmt.NoError(t, err)

	modName, err := mod.GetModuleName()
	errfmt.NoError(t, err)
	assert.Equal(t, "github.com/gongt/go", modName)

	meLoc, err := mod.CalculateImportPath(me)
	errfmt.NoError(t, err)
	assert.Equal(t, "github.com/gongt/go/pkg/source_code", meLoc)
}
