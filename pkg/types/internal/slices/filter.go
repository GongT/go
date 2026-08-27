// @exported
package slices

// MapOf 生成一个map，key为keys中的元素，value为value
func MapOf[K comparable, V any](keys []K, value V) map[K]V {
	m := make(map[K]V, len(keys))
	for _, key := range keys {
		m[key] = value
	}
	return m
}

// EmptyMapOf 生成一个map，key为keys中的元素，value为struct{}{}
func EmptyMapOf[K comparable](keys []K) map[K]struct{} {
	m := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		m[key] = struct{}{}
	}
	return m
}

// IntersectInplace 将a中不存在于b的元素移除，即a和b的交集
//
// # 修改a的内容
//
// 重复元素会被移除
func IntersectInplace[T comparable](a *[]T, b []T) {
	set := EmptyMapOf(b)

	n := 0
	for _, item := range *a {
		if _, ok := set[item]; ok {
			(*a)[n] = item
			delete(set, item) // 避免a中重复元素
			n++
		}
	}
	*a = (*a)[:n]
}
