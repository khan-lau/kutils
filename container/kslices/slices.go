package kslices

import (
	"strings"
)

// Equal reports whether two slices are equal: the same length and all
// elements equal. If the lengths are different, Equal returns false.
// Otherwise, the elements are compared in increasing index order, and the
// comparison stops at the first unequal pair.
// Floating point NaNs are not considered equal.
func Equal[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[E comparable](s []E, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
func Contains[E comparable](s []E, v E) bool {
	return Index(s, v) >= 0
}

// EqualFunc reports whether two slices are equal using a comparison
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
func EqualFunc[E1, E2 any](s1 []E1, s2 []E2, eq func(E1, E2) bool) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v1 := range s1 {
		v2 := s2[i]
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}

// IndexFunc 是一个泛型函数，用于在切片s中查找第一个满足函数f的元素，并返回其索引
//
// 参数:
//
//	@param E 是切片的元素类型，
//	@param f 是一个接收E类型参数并返回bool值的函数
//
// 返回值:
//
//	@return int -1: 没有找到, >=0: 找到，返回索引值
func IndexFunc[E any](s []E, f func(E) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

// ContainsFunc reports whether at least one
// element e of s satisfies f(e).
func ContainsFunc[E any](s []E, f func(E) bool) bool {
	return IndexFunc(s, f) >= 0
}

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// In the returned slice r, r[i] == v[0].
// Insert panics if i is out of range.
// This function is O(len(s) + len(v)).
func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	tot := len(s) + len(v)
	if tot <= cap(s) {
		s2 := s[:tot]
		copy(s2[i+len(v):], s[i:])
		copy(s2[i:], v)
		return s2
	}
	s2 := make(S, tot)
	copy(s2, s[:i])
	copy(s2[i:], v)
	copy(s2[i+len(v):], s[i:])
	return s2
}

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-j), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete might not modify the elements s[len(s)-(j-i):len(s)]. If those
// elements contain pointers you might consider zeroing those elements so that
// objects they reference can be garbage collected.
func Delete[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j] // bounds check

	return append(s[:i], s[j:]...)
}

// Replace replaces the elements s[i:j] by the given v, and returns the
// modified slice. Replace panics if s[i:j] is not a valid slice of s.
func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	_ = s[i:j] // verify that i:j is a valid subslice
	tot := len(s[:i]) + len(v) + len(s[j:])
	if tot <= cap(s) {
		s2 := s[:tot]
		copy(s2[i+len(v):], s[j:])
		copy(s2[i:], v)
		return s2
	}
	s2 := make(S, tot)
	copy(s2, s[:i])
	copy(s2[i:], v)
	copy(s2[i+len(v):], s[j:])
	return s2
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[S ~[]E, E any](s S) S {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}
	return append(S([]E{}), s...)
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Compact modifies the contents of the slice s; it does not create a new slice.
// When Compact discards m elements in total, it might not modify the elements
// s[len(s)-m:len(s)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage collected.
func Compact[S ~[]E, E comparable](s S) S {
	if len(s) < 2 {
		return s
	}
	i := 1
	for k := 1; k < len(s); k++ {
		if s[k] != s[k-1] {
			if i != k {
				s[i] = s[k]
			}
			i++
		}
	}
	return s[:i]
}

// CompactFunc is like Compact but uses a comparison function.
func CompactFunc[S ~[]E, E any](s S, eq func(E, E) bool) S {
	if len(s) < 2 {
		return s
	}
	i := 1
	for k := 1; k < len(s); k++ {
		if !eq(s[k], s[k-1]) {
			if i != k {
				s[i] = s[k]
			}
			i++
		}
	}
	return s[:i]
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func Grow[S ~[]E, E any](s S, n int) S {
	if n < 0 {
		panic("cannot be negative")
	}
	if n -= cap(s) - len(s); n > 0 {
		// (https://go.dev/issue/53888): Make using []E instead of S
		// to workaround a compiler bug where the runtime.growslice optimization
		// does not take effect. Revert when the compiler is fixed.
		s = append([]E(s)[:cap(s)], make([]E, n)...)[:len(s)]
	}
	return s
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func Clip[S ~[]E, E any](s S) S {
	return s[:len(s):len(s)]
}

// 加工slice中的各个元素, 例如TrimSpace 每个元素
func Process[E any](s []E, callback func(val E) E) {
	for i := 0; i < len(s); i++ {
		s[i] = callback(s[i])
	}
}

// 获取 string的切片, 不足指定长度, 补0x0; 超过指定长度, 截取指定长度
func StringToSlice(s string, length int) []byte {
	b := []byte(s)
	if (len(s)) >= length {
		return b[:length]
	} else {
		slice := make([]byte, length)
		copy(slice, b)
		return slice
	}
}

// 从切换过滤出一个新的切片
func FilterFunc[T comparable](s []T, callback func(val T) bool) []T {
	slice := make([]T, 0, len(s))
	for _, item := range s {
		if callback(item) {
			slice = append(slice, item)
		}
	}
	return slice
}

// StringItemHasPrefix 判断字符串切片 s 中是否存在以 prefix 为前缀的元素
//
// 参数：
//
//	@param s []string - 字符串切片
//	@param prefix string - 前缀字符串
//
// 返回值：
//
//	@return bool - 如果 s 中存在以 prefix 为前缀的元素，则返回 true；否则返回 false
func StringItemHasPrefix(s []string, prefix string) bool {
	for _, item := range s {
		if strings.HasPrefix(item, prefix) {
			return true
		}
	}
	return false
}

// StringItemHasPrefix 判断字符串切片 s 中是否存在以 suffix 为后缀的元素
//
// 参数：
//
//	@param s []string - 字符串切片
//	@param suffix string - 后缀字符串
//
// 返回值：
//
//	@return bool - 如果 s 中存在以 suffix 为后缀的元素，则返回 true；否则返回 false
func StringItemHasSuffix(s []string, suffix string) bool {
	for _, item := range s {
		if strings.HasSuffix(item, suffix) {
			return true
		}
	}
	return false
}

// @bref SplitSliceByLimit 将切片s按照指定的长度limit进行分割，返回分割后的二维切片
//
// 参数：
//
//	@param s：待分割的切片
//	@param limit：每个子切片的最大长度
//
// 返回值：
//
//	@return [][]T：分割后的二维切片
//
// 示例：
//
//	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	limit := 3
//	result := SplitSliceByLimit(s, limit)
//	fmt.Println(result) // 输出：[[1 2 3] [4 5 6] [7 8 9] [10]]
func SplitSliceByLimit[T any](s []T, limit int) [][]T {
	if len(s) <= 0 || limit <= 0 {
		return [][]T{}
	}

	if limit >= len(s) {
		return [][]T{s}
	}

	num := len(s) / limit
	if len(s)%limit != 0 {
		num++
	}
	result := make([][]T, 0, num)
	for i := 0; i < len(s); i += limit {
		end := i + limit
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}
