package fpath

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/require"
)

func panicif(err error) {
	if err != nil {
		panic(err)
	}
}

func test_content(path string, expected string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if string(data) != expected {
		return fmt.Errorf("content mismatch: expected %q, got %q", expected, string(data))
	}
	return nil
}

type Files struct {
	TextFile1    string
	MissingFile1 string
	TextFile2    string

	LinkSimple    string
	LinkSimpleBrk string

	LinkToLink string

	LinkUpper string
	LinkDown  string

	LinkDir        string
	LinkMissingDir string
	LinkDirBrk     string
	ExpectedDir    string

	SelfLink string

	FindThroughLink       string
	FindThroughLink2      string
	FindThroughDirUpper   string
	LinkThroughDirRecurse string
}

func createTestFiles(dir string) Files {
	var r Files

	r.TextFile1 = filepath.Join(dir, "a/b/c/d/e/f/g/file1.txt")
	panicif(os.MkdirAll(filepath.Dir(r.TextFile1), 0o755))
	panicif(os.WriteFile(r.TextFile1, []byte("content1"), 0o600))

	r.TextFile2 = filepath.Join(dir, "a/b/c/file2.txt")
	panicif(os.WriteFile(r.TextFile2, []byte("content2"), 0o600))

	r.LinkSimple = filepath.Join(dir, "a/b/c/d/e/f/g/link-simple")
	panicif(os.Symlink("./file1.txt", r.LinkSimple))

	r.LinkSimpleBrk = filepath.Join(dir, "a/b/c/d/e/f/g/broken-simple")
	panicif(os.Symlink("./missing-file", r.LinkSimpleBrk))
	r.MissingFile1 = filepath.Join(dir, "a/b/c/d/e/f/g/missing-file")

	r.LinkUpper = filepath.Join(dir, "a/b/c/d/e/f/g/link-upper")
	panicif(os.Symlink("../../../../file2.txt", r.LinkUpper))

	r.LinkToLink = filepath.Join(dir, "a/b/c/d/e/f/g/link-to-link")
	panicif(os.Symlink("./link-upper", r.LinkToLink))

	r.LinkDown = filepath.Join(dir, "link-down")
	panicif(os.Symlink("a/b/c/d/e/f/g/file1.txt", r.LinkDown))

	r.LinkDir = filepath.Join(dir, "dir/link-e")
	panicif(os.MkdirAll(filepath.Dir(r.LinkDir), 0o755))
	panicif(os.Symlink("../a/b/c/d/e", r.LinkDir))
	panicif(os.Symlink("./link-e", filepath.Join(dir, "dir/r1")))
	panicif(os.Symlink("./r1", filepath.Join(dir, "dir/r2")))

	r.LinkDirBrk = filepath.Join(dir, "dir/link-broken")
	r.LinkMissingDir = filepath.Join(dir, "dir/ptr_to_missing")
	r.ExpectedDir = filepath.Join(dir, "dir/no-such-dir")
	panicif(os.Symlink("./ptr_to_missing", r.LinkDirBrk))
	panicif(os.Symlink("no-such-dir", r.LinkMissingDir))

	r.SelfLink = filepath.Join(dir, "self-link")
	panicif(os.Symlink("self-link", r.SelfLink))

	r.FindThroughLink = dir + "/dir/link-e/f/g/file1.txt"
	r.FindThroughLink2 = dir + "/dir/link-e/f/g/link-simple"
	r.FindThroughDirUpper = dir + "/dir/link-e/../../d/e/f/g/link-upper"
	r.LinkThroughDirRecurse = filepath.Join(dir, "dir/r2/f/g/link-simple")

	panicif(test_content(r.LinkSimple, "content1"))
	panicif(test_content(r.LinkUpper, "content2"))

	return r
}

func TestRealpath(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	tempDir := t.TempDir()
	files := createTestFiles(tempDir)

	file, _ := myenv.CurrentFileLine()
	_, base := filepath.Split(file)
	p := New(base)
	p.MustRealpath()
	require.Equal(t, file, p.String())

	t.Run("Optional", func(t *testing.T) {
		allCases(t, files, execOptional, 0)
	})

	t.Run("Existing", func(t *testing.T) {
		allCases(t, files, execExisting, 1)
	})

	t.Run("Missing", func(t *testing.T) {
		allCases(t, files, execMissing, 2)
	})
}

func allCases(t *testing.T, files Files, exec func(*Path) error, mode int) {
	var p *Path

	p = New(files.LinkSimple)
	require.NoError(t, exec(p))
	require.Equal(t, files.TextFile1, p.String())

	p = New(files.LinkSimpleBrk)
	if mode == 1 { // Existing模式必须存在
		require.Error(t, exec(p))
		require.Equal(t, files.LinkSimpleBrk, p.String(), "出错时不能修改对象")
	} else {
		require.NoError(t, exec(p))
		require.Equal(t, files.MissingFile1, p.String())
	}

	p = New(files.LinkUpper)
	require.NoError(t, exec(p))
	require.Equal(t, files.TextFile2, p.String())

	p = New(files.LinkToLink)
	require.NoError(t, exec(p))
	require.Equal(t, files.TextFile2, p.String())

	p = New(files.LinkDown)
	require.NoError(t, exec(p))
	require.Equal(t, files.TextFile1, p.String())

	p = New(files.FindThroughLink)
	require.NoError(t, exec(p))
	require.Equal(t, files.TextFile1, p.String())

	p = New(files.FindThroughLink2)
	require.NoError(t, exec(p))
	require.Equal(t, files.TextFile1, p.String())

	p = New(files.FindThroughDirUpper)
	require.NoError(t, exec(p))
	require.Equal(t, files.TextFile2, p.String())

	p = New(files.LinkThroughDirRecurse)
	require.NoError(t, exec(p))
	require.Equal(t, files.TextFile1, p.String())

	f := filepath.Join(files.LinkDirBrk, "file1.txt")
	p = New(f)
	if mode == 2 { // Missing模式允许不存在的路径
		require.NoError(t, exec(p))
		require.Equal(t, filepath.Join(files.ExpectedDir, "file1.txt"), p.String())
	} else {
		require.Error(t, exec(p))
		require.Equal(t, f, p.String(), "出错时不能修改对象")
	}

	p = New(files.SelfLink)
	require.Error(t, exec(p))
	require.Equal(t, files.SelfLink, p.String(), "出错时不能修改对象")
}

func execOptional(p *Path) error {
	return p.Realpath()
}
func execExisting(p *Path) error {
	return p.RealpathExisting()
}
func execMissing(p *Path) error {
	return p.RealpathMissing()
}
