package errors

import (
	stderrors "errors"
	"fmt"

	"github.com/gongt/go/pkg/errors/internal"
	"github.com/gongt/go/pkg/errors/internal/types"
)

// 接口
type E = internal.E
type EE = internal.EE

// 错误对象类型
type Err = *types.ErrorObjectBase
type ErrJoin = *types.ErrorObjectJoined
type ErrWrap = *types.ErrorObjectWrapped

// 生成一个新的错误，带有栈信息
//   - 替代[errors.New]
//   - 支持[fmt.Sprintf]格式化
func NewAnonymous(message string, args ...any) Err {
	if len(args) == 0 {
		return types.Create(stderrors.New(message), 1)
	} else {
		return types.Create(fmt.Errorf(message, args...), 1)
	}
}

// [stderrors.New]的别名
//   - 通常用于初始化全局常量
//   - 返回的template本身是一个error，可用于[NewInstance]的message参数
//   - 但更可以直接用[types.ErrorTemplate.New]
func NewTemplate(message string) *types.ErrorTemplate {
	return types.CreateTemplate(message)
}

// NewT 是[NewTemplate]的别名
func NewT(message string) *types.ErrorTemplate {
	return types.CreateTemplate(message)
}

// 生成一个新的错误，带有栈信息
//   - 第一个参数是一个error对象，通常是全局常量，可用于[errors.Is]判断
//   - 此error的字符串将被用于[fmt.Sprintf]格式化
//   - 匿名的用[NewAnonymous]创建
func NewInstance(template error, args ...any) internal.EE {
	if template == nil {
		panic("errors.NewInstance: message is nil")
	}

	var err internal.EE
	switch e := template.(type) {
	case internal.EE:
		err = e
	case *types.ErrorTemplate:
		return e.New(args...)
	default:
		err = types.Create(template, 1)
	}
	err.OverrideMessage(fmt.Sprintf(template.Error(), args...))
	return err
}

// 如果err不是本库的错误类型，则包装它
func Ensure(err error) internal.EE {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case internal.EE:
		return e
	case *types.ErrorTemplate:
		panic("errors.Ensure: 不接受*types.ErrorTemplate类型")
	default:
		return types.Create(err, 1)
	}
}

// 如果err没有栈信息，则为其添加
func EnsureTrace(err error) error {
	return types.CreateEnsureStackTrace(err, 1)
}

// 和[stderrors.Wrap]做相同的事
func Extend(err error, message string, args ...any) ErrWrap {
	if err == nil {
		return nil
	}
	return types.Wrap(1, err, false, message, args)
}

// 将err用template提前定义的错误包装起来
//   - 大多数情况下和普通的Extend效果相同
//   - 唯一区别在于可以用[stderrors.Is]同时判断两个错误
//   - “template”不会出现在任何一条链上
//
// Example:
//
//	tmpl := NewTemplate("无法读取文件%s")
//	_, err := os.ReadFile(file)
//	if err != nil {
//	  err = ExtendWith(err, tmpl, file)
//	}
//	errors.Is(err, os.ErrNotExist) // true
//	errors.Is(err, tmpl) // true
func ExtendWith(err error, template error, args ...any) ErrWrap {
	if err == nil {
		return nil
	}
	e := Extend(err, template.Error(), args...)
	e.AlsoBe(template)
	return e
}

// [Extend]但覆盖栈信息
func ExtendTrace(err error, message string, args ...any) ErrWrap {
	if err == nil {
		return nil
	}
	return types.Wrap(1, err, true, message, args)
}

// [stderrors.Join]
func Join(errs ...error) internal.EE {
	return types.Join(1, errs, false)
}

// [stderrors.Join] 修改“发生了多个错误”消息
func JoinExt(message string, errs ...error) internal.EE {
	return types.JoinMessage(1, message, errs, false)
}

// [JoinWith]
func Concat(errs []error, message string, args ...any) internal.EE {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}
	return types.JoinMessage(1, message, errs, false)
}
