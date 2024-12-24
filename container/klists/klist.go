package klists

import (
	"fmt"
	"strings"
)

// Element is an element of a linked list.
type KElement[E comparable] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *KElement[E]

	// The list to which this element belongs.
	list *KList[E]

	// The value stored with this element.
	Value E
}

// Next returns the next list element or nil.
func (e *KElement[E]) Next() *KElement[E] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *KElement[E]) Prev() *KElement[E] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

////////////////////////////////////////////////////////////////////

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type KList[E comparable] struct {
	root KElement[E] // sentinel list element, only &root, root.prev, and root.next are used
	len  int         // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *KList[E]) Init() *KList[E] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New[E comparable]() *KList[E] { return new(KList[E]).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *KList[E]) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *KList[E]) Front() *KElement[E] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *KList[E]) Back() *KElement[E] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *KList[E]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *KList[E]) insert(e, at *KElement[E]) *KElement[E] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	// fmt.Printf("test -- %v\n", e)
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *KList[E]) insertValue(v E, at *KElement[E]) *KElement[E] {
	return l.insert(&KElement[E]{Value: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *KList[E]) remove(e *KElement[E]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

// move moves e to next to at.
func (l *KList[E]) move(e, at *KElement[E]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *KList[E]) Remove(e *KElement[E]) E {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *KList[E]) PushFront(v E) *KElement[E] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

func (l *KList[E]) PushFrontSlice(v ...E) *KElement[E] {
	l.lazyInit()
	if len(v) == 0 {
		return nil
	}
	var ret *KElement[E]
	for _, e := range v {
		ret = l.insertValue(e, &l.root)
	}
	return ret
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *KList[E]) PushBack(v E) *KElement[E] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *KList[E]) PushBackSlice(v ...E) *KElement[E] {
	l.lazyInit()
	if len(v) == 0 {
		return nil
	}
	var ret *KElement[E]
	for _, e := range v {
		l.insertValue(e, l.root.prev)
	}
	return ret
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *KList[E]) InsertBefore(v E, mark *KElement[E]) *KElement[E] {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *KList[E]) InsertAfter(v E, mark *KElement[E]) *KElement[E] {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *KList[E]) MoveToFront(e *KElement[E]) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *KList[E]) MoveToBack(e *KElement[E]) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *KList[E]) MoveBefore(e, mark *KElement[E]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *KList[E]) MoveAfter(e, mark *KElement[E]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *KList[E]) PushBackList(other *KList[E]) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *KList[E]) PushFrontList(other *KList[E]) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

func (l *KList[E]) Clear() {
	var next *KElement[E]
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		l.Remove(e)
	}
}

func (l *KList[E]) FindIf(callback func(v E) bool) *E {
	var next *KElement[E]
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		if nil != callback && callback(e.Value) {
			result := e.Value
			return &result
		}
	}
	return nil
}

func (l *KList[E]) FindAllIf(callback func(v E) bool) []E {
	var next *KElement[E]
	result := make([]E, 0, l.Len())
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		if nil != callback && callback(e.Value) {
			result = append(result, e.Value)
		}
	}
	return result
}

func (l *KList[E]) PopIf(callback func(v E) bool) *E {
	var next *KElement[E]
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		if nil != callback && callback(e.Value) {
			result := e.Value
			l.Remove(e)
			return &result
		}
	}
	return nil
}

func (l *KList[E]) PopAllIf(callback func(v E) bool) []E {
	result := make([]E, 0, l.Len())
	var next *KElement[E]
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		if nil != callback && callback(e.Value) {
			result = append(result, e.Value)
			l.Remove(e)
		}
	}
	return result
}

func (l *KList[E]) PopFront() *E {
	result := l.Remove(l.Front())
	return &result
	// iter := l.Front()
	// if nil == iter {
	// 	return nil
	// }
	// val := iter.Value
	// l.Remove(iter)
	// return &val
}

func (l *KList[E]) PopBack() *E {
	result := l.Remove(l.Back())
	return &result

	// iter := l.Back()
	// if nil == iter {
	// 	return nil
	// }
	// val := iter.Value
	// l.Remove(iter)
	// return &val
}

// 获取指定序号的元素
func (l *KList[E]) At(index int) *E {
	idx := 0
	for e := l.Front(); e != nil; e = e.Next() {
		if idx == index {
			return &e.Value
		}
		idx++
	}
	return nil
}

func (l KList[E]) ToJson5(ident string) string {
	var sb strings.Builder
	var next *KElement[E]
	// sb.WriteString("[\n")
	i := 0
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		sb.WriteString(fmt.Sprintf("%s%#v", ident, e.Value))
		if i < l.Len()-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
		i++
	}
	// sb.WriteString("]")
	return sb.String()
}

func (l KList[E]) ToString() string {
	return l.String()
}

func (l KList[E]) String() string {
	var sb strings.Builder
	var next *KElement[E]
	sb.WriteString("[")
	i := 0
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		sb.WriteString(fmt.Sprintf("%#v", e.Value))
		if i < l.Len()-1 {
			sb.WriteString(",")
		}
		i++
	}
	sb.WriteString("]")
	return sb.String()
}
