package reader

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"iter"
	"unsafe"
)

type FileInfo struct {
	*ast.File

	fSet *token.FileSet
	file *token.File

	// 文件路径
	Filename string
	RawBytes []byte
	_package *PackageInfo
}

func NewFileInfo(ast *ast.File, fSet *token.FileSet, path string, content []byte, pkg *PackageInfo) FileInfo {
	return FileInfo{
		File:     ast,
		fSet:     fSet,
		file:     fSet.File(ast.FileStart),
		Filename: path,
		RawBytes: content,
		_package: pkg,
	}
}

func (f *FileInfo) CloneAt(sb io.Writer, start token.Pos, end token.Pos) {
	if start == token.NoPos || end == token.NoPos {
		panic(fmt.Sprintf("无效的起止位置: start=%v, end=%v", start, end))
	}

	offsetStart := f.file.Offset(start)
	offsetEnd := f.file.Offset(end)

	sb.Write(f.RawBytes[offsetStart:offsetEnd])
}

func (f *FileInfo) LineCount() int {
	return f.file.LineCount()
}

func (f *FileInfo) Container() *PackageInfo {
	return f._package
}

func (f *FileInfo) Position(pos token.Pos) token.Position {
	return f.fSet.Position(pos)
}

func (f *FileInfo) Lines() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		lastOffset := 0
		for index, line_start := range f.file.Lines()[1:] {
			if lastOffset == 0 {
				lastOffset = line_start
				continue // 第一行
			}

			line := unsafe.String(unsafe.SliceData(f.RawBytes[lastOffset:]), line_start-lastOffset-1)
			if !yield(index, line) {
				return
			}
			lastOffset = line_start
		}
	}
}
