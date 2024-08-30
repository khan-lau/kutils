package klists

import (
	"container/list"
)

// 将列表转换为切片
func ToSlice[T any](l *list.List) []T {
	s := make([]T, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value.(T))
	}
	return s
}

// 将列表转换为切片
func ToSliceIf[T any](l *list.List, callback func(item T) bool) []T {
	s := make([]T, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		v, _ := e.Value.(T)
		if callback(v) {
			s = append(s, v)
		}
	}
	return s
}

// 将列表转换为切片
func ToKSlice[T comparable](l *KList[T]) []T {
	s := make([]T, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value)
	}
	return s
}

func FilterFunc[T comparable](l *KList[T], callback func(v T) bool) *KList[T] {
	nl := New[T]()
	for e := l.Front(); e != nil; e = e.Next() {
		if callback(e.Value) {
			nl.PushBack(e.Value)
		}
	}
	return nl
}

// splitKList 将传入的KList类型的切片l按照给定的limit进行拆分，返回拆分后的结果
//
// 参数：
//
//	l *klists.KList[T] - 待拆分的KList类型的切片
//	limit int - 拆分后每个子切片的最大长度
//
// 返回值：
//
//	[]T - 拆分后的结果切片
func SplitKList[T comparable](l *KList[T], limit int) [][]T {
	len := l.Len()

	if limit <= 0 || len <= 0 {
		return [][]T{}
	}

	if len <= limit {
		return [][]T{ToKSlice(l)}
	} else {
		num := len / limit
		if len%limit != 0 {
			num += 1
		}
		result := make([][]T, num)
		idx := 0
		for iter := l.Front(); iter != nil; iter = iter.Next() {
			val := iter.Value
			offset := idx / limit
			if result[offset] == nil {
				result[offset] = make([]T, 0, limit)
			}
			result[offset] = append(result[offset], val)
			idx++
		}
		return result
	}
}
