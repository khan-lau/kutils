package kmaps

// 定义一个泛型的 map 类型
type AnyMap[K comparable, V any] map[K]V

type ComparableMap[K comparable, V comparable] map[K]V

// 判断map中对应的key是否存在
func HasKey[K comparable, V any](m AnyMap[K, V], key K) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

// 判断map中是否存在对应的val
func FindValue[K comparable, V comparable](m ComparableMap[K, V], val V) (*K, bool) {
	for k, v := range m {
		if v == val {
			return &k, true
		}
	}
	return nil, false
}

// 清除 map中的所有元素
func Clear[K comparable, V any](m AnyMap[K, V]) {
	for k := range m {
		delete(m, k)
	}
}
