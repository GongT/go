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

func (ip *IPath) String() string {
	return ip.value.String()
}
func (ip *IPath) Raw() string {
	return ip.value.Raw()
}
func (ip *IPath) IsAbs() bool {
	return ip.value.IsAbs()
}
func (ip *IPath) NeedsNormalize() bool {
	// NeedsNormalize虽然会改动状态，但不会影响实际取值
	return ip.value.NeedsNormalize()
}
func (ip *IPath) Normalize() *IPath {
	if ip.value.NeedsNormalize() {
		return &IPath{value: ip.value.Clone().Normalize()}
	} else {
		return ip
	}
}
func (ip *IPath) Clone() *IPath {
	// 不可变的Clone()方法直接返回自身即可
	return ip
}

// Mutate()复制一个可变的Path对象
func (ip *IPath) Mutate() *Path {
	return ip.value.Clone()
}
func (ip *IPath) Join(others ...string) *IPath {
	return &IPath{value: ip.value.Clone().Join(others...)}
}
func (ip *IPath) JoinWith(target rawer) *IPath {
	return &IPath{value: ip.value.Clone().JoinWith(target)}
}
func (ip *IPath) Resolve(others ...string) *IPath {
	return &IPath{value: ip.value.Clone().Resolve(others...)}
}
func (ip *IPath) ResolveWith(target rawer) *IPath {
	return &IPath{value: ip.value.Clone().ResolveWith(target)}
}
func (ip *IPath) Dir() *IPath {
	return &IPath{value: ip.value.Clone().Dir()}
}
func (ip *IPath) SetDir(dir rawer) *IPath {
	return &IPath{value: ip.value.Clone().SetDir(dir)}
}
func (ip *IPath) Base() *File {
	return ip.value.Base()
}
func (ip *IPath) SetBase(suffix *File) *IPath {
	return &IPath{value: ip.value.Clone().SetBase(suffix)}
}
func (ip *IPath) SetFilename(suffix string) *IPath {
	return &IPath{value: ip.value.Clone().SetFilename(suffix)}
}
func (ip *IPath) Realpath() (*IPath, error) {
	p := ip.value.Clone()
	if err := p.Realpath(); err != nil {
		return nil, err
	}
	return &IPath{value: p}, nil
}
func (ip *IPath) RealpathExisting() (*IPath, error) {
	p := ip.value.Clone()
	if err := p.RealpathExisting(); err != nil {
		return nil, err
	}
	return &IPath{value: p}, nil
}
func (ip *IPath) RealpathMissing() (*IPath, error) {
	p := ip.value.Clone()
	if err := p.RealpathMissing(); err != nil {
		return nil, err
	}
	return &IPath{value: p}, nil
}
func (ip *IPath) MustRealpath() *IPath {
	p := ip.value.Clone()
	p.MustRealpath()
	return &IPath{value: p}
}
func (ip *IPath) MustRealpathExisting() *IPath {
	p := ip.value.Clone()
	p.MustRealpathExisting()
	return &IPath{value: p}
}
func (ip *IPath) MustRealpathMissing() *IPath {
	p := ip.value.Clone()
	p.MustRealpathMissing()
	return &IPath{value: p}
}
