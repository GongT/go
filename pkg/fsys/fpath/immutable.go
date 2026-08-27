package fpath

// ImmutablePath 不可变的路径对象
type IPath struct {
	value *Path
}

func INew(p string) *IPath {
	if p == "" {
		p = "."
	}
	obj := &IPath{
		value: New(p),
	}

	return obj
}

// Deprecated: 应使用Raw()
//
// See: [Path.String]
func (ip *IPath) String() string {
	return ip.value.String()
}

// See: [Path.Raw]
func (ip *IPath) Raw() string {
	return ip.value.Raw()
}

// See: [Path.IsAbs]
func (ip *IPath) IsAbs() bool {
	return ip.value.IsAbs()
}

// See: [Path.NeedsNormalize]
func (ip *IPath) NeedsNormalize() bool {
	// NeedsNormalize虽然会改动状态，但不会影响实际取值
	return ip.value.NeedsNormalize()
}

// See: [Path.Normalize]
func (ip *IPath) Normalize() *IPath {
	if ip.value.NeedsNormalize() {
		return &IPath{value: ip.value.Clone().Normalize()}
	} else {
		return ip
	}
}

// See: [Path.Clone]
func (ip *IPath) Clone() *IPath {
	// 不可变的Clone()方法直接返回自身即可
	return ip
}

// Mutate()复制一个可变的Path对象
func (ip *IPath) Mutate() *Path {
	return ip.value.Clone()
}

// See: [Path.Join]
func (ip *IPath) Join[T pathOrFileLike](others ...T) *IPath {
	return &IPath{value: ip.value.Clone().Join(others...)}
}

// See: [Path.Resolve]
func (ip *IPath) Resolve(others ...string) *IPath {
	return &IPath{value: ip.value.Clone().Resolve(others...)}
}

// See: [Path.ResolveWith]
func (ip *IPath) ResolveWith(target rawer) *IPath {
	return &IPath{value: ip.value.Clone().ResolveWith(target)}
}

// See: [Path.Dir]
func (ip *IPath) Dir() *IPath {
	return &IPath{value: ip.value.Clone().Dir()}
}

// See: [Path.Parent]
func (ip *IPath) Parent() *IPath {
	return &IPath{value: ip.value.Clone().Parent()}
}

// See: [Path.SetDir]
func (ip *IPath) SetDir[T PathLike](dir T) *IPath {
	return &IPath{value: ip.value.Clone().SetDir(dir)}
}

// See: [Path.Base]
func (ip *IPath) Base() *File {
	return ip.value.Base()
}

// See: [Path.SetBase]
func (ip *IPath) SetBase[T fileLike](suffix T) *IPath {
	return &IPath{value: ip.value.Clone().SetBase(suffix)}
}

// deprecated: use SetBase instead
func (ip *IPath) SetFilename(suffix string) *IPath {
	return &IPath{value: ip.value.Clone().SetFilename(suffix)}
}

// See: [Path.Realpath]
func (ip *IPath) Realpath() (*IPath, error) {
	p := ip.value.Clone()
	if err := p.Realpath(); err != nil {
		return nil, err
	}
	return &IPath{value: p}, nil
}

// See: [Path.RealpathExisting]
func (ip *IPath) RealpathExisting() (*IPath, error) {
	p := ip.value.Clone()
	if err := p.RealpathExisting(); err != nil {
		return nil, err
	}
	return &IPath{value: p}, nil
}

// See: [Path.RealpathMissing]
func (ip *IPath) RealpathMissing() (*IPath, error) {
	p := ip.value.Clone()
	if err := p.RealpathMissing(); err != nil {
		return nil, err
	}
	return &IPath{value: p}, nil
}

// See: [Path.MustRealpath]
func (ip *IPath) MustRealpath() *IPath {
	p := ip.value.Clone()
	p.MustRealpath()
	return &IPath{value: p}
}

// See: [Path.MustRealpathExisting]
func (ip *IPath) MustRealpathExisting() *IPath {
	p := ip.value.Clone()
	p.MustRealpathExisting()
	return &IPath{value: p}
}

// See: [Path.MustRealpathMissing]
func (ip *IPath) MustRealpathMissing() *IPath {
	p := ip.value.Clone()
	p.MustRealpathMissing()
	return &IPath{value: p}
}
