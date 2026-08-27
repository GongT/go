package fpath

import (
	"path/filepath"
	"strings"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/fsys/fpath/internal"
)

var PathErr = internal.PathErr

type rawer interface {
	Raw() string
}

// 路径对象，所有写操作都会修改原对象
type Path struct {
	value string

	// 规范化缓存
	canonicalizeCache string
}

// 创建一个新的路径对象
func New(ps ...string) *Path {
	p := strings.Join(ps, "/")

	if p == "" {
		p = "."
	}

	return &Path{
		value:             p,
		canonicalizeCache: "",
	}
}

func (p *Path) push(s string) {
	internal.AssertValidPath(s)
	p.value += "/" + s
	p.canonicalizeCache = ""
}
func (p *Path) replace(s string) {
	internal.AssertValidPath(s)
	p.value = s
	p.canonicalizeCache = ""
}

// 返回路径字符串
//
// Deprecated: 应使用Raw()，
// 绝大多数情况下用户程序不需要规范化路径，
// 反而会破坏符号链接
func (p *Path) String() string {
	p.NeedsNormalize()
	return p.canonicalizeCache
}

func (p *Path) Raw() string {
	if myenv.IsDebug {
		if p.value == "" {
			panic(PathErr.New("Raw: 路径为空"))
		}
	}
	return p.value
}

func (p *Path) Immutable() *IPath {
	return &IPath{value: p.Clone()}
}

// 规范化路径
//   - 将路径中的分隔符统一为正斜杠
//   - 并去除多余的分隔符、点号之类的东西
func (p *Path) Normalize() *Path {
	if p.NeedsNormalize() {
		p.value = p.canonicalizeCache
	}
	return p
}

// 路径是否需要规范化
func (p *Path) NeedsNormalize() bool {
	if p.canonicalizeCache == "" {
		p.canonicalizeCache = internal.Clean(p.value)
	}
	return p.canonicalizeCache != p.value
}

// 复制
func (p *Path) Clone() *Path {
	n := &Path{
		value:             p.value,
		canonicalizeCache: p.canonicalizeCache,
	}
	return n
}

// 判断路径是否为绝对路径
func (p *Path) IsAbs() bool {
	return filepath.IsAbs(p.value)
}

// 判断路径是否为根目录（且绝对）
func (p *Path) IsRoot() bool {
	if !p.IsAbs() {
		return false
	}
	p.NeedsNormalize()

	if myenv.IsWindows {
		return filepath.VolumeName(p.canonicalizeCache) == p.canonicalizeCache
	} else {
		return p.canonicalizeCache == "/"
	}
}

// 推入路径片段
func (p *Path) Join[T pathOrFileLike](others ...T) *Path {
	for _, other := range others {
		p.push(ToString(other))
	}
	return p
}

// [Path.Join] 时判断是否为绝对路径，如果是，则直接替换
func (p *Path) Resolve(segments ...string) *Path {
	segments, rooted := internal.ShrinkArguments(segments)
	if rooted {
		p.replace(segments[0])
		segments = segments[1:]
	}
	for _, s := range segments {
		p.push(s)
	}
	return p
}

// [Path.Join] 时判断是否为绝对路径，如果是，则直接替换
func (p *Path) ResolveWith(target rawer) *Path {
	return p.Resolve(target.Raw())
}

// 上一级目录 = Join("..")
func (p *Path) Parent() *Path {
	p.push("..")
	return p
}

// 逻辑上级目录 [path.Dir] 在路径中有符号链接的时候，有些非常特殊的情况会产生错误路径
//
// 建议只在此路径表示一个文件时使用它获取目录
func (p *Path) Dir() *Path {
	pos := strings.LastIndexByte(p.value, '/')
	if pos > 0 {
		p.replace(p.value[:pos])
	} else if pos == 0 {
		p.replace("/")
	} else {
		p.replace(".")
	}
	return p
}

// 修改目录，保留文件名
func (p *Path) SetDir[T PathLike](dir T) *Path {
	d := ToString(dir)
	internal.AssertValidPath(d)
	p.replace(d + "/" + filepath.Base(p.value))
	return p
}

// 各种获取路径的函数
//
// 修改返回的对象不会影响当前路径对象
func (p *Path) Base() *File {
	return NewFile(filepath.Base(p.value))
}

// 替换文件名 就是 .Dir().Join(suffix)
func (p *Path) SetBase[T fileLike](suffix T) *Path {
	suffixStr := ToString(suffix)
	p.Dir().Join(suffixStr)
	return p
}

// deprecated: 改为 [Path.SetBase]
func (p *Path) SetFilename(suffix string) *Path {
	p.Dir().Join(suffix)
	return p
}
