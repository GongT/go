package fsys

import (
	"os"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/fsys/fpath"
)

// ReadFileOrEmpty 如果文件存在则读取文件内容，否则返回空字节
//
// 拒绝读取非普通文件（但可以读取符号链接指向的普通文件）
func ReadFileOrEmpty[T fpath.PathLike](ipath T) ([]byte, error) {
	path := fpath.ToString(ipath)
	if stat, err := os.Stat(path); err != nil || stat.IsDir() {
		if os.IsNotExist(err) {
			// 允许不存在
			return nil, nil
		} else {
			return nil, errors.NewAnonymous("读取的文件不是普通文件").WithDetails("path", path)
		}
	} else if stat.Mode().IsRegular() {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				// 允许不存在
				return nil, nil
			}
			return nil, errors.Extend(err, "无法读取输入文件").WithDetails("path", path)
		}
		return data, nil
	} else {
		return nil, errors.NewAnonymous("读取的文件不是普通文件").WithDetails("path", path)
	}
}

// IsEntityExists 判断文件（夹）是否存在
func IsEntityExists[T fpath.PathLike](ipath T) bool {
	path := fpath.ToString(ipath)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// IsEntityExistsL 判断文件（夹）是否存在（使用 Lstat）
func IsEntityExistsL[T fpath.PathLike](ipath T) bool {
	path := fpath.ToString(ipath)
	if _, err := os.Lstat(path); err != nil {
		return false
	}
	return true
}

// IsFileNormalL 判断文件是否是普通文件（非目录、非设备、非套接字、非符号链接等）
func IsFileNormalL[T fpath.PathLike](ipath T) (bool, error) {
	path := fpath.ToString(ipath)
	if stat, err := os.Lstat(path); err != nil {
		return false, err
	} else {
		return stat.Mode().IsRegular(), nil
	}
}

// IsFileNormal 判断文件是否是普通文件（非目录、非设备、非套接字等）
//
// 如果是符号链接，则判断符号链接指向的文件是否是普通文件
func IsFileNormal[T fpath.PathLike](ipath T) (bool, error) {
	path := fpath.ToString(ipath)
	if stat, err := os.Stat(path); err != nil {
		return false, err
	} else {
		return stat.Mode().IsRegular(), nil
	}
}

// EnsureDirExists 确保目录存在，如果不存在则创建
func EnsureDirExists[T fpath.PathLike](ipath T) error {
	path := fpath.ToString(ipath)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return errors.Extend(err, "无法创建目录").WithDetails("path", path)
	}
	return nil
}

// EnsureLinkTarget 确保符号链接linkPath存在且指向targetPath，如果不存在则创建
//
// 如果已存在但不是符号链接，则返回错误
func EnsureLinkTarget[T1 fpath.PathLike, T2 fpath.PathLike](linkPath T1, targetPath T2) error {
	link := fpath.ToString(linkPath)
	target := fpath.ToString(targetPath)

	if stat, err := os.Lstat(link); err != nil {
		if os.IsNotExist(err) {
			// 链接不存在，创建
			if err := os.Symlink(target, link); err != nil {
				return errors.Extend(err, "无法创建符号链接").WithDetails("path", link, "target", target)
			}
			return nil
		} else {
			return errors.Extend(err, "无法获取符号链接信息").WithDetails("path", link)
		}
	} else {
		if stat.Mode()&os.ModeSymlink == 0 {
			return errors.NewAnonymous("路径已存在但不是符号链接").WithDetails("path", link, "mode", ModeName(stat.Mode()))
		}
		// 已存在且是符号链接，检查指向
		currentTarget, err := os.Readlink(link)
		if err != nil {
			return errors.Extend(err, "无法读取符号链接内容").WithDetails("path", link)
		}

		if currentTarget == target {
			// 链接指向正确
			return nil
		}

		if err := os.Remove(link); err != nil {
			return errors.Extend(err, "无法删除错误的符号链接").WithDetails("path", link)
		}
		if err := os.Symlink(target, link); err != nil {
			return errors.Extend(err, "无法创建符号链接").WithDetails("path", link, "target", target)
		}

		return nil
	}
}

// ModeName 返回文件模式的字符串表示
func ModeName(mode os.FileMode) string {
	if mode&os.ModeDir != 0 {
		return "目录"
	} else if mode&os.ModeSymlink != 0 {
		return "符号链接"
	} else if mode&os.ModeNamedPipe != 0 {
		return "命名管道"
	} else if mode&os.ModeSocket != 0 {
		return "套接字"
	} else if mode&os.ModeDevice != 0 {
		return "设备"
	} else if mode&os.ModeCharDevice != 0 {
		return "字符设备"
	} else if mode&os.ModeIrregular != 0 {
		return "未知类型"
	} else {
		return "文件"
	}
}
