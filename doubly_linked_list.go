package godatastructures

import "sync"

// Generic Doubly linked list
type DoublyLinkedList[val comparable] struct {
	size  int
	first *DoublyLinkedListNode[val]
	last  *DoublyLinkedListNode[val]
	mutex sync.RWMutex
}

type DoublyLinkedListNode[val comparable] struct {
	value val
	next  *DoublyLinkedListNode[val]
	prev  *DoublyLinkedListNode[val]
	mutex sync.RWMutex
}

// constructor
func NewDoublyLinkedListNode[val comparable](v val) *DoublyLinkedListNode[val] {
	n := DoublyLinkedListNode[val]{
		value: v,
	}

	return &n
}

func (n *DoublyLinkedListNode[val]) Value() val {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.value
}

func (n *DoublyLinkedListNode[val]) SetValue(v val) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.value = v
}

// constructor
func NewDoublyLinkedList[val comparable]() *DoublyLinkedList[val] {
	l := DoublyLinkedList[val]{}
	return &l
}

// Removes all nodes from the list
func (l *DoublyLinkedList[val]) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.size = 0
	l.first = nil
	l.last = nil
}

// Start of list functions

// Adds a new value at the start of the list
func (l *DoublyLinkedList[val]) AddFirst(v val) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.first != nil {
		l.first.mutex.Lock()
		defer l.first.mutex.Unlock()
	}

	n := NewDoublyLinkedListNode(v)
	n.next = l.first

	if l.first != nil {
		l.first.prev = n
	}
	l.first = n
	l.size++
	if l.last == nil {
		l.last = l.first
	}
}

// Removes and returns the value at the start of the list
// boolean value indicates the presence of a value
func (l *DoublyLinkedList[val]) RemoveFirst() (v val, ok bool) {

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.first != nil {
		l.first.mutex.Lock()
		defer l.first.mutex.Unlock()

		ok = true
		v = l.first.value
		l.first = l.first.next
		if l.first != nil {
			l.first.prev = nil
		} else {
			l.last = nil
		}
		l.size--
		return
	}
	ok = false
	v = *new(val)
	return
}

// Returns the value at the start of the list
// boolean value indicates the presence of a value
func (l *DoublyLinkedList[val]) PeekFirst() (v val, ok bool) {

	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if l.first != nil {
		l.first.mutex.RLock()
		defer l.first.mutex.RUnlock()

		v = l.first.value
		ok = true
		return
	}
	v = *new(val)
	ok = false
	return
}

// End of list funcs

// Adds a new value at the end of the list
func (l *DoublyLinkedList[val]) AddLast(v val) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	n := NewDoublyLinkedListNode(v)
	n.prev = l.last
	if l.last != nil {
		l.last.mutex.Lock()
		defer l.last.mutex.Unlock()
		l.last.next = n
		l.last = n
		l.size += 1
	} else {
		// list is empty
		l.first = n
		l.last = n
		l.size = 1
	}
}

// Removes and returns the value at the end of the list
// boolean value indicates the presence of a value
func (l *DoublyLinkedList[val]) RemoveLast() (v val, ok bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.last != nil {
		l.last.mutex.Lock()
		defer l.last.mutex.Unlock()
		if l.last.prev != nil {
			l.last.prev.mutex.Lock()
			defer l.last.prev.mutex.Unlock()
		}

		v = l.last.value
		ok = true
		l.last = l.last.prev
		if l.last != nil {
			l.last.next = nil
		} else {
			l.first = nil
		}
		l.size--
		return
	}
	v = *new(val)
	ok = false
	return
}

// Returns the value at the end of the list
// boolean value indicates the presence of a value
func (l *DoublyLinkedList[val]) PeekLast() (v val, ok bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if l.last != nil {
		l.last.mutex.RLock()
		defer l.last.mutex.RUnlock()

		v = l.last.value
		ok = true
		return
	}
	v = *new(val)
	ok = false
	return

}

// Arbitary position functions

// Adds a new node after an existing node in the list
func (l *DoublyLinkedList[val]) AddAfter(existingNode *DoublyLinkedListNode[val], newNode *DoublyLinkedListNode[val]) {
	if existingNode == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	existingNode.mutex.Lock()
	defer existingNode.mutex.Unlock()
	newNode.mutex.Lock()
	defer newNode.mutex.Unlock()

	next := existingNode.next

	if next != nil {
		next.mutex.Lock()
		defer next.mutex.Unlock()
		next.prev = newNode
	} else {
		l.last = newNode
	}

	newNode.next = next
	newNode.prev = existingNode

	existingNode.next = newNode

	l.size++
}

// Adds a new node before an existing node in the list
func (l *DoublyLinkedList[val]) AddBefore(existingNode *DoublyLinkedListNode[val], newNode *DoublyLinkedListNode[val]) {
	if existingNode == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	existingNode.mutex.Lock()
	defer existingNode.mutex.Unlock()
	newNode.mutex.Lock()
	defer newNode.mutex.Unlock()

	prev := existingNode.prev

	if prev != nil {
		prev.mutex.Lock()
		defer prev.mutex.Unlock()
		prev.next = newNode
	}

	newNode.prev = prev
	newNode.next = existingNode
	existingNode.prev = newNode

	if newNode.prev == nil {
		l.first = newNode
	}

	l.size++
}

// Determines whether a value is in the list
func (l *DoublyLinkedList[val]) Contains(v val) bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	for curr := l.first; curr != nil; curr = curr.next {
		if curr.value == v {
			return true
		}
	}
	return false
}

// Finds the first node that contains the specified value
// boolean value indicates the presence of a value
func (l *DoublyLinkedList[val]) FindFirst(v val) (node *DoublyLinkedListNode[val], ok bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	for curr := l.first; curr != nil; curr = curr.next {
		if curr.value == v {
			return curr, true
		}
	}
	return nil, false
}

// Finds the last node that contains the specified value
// boolean value indicates the presence of a value
func (l *DoublyLinkedList[val]) FindLast(v val) (node *DoublyLinkedListNode[val], ok bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	for curr := l.last; curr != nil; curr = curr.prev {
		if curr.value == v {
			return curr, true
		}
	}
	return nil, false
}

// Finds the first node that matches the predicate
// boolean value indicates the presence of a value
func (l *DoublyLinkedList[val]) FindFirstFunc(predicate func(v val) bool) (node *DoublyLinkedListNode[val], ok bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	for curr := l.first; curr != nil; curr = curr.next {
		if predicate(curr.value) {
			return curr, true
		}
	}
	return nil, false
}

// Finds the last node that matches the predicate
// boolean value indicates the presence of a value
func (l *DoublyLinkedList[val]) FindLastFunc(predicate func(v val) bool) (node *DoublyLinkedListNode[val], ok bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	for curr := l.last; curr != nil; curr = curr.prev {
		if predicate(curr.value) {
			return curr, true
		}
	}
	return nil, false
}

// Disconnects the node from the list
func (l *DoublyLinkedList[val]) Unlink(n *DoublyLinkedListNode[val]) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.next != nil {
		n.next.mutex.Lock()
		defer n.next.mutex.Unlock()
		n.next.prev = n.prev
	}
	if n.prev != nil {
		n.prev.mutex.Lock()
		defer n.prev.mutex.Unlock()
		n.prev.next = n.next
	}

	if l.first == n {
		l.first = n.next
		l.first.prev = nil
	}
	if l.last == n {
		l.last = n.prev
		l.last.next = nil
	}

	n.next = nil
	n.prev = nil

	// reduce length
	l.size--
}

// Moves the node to the first position of the list
func (l *DoublyLinkedList[val]) ToFirst(n *DoublyLinkedListNode[val]) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if n == l.first {
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()
	if n.prev != nil {
		n.prev.mutex.Lock()
		defer n.prev.mutex.Unlock()
	}
	if n.next != nil {
		n.next.mutex.Lock()
		defer n.next.mutex.Unlock()
	}

	if n == l.last {
		l.last = n.prev
		l.last.next = nil
	}

	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}

	n.prev = nil
	n.next = l.first
	l.first.prev = n
	l.first = n
}

// Moves the node to the last position of the list
func (l *DoublyLinkedList[val]) ToLast(n *DoublyLinkedListNode[val]) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if n == l.last {
		return
	}

	if n.prev != nil {
		n.prev.mutex.Lock()
		defer n.prev.mutex.Unlock()
	}
	if n.next != nil {
		n.next.mutex.Lock()
		defer n.next.mutex.Unlock()
	}

	if n == l.first {
		l.first = l.first.next
		l.first.prev = nil
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}

	n.next = nil
	n.prev = l.last
	l.last.next = n
	l.last = n
}

func (l *DoublyLinkedList[val]) Slice() (slice []val) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	slice = make([]val, l.size)
	idx := 0
	for curr := l.first; curr != nil; curr = curr.next {
		slice[idx] = curr.value
		idx++
	}

	return
}

func (l *DoublyLinkedList[val]) FromSlice(slice []val) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.first = nil
	l.last = nil
	l.size = len(slice)
	if len(slice) == 0 {
		return
	}
	n := NewDoublyLinkedListNode(slice[0])
	l.first = n
	for _, v := range slice[1:] {
		n.next = NewDoublyLinkedListNode(v)
		n.next.prev = n
		n = n.next
	}
	l.last = n
}

func (l *DoublyLinkedList[val]) ReverseSlice() (slice []val) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	slice = make([]val, l.size)
	idx := 0
	for curr := l.last; curr != nil; curr = curr.prev {
		slice[idx] = curr.value
		idx++
	}

	return
}

// higher order functions

func (l *DoublyLinkedList[val]) Do(f func(v val)) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for curr := l.first; curr != nil; curr = curr.next {
		f(curr.value)
	}
}

func (l *DoublyLinkedList[val]) DoReverse(f func(v val)) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for curr := l.last; curr != nil; curr = curr.prev {
		f(curr.value)
	}
}
