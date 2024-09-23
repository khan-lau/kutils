package kmaps

// 定义一个泛型的 map 类型
type AnyMap[K comparable, V any] map[K]V

type ComparableMap[K comparable, V comparable] map[K]V

// 获取map中的所有key
func Keys[K comparable, V any](m AnyMap[K, V]) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// 获取map中的所有Values
func Values[K comparable, V any](m AnyMap[K, V]) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

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

// Equal reports whether two maps contain the same key/value pairs.
// Values are compared using ==.
func Equal[K, V comparable](m1 AnyMap[K, V], m2 AnyMap[K, V]) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

// EqualFunc is like Equal, but compares values using eq.
// Keys are still compared with ==.
func EqualFunc[K comparable, V any](m1 AnyMap[K, V], m2 AnyMap[K, V], eq func(V, V) bool) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || !eq(v1, v2) {
			return false
		}
	}
	return true
}

// 清除 map中的所有元素
func Clear[K comparable, V any](m AnyMap[K, V]) {
	for k := range m {
		delete(m, k)
	}
}

func Copy[K comparable, V any](dst AnyMap[K, V], src AnyMap[K, V]) {
	for k, v := range src {
		dst[k] = v
	}
}

// DeleteFunc deletes any key/value pairs from m for which del returns true.
func DeleteFunc[K comparable, V any](m AnyMap[K, V], del func(K, V) bool) {
	for k, v := range m {
		if del(k, v) {
			delete(m, k)
		}
	}
}
