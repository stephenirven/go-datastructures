package godatastructures

import (
	"sync"
)

// Generic Heap structure using supplied compare function
type Heap[val any] struct {
	slice   []val
	compare func(v1, v2 val) int
	mutex   sync.RWMutex
}

// constructor
func NewHeap[val any](capacity int, f func(val, val) int) *Heap[val] {
	h := Heap[val]{
		slice:   make([]val, 0),
		compare: f,
		mutex:   sync.RWMutex{},
	}

	return &h
}

// Get the item at the top of the heap without removing it
func (h *Heap[val]) Peek() (value val, ok bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if len(h.slice) == 0 {
		ok = false
		return
	}
	ok = true
	value = h.slice[0]
	return
}

// Get and remove the value from the top of the heap
func (h *Heap[val]) Get() (value val, ok bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if len(h.slice) == 0 {
		ok = false
		return
	}
	ok = true
	value = h.slice[0]
	h.slice[0], h.slice[len(h.slice)-1] = h.slice[len(h.slice)-1], h.slice[0]
	h.slice = h.slice[:len(h.slice)-1]
	h.bubbleDown()

	return
}

// Puts the value on the heap
func (h *Heap[val]) Put(value val) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.slice = append(h.slice, value)
	h.bubbleUp()
}

// Gets the size of the heap
func (h *Heap[val]) Size() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.slice)
}

func (h *Heap[val]) bubbleUp() {

	idx := len(h.slice) - 1

	parentIdx, ok := parent(idx)

	for ok {
		if h.compare(h.slice[parentIdx], h.slice[idx]) > 0 {
			h.slice[parentIdx], h.slice[idx] = h.slice[idx], h.slice[parentIdx]
			idx = parentIdx
			parentIdx, ok = parent(idx)
		} else {
			return
		}
	}
}

func (h *Heap[val]) bubbleDown() {
	idx := 0

	childIdx, ok := h.bestChild(idx)

	for ok {
		cmp := h.compare(h.slice[idx], h.slice[childIdx])
		if cmp > 0 {
			h.slice[idx], h.slice[childIdx] = h.slice[childIdx], h.slice[idx]
			idx = childIdx
			childIdx, ok = h.bestChild(idx)
		} else {
			return
		}
	}

}

// Get the parent index of the supplied index
func parent(idx int) (parentIdx int, ok bool) {
	if idx == 0 {
		ok = false
		return
	}
	ok = true
	parentIdx = (idx - 1) / 2
	return
}

func (h *Heap[val]) bestChild(index int) (child int, ok bool) {
	child1 := (2 * index) + 1
	child2 := child1 + 1

	if child2 >= len(h.slice) {
		if child1 >= len(h.slice) {
			ok = false
			return
		}
		ok = true
		child = child1
		return
	}

	ok = true
	if h.compare(h.slice[child1], h.slice[child2]) > 0 {
		child = child2
		return
	}
	child = child1
	return
}
