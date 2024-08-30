package ksets

type Empty struct{}
type KSet[T comparable] map[T]Empty

// @bref New 函数用于创建一个新的KSet[T]类型的集合, 其中T为可比较类型
//
// 参数:
//
//	@param cap是一个可选参数，用于指定集合的初始容量, 如果cap未指定，则使用默认容量创建集合
//
// 返回值:
//
//	@return KSet[T]类型的集合
func New[T comparable](cap ...int) KSet[T] {
	var set KSet[T]
	if len(cap) == 0 {
		set = make(KSet[T])
	} else {
		set = make(KSet[T], cap[0])
	}
	return set
}

// @bref Insert 向KSet中插入一个或多个元素
//
// 参数：
//
//	@param m：一个KSet类型的集合
//	@param items：要插入的元素，可变参数，类型需为KSet集合元素的类型T，且T类型需实现comparable接口
func Insert[T comparable](m KSet[T], items ...T) {
	for _, item := range items {
		m[item] = Empty{}
	}
}

// @bref Delete 函数从KSet类型的集合m中删除给定的items元素
//
// 参数：
//
//	@param m：KSet[T]类型的集合，表示要删除元素的集合
//	@param items：T类型的可变参数，表示要删除的元素
func Delete[T comparable](m KSet[T], items ...T) {
	for _, item := range items {
		delete(m, item)
	}
}

// @bref Has 判断给定元素是否在集合中
//
// 参数：
//
//	@param m KSet[T] - 要查找的集合
//	@param item T - 要查找的元素
//
// 返回值：
//
//	@return bool - 如果元素在集合中，返回true；否则返回false
func Has[T comparable](m KSet[T], item T) bool {
	_, ok := m[item]
	return ok
}

// List 函数将KSet[T]类型的集合m中的元素转换成一个T类型的切片并返回, 其中T必须满足comparable接口，即可以进行比较操作
//
// 返回值:
//
//	@return []T T类型的切片，包含了集合m中的所有元素
func List[T comparable](m KSet[T]) []T {
	list := make([]T, 0, len(m))
	for item := range m {
		list = append(list, item)
	}
	return list
}
