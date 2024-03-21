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
