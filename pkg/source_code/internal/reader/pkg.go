package reader

import "golang.org/x/tools/go/packages"

type PackageInfo struct {
	pkg *packages.Package
}

func NewPackageInfo(pkg *packages.Package) *PackageInfo {
	return &PackageInfo{
		pkg: pkg,
	}
}

// 包名（package xxx）
func (p *PackageInfo) Name() string {
	return p.pkg.Name
}

// 用于导入的包路径，从url开始
// func (p *PackageInfo) ImportPath() string {
// 	return p.pkg.PkgPath
// }

// 包路径
func (p *PackageInfo) Path() string {
	return p.pkg.PkgPath
}
