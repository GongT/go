package deep_init

import (
	"reflect"
	"strings"

	"github.com/gongt/go/pkg/reflection"
)

// 深度初始化一个值
// 返回所有走过的节点的指针（不仅是新添加的）
func DeepInitialize(v any) (results []any) {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Pointer || val.IsNil() || !val.Elem().CanSet() {
		panic("DeepInitialize需要一个非空指针")
	}

	var walk func(v reflect.Value, level int)
	walk = func(v reflect.Value, level int) {
		if !v.CanSet() {
			debugf("%s地址不可写", strings.Repeat("  ", level))
			return
		}

		if v.Kind() == reflect.Pointer {
			var err error
			v, err = reflection.InstantiatePointers(v.Addr())
			if err != nil {
				debugf("%s初始化指针失败: %v", strings.Repeat("  ", level), err)
				return
			}
			if v.IsNil() {
				debugf("%s指针为nil，创建空值并指向它", strings.Repeat("  ", level))
				newVal := reflection.InstantiateType(v.Type())
				v.Set(newVal.Addr())
			}
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			debugf("%s结构体", strings.Repeat("  ", level))
			// if v.CanSet(){}

			for name, field := range v.Fields() {
				debugf("%s- 字段: %s (%s)", strings.Repeat("  ", level), name.Name, field.Type())
				if field.CanSet() {
					walk(field, level+1)
				} else {
					debugf("%s  地址不可写", strings.Repeat("  ", level))
				}
			}
		case reflect.Slice, reflect.Array:
			debugf("%s数组/切片", strings.Repeat("  ", level))
			if v.IsNil() {
				debugf("%s - 初始化", strings.Repeat("  ", level))
				v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			} else {
				debugf("%s - 已有", strings.Repeat("  ", level))
			}
		case reflect.Map:
			debugf("%s映射", strings.Repeat("  ", level))
			if v.IsNil() {
				debugf("%s - 初始化", strings.Repeat("  ", level))
				v.Set(reflect.MakeMap(v.Type()))
			} else {
				debugf("%s - 已有", strings.Repeat("  ", level))
			}
			/* 通道隐式初始化问题太多了，暂不支持 */
			// case reflect.Chan:
		// 	logger.DLog(logger.DEEP_INIT,"%s通道", strings.Repeat("  ", level))
		// 	if v.IsNil() {
		// 		logger.DLog(logger.DEEP_INIT,"%s - 初始化", strings.Repeat("  ", level))
		// 		v.Set(reflect.MakeChan(v.Type(), 0))
		// 	} else {
		// 		logger.DLog(logger.DEEP_INIT,"%s - 已有", strings.Repeat("  ", level))
		// 	}
		default:
			debugf("%s跳过类型: %s", strings.Repeat("  ", level), v.Type())
		}

		results = append(results, v.Addr().Interface())
	}

	walk(val.Elem(), 0)
	return
}
