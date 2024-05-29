package godatastructures

// Generic Doubly linked list
type List[val comparable] struct {
	size  int
	first *Node[val]
	last  *Node[val]
}

type Node[val comparable] struct {
	value val
	next  *Node[val]
	prev  *Node[val]
}

// constructor
func NewNode[val comparable](v val) *Node[val] {
	n := Node[val]{
		value: v,
	}

	return &n
}

// Get the next node or nil
func (n *Node[val]) Next() *Node[val] {
	if n.next != nil {
		return n.next
	}
	return nil
}

// Get the previous node or nil
func (n *Node[val]) Prev() *Node[val] {
	if n.prev != nil {
		return n.prev
	}
	return nil
}

// Get the value on the node
func (n *Node[val]) Value() val {
	return n.value
}

// Set the value on the node
func (n *Node[val]) SetValue(v val) {
	n.value = v
}

// constructor
func NewList[val comparable]() *List[val] {
	l := List[val]{}
	return &l
}

// Get first node or nil
func (l *List[val]) First() *Node[val] {
	return l.first
}

// Get last node or nil
func (l *List[val]) Last() *Node[val] {
	return l.last
}

// Get the size of the list
func (l *List[val]) Size() int {
	return l.size
}

// Remove all nodes from the list
func (l *List[val]) Clear() {
	l.size = 0
	l.first = nil
	l.last = nil
}

// Adds a new value at the start of the list
func (l *List[val]) AddFirst(v val) {

	n := NewNode(v)
	n.next = l.first
	if l.first != nil {

		l.first.prev = n
		l.first = n
		l.size += 1
	} else {
		// list is empty
		l.first = n
		l.last = n
		l.size = 1
	}

}

// Remove and return the value at the start of the list
// boolean ok indicates the presence of a value
func (l *List[val]) RemoveFirst() (v val, ok bool) {

	if l.first != nil {

		ok = true
		v = l.first.value

		if l.first.next != nil {

			l.first = l.first.next
			l.first.prev = nil
			l.size--
		} else {
			// we removed all nodes
			l.first = nil
			l.last = nil
			l.size = 0
		}

		return
	}
	ok = false
	v = *new(val)
	return
}

// Return the value at the start of the list
// boolean ok indicates the presence of a value
func (l *List[val]) PeekFirst() (v val, ok bool) {

	if l.first != nil {
		v = l.first.value
		ok = true
		return
	}
	v = *new(val)
	ok = false
	return
}

// End of list funcs

// Add a new value at the end of the list
func (l *List[val]) AddLast(v val) {

	n := NewNode(v)
	n.prev = l.last
	if l.last != nil {

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

// Remove and returns the value at the end of the list
// boolean ok indicates the presence of a value
func (l *List[val]) RemoveLast() (v val, ok bool) {

	if l.last != nil {

		ok = true
		v = l.last.value

		if l.last.prev != nil {

			l.last = l.last.prev
			l.last.next = nil
			l.size--
		} else {
			// we removed all nodes
			l.first = nil
			l.last = nil
			l.size = 0
		}
		return
	}
	v = *new(val)
	ok = false
	return
}

// Return the value at the end of the list
// boolean ok indicates the presence of a value
func (l *List[val]) PeekLast() (v val, ok bool) {
	if l.last != nil {
		v = l.last.value
		ok = true
		return
	}
	v = *new(val)
	ok = false
	return

}

// Add a new node after an existing node in the list
func (l *List[val]) AddAfter(existingNode *Node[val], newNode *Node[val]) {
	if existingNode == nil {
		return
	}
	next := existingNode.next

	if next != nil {
		next.prev = newNode
	} else {
		l.last = newNode
	}

	newNode.next = next
	newNode.prev = existingNode

	existingNode.next = newNode

	l.size++
}

// Add a new node before an existing node in the list
func (l *List[val]) AddBefore(existingNode *Node[val], newNode *Node[val]) {
	if existingNode == nil {
		return
	}

	prev := existingNode.prev

	if prev != nil {
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

// Determine whether a value is in the list
func (l *List[val]) Contains(v val) bool {

	for curr := l.first; curr != nil; curr = curr.next {
		if curr.value == v {
			return true
		}
	}
	return false
}

// Find the first node that contains the specified value
// boolean ok indicates the presence of a value
func (l *List[val]) FindFirst(v val) (node *Node[val], ok bool) {

	for curr := l.first; curr != nil; curr = curr.next {
		if curr.value == v {
			return curr, true
		}
	}
	return nil, false
}

// Find the last node that contains the specified value
// boolean ok indicates the presence of a value
func (l *List[val]) FindLast(v val) (node *Node[val], ok bool) {

	for curr := l.last; curr != nil; curr = curr.prev {
		if curr.value == v {
			return curr, true
		}
	}
	return nil, false
}

// Find the first node that matches the predicate
// boolean ok indicates the presence of a value
func (l *List[val]) FindFirstFunc(predicate func(v val) bool) (node *Node[val], ok bool) {

	for curr := l.first; curr != nil; curr = curr.next {
		if predicate(curr.value) {
			return curr, true
		}
	}
	return nil, false
}

// Find the last node that matches the predicate
// boolean ok indicates the presence of a value
func (l *List[val]) FindLastFunc(predicate func(v val) bool) (node *Node[val], ok bool) {

	for curr := l.last; curr != nil; curr = curr.prev {
		if predicate(curr.value) {
			return curr, true
		}
	}
	return nil, false
}

// Disconnect the node from the list
func (l *List[val]) Unlink(n *Node[val]) {

	if l.first == n { // if the node is the first
		l.first = n.next // set the first to be the next
	}

	if l.last == n { // if the node is the last
		l.last = n.prev // set the last to be the prev
	}

	if n.next != nil {
		n.next.prev = n.prev
	}
	if n.prev != nil {
		n.prev.next = n.next
	}

	n.next = nil
	n.prev = nil

	l.size--
}

// Move the node to the first position of the list
func (l *List[val]) ToFirst(n *Node[val]) {
	if n == l.first {
		return
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

// Move the node to the last position of the list
func (l *List[val]) ToLast(n *Node[val]) {

	if n == l.last {
		return
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

// Return a slice of the list values
func (l *List[val]) Slice() (slice []val) {

	slice = make([]val, l.size)
	idx := 0
	for curr := l.first; curr != nil; curr = curr.next {
		slice[idx] = curr.value
		idx++
	}

	return
}

// Replace the list values with those from the slice
func (l *List[val]) FromSlice(slice []val) {

	l.first = nil
	l.last = nil
	l.size = len(slice)
	if len(slice) == 0 {
		return
	}
	n := NewNode(slice[0])
	l.first = n
	for _, v := range slice[1:] {
		n.next = NewNode(v)
		n.next.prev = n
		n = n.next
	}
	l.last = n
}

// Return a slice of the list values in reverse order
func (l *List[val]) ReverseSlice() (slice []val) {

	slice = make([]val, l.size)
	idx := 0
	for curr := l.last; curr != nil; curr = curr.prev {
		slice[idx] = curr.value
		idx++
	}

	return
}

// Apply the provided function to each node in the list
// function must NOT modify the list when used with concurrency
func (l *List[val]) Do(f func(v val)) {

	for curr := l.first; curr != nil; curr = curr.next {
		f(curr.value)
	}
}

// Apply the provided function to each node on the list in reverse
// function must NOT modify the list when used with concurrency
func (l *List[val]) DoReverse(f func(v val)) {

	for curr := l.last; curr != nil; curr = curr.prev {
		f(curr.value)
	}
}
