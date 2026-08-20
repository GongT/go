package errors

import (
	"github.com/gongt/go/pkg/errors/internal"
	"github.com/gongt/go/pkg/errors/internal/iterator"
	"github.com/gongt/go/pkg/errors/internal/types"
)

func detailWriter(e error) internal.DetailerWriter {
	switch v := e.(type) {
	case internal.DetailerWriter:
		return v
	default:
		return types.Create(e, 1)
	}
}

// 设置错误的附加数据
//   - 如果e是Err类型，则直接修改它
//   - 如果e不是Err类型，则返回一个新的Err
func SetDetails(e error, details map[string]any) EE {
	if e == nil {
		return nil
	}
	return detailWriter(e).SetDetails(details)
}

// 设置错误的附加数据
//   - 如果e是Err类型，则直接修改它
//   - 如果e不是Err类型，则返回一个新的Err
func SetDetail(e error, key string, value any) EE {
	if e == nil {
		return nil
	}
	return detailWriter(e).SetDetail(key, value)
}

// [SetDetails] 模拟常见标准写法
//
// 但此函数会直接修改e参数，而非创建一个新的错误对象
func WithDetails(e error, detail_pairs ...any) EE {
	if e == nil {
		return nil
	}
	return detailWriter(e).WithDetails(detail_pairs...)
}

// 返回一个包含所有detail的map
//   - 如果err为nil，则返回空map
func GetDetails(err error) map[string]any {
	result := make(map[string]any)
	for details := range iterator.IterEveryDetail(err) {
		for k, v := range details {
			if _, ok := result[k]; ok {
				// 如果已经存在同名的detail，则不覆盖
				continue
			}
			result[k] = v
		}
	}
	return result
}

// 找到details中的code字段，并尝试将其转换为整数
func GetCode(err error) (code int, found bool) {
	return iterator.GetCode(err)
}

// 找到details中的reason字段，并转换成error类型
func GetReason(err error) (reason error, found bool) {
	return iterator.GetReason(err)
}
