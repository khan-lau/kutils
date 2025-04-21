// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kring

// Ring 是环形链表（或环）的一个元素。
// 环没有起点或终点；指向任何环元素的指针都可以作为整个环的引用。
// 空环表示为 nil Ring 指针。Ring 的零值是一个包含 nil Value 的单元素环。
type Ring[T any] struct {
	next, prev *Ring[T] // 指向下一个和前一个环元素的指针
	Value      T        // 环元素的值
	len        int      // 缓存长度
}

// New 函数创建一个长度为 n 的 Ring 结构体指针，泛型参数 T 为 Ring 存储的数据类型。
// 如果 n 小于等于 0，则返回 nil。
//
// 参数:
//     n: Ring 的长度
//
// 返回值:
//     *Ring[T]: 指向 Ring 结构体的指针，如果 n 小于等于 0，则返回 nil
func New[T any](n int) *Ring[T] {
	if n <= 0 {
		return nil
	}
	// r := new(Ring[T])
	r := &Ring[T]{len: n}
	p := r
	for i := 1; i < n; i++ {
		p.next = &Ring[T]{prev: p, len: n}
		p = p.next
	}
	p.next = r
	r.prev = p
	return r
}

// Prev 返回当前节点的下一个节点
//
// 如果当前节点无后一个节点（即当前节点是环的最后一个节点），则调用 init 方法初始化环，并返回初始化后的环的最后节点
//
// 返回值：
//    *Ring[T]：返回后一个节点，如果没有后一个节点则返回初始化后的环的最后节点
func (that *Ring[T]) Next() *Ring[T] {
	if that == nil {
		return nil
	}
	if that.next == nil {
		return that.init()
	}
	return that.next
}

// Prev 返回当前节点的前一个节点
//
// 如果当前节点没有前一个节点（即当前节点是环的最后一个节点），则调用 init 方法初始化环，并返回初始化后的环的第一个节点
//
// 返回值：
//    *Ring[T]：返回前一个节点，如果没有前一个节点则返回初始化后的环的第一个节点
func (that *Ring[T]) Prev() *Ring[T] {
	if that == nil {
		return nil
	}
	if that.next == nil {
		return that.init()
	}
	return that.prev
}

// Move 方法用于移动 Ring 指针。
//
// 参数：
//     n int：移动步数，可以为正数、负数或零。正数表示向前移动，负数表示向后移动。
//
// 返回值：
//     *Ring[T]：移动后的 Ring 指针。
func (that *Ring[T]) Move(n int) *Ring[T] {
	if that == nil {
		return nil
	}
	if that.next == nil {
		return that.init()
	}

	n = n % that.Len() // 优化：预计算 n % Len()
	switch {
	case n < 0:
		for ; n < 0; n++ {
			that = that.prev
		}
	case n > 0:
		for ; n > 0; n-- {
			that = that.next
		}
	}
	return that
}

// Link 函数用于将两个环（Ring）s 和 that 链接在一起。
//
// 参数：
//     s: *Ring[T] 类型，表示需要链接的另一个环。
//
// 返回值：
//     *Ring[T] 类型，返回链接后的新环的尾节点。
//
// 说明：
//     该函数将环 s 链接到环 that 后面，并返回链接后的新环的尾节点。
//     如果 s 为 nil，则该函数不会进行任何操作并返回 that 的尾节点。
//     如果 that 和 s 是同一个环，则该函数将删除 that 和 s 之间的节点，并更新环的长度。
//     如果 that 和 s 是不同的环，则该函数将两个环合并为一个环，并更新环的长度。
func (that *Ring[T]) Link(s *Ring[T]) *Ring[T] {
	if that == nil {
		return nil
	}
	n := that.Next()
	if s != nil {
		p := s.Prev()
		// 注意：不能使用多重赋值，因为 LHS 的求值顺序未定义。
		that.next = s
		s.prev = that
		n.prev = p
		p.next = n

		// 更新 len
		rLen := that.Len()
		sLen := s.Len()
		if that == s { // 同一环：移除 r 到 s 的节点
			removedLen := 0
			for p := n; p != s; p = p.next {
				removedLen++
			}
			that.updateLen(rLen - removedLen)
			n.updateLen(removedLen)
		} else { // 不同环：合并
			that.updateLen(rLen + sLen)
		}
	}
	return n
}

// Unlink 从当前环中删除第 n 个节点，并返回删除节点后的新环
// 如果当前环为 nil 或者 n <= 0，则返回 nil
// 如果 n 为 0，则返回 nil
// 参数:
//   n: int - 要删除的节点位置
// 返回值:
//   *Ring[T] - 删除节点后的新环
func (that *Ring[T]) Unlink(n int) *Ring[T] {
	if that == nil || n <= 0 {
		return nil
	}
	n = n % that.Len() // 优化：预计算 n % Len()
	if n == 0 {
		return nil
	}
	subRing := that.Link(that.Move(n + 1))
	if subRing != that.Next() { // 更新子环的 len
		subRing.updateLen(n)
	}
	return subRing
}

// Count 返回环形链表中的节点数
//
// 如果环形链表为空，则返回 0
func (r *Ring[T]) Count() int {
	n := 0
	if r != nil {
		n = 1
		for p := r.Next(); p != r; p = p.next {
			n++
		}
	}
	return n
}

func (r *Ring[T]) Len() int {
	if r == nil {
		return 0
	}
	return r.len
}

// Do 按正向顺序对环的每个元素调用函数 f. 如果 f 修改了 *r, Do 的行为未定义
func (r *Ring[T]) Do(f func(any)) {
	if r != nil {
		f(r.Value)
		for p := r.Next(); p != r; p = p.next {
			f(p.Value)
		}
	}
}

///////////////////////////////////////////////////////////////////////////

// init 方法用于初始化一个包含单个元素的环形链表。
//
// 返回值：
//   - *Ring[T]: 指向Ring类型的指针，已初始化为单元素环形链表。
func (that *Ring[T]) init() *Ring[T] {
	that.next = that
	that.prev = that
	that.len = 1 // 单元素环 len=1
	return that
}

// updateLen 函数用于更新 Ring 结构中每个节点的长度。
//
// 参数：
// - newLen int: 更新后的长度。
func (that *Ring[T]) updateLen(newLen int) {
	if that == nil {
		return
	}
	that.len = newLen
	for p := that.Next(); p != that; p = p.next {
		p.len = newLen
	}
}
