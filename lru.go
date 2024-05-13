package godatastructures

import (
	"sync"
)

type CacheEntry[key comparable, val any] struct {
	Key   key
	Value val
}

// Cache with Least Recently Used eviction policy
type LRUCache[key comparable, val any] struct {
	mutex    sync.RWMutex
	capacity int
	values   *Map[key, *DoublyLinkedListNode[*CacheEntry[key, val]]]
	list     *DoublyLinkedList[*CacheEntry[key, val]]
}

func NewLRUCache[key comparable, val any](cap int) (cache *LRUCache[key, val], ok bool) {

	if cap < 1 {
		return cache, false
	}

	cache = &LRUCache[key, val]{
		mutex:    sync.RWMutex{},
		capacity: cap,
		values:   NewMap[key, *DoublyLinkedListNode[*CacheEntry[key, val]]](1),
		list:     NewDoublyLinkedList[*CacheEntry[key, val]](),
	}
	return cache, true
}

// Return the capacity of the cache
func (l *LRUCache[key, val]) Cap() int {
	return l.capacity
}

// Return the current length of the cache
func (l *LRUCache[key, val]) Len() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.values.Size()
}

// Check if the key is in the cache.
// Does not affect recency
func (l *LRUCache[key, val]) Contains(k key) bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.values.ContainsKey(k)
}

// Return the value associated with the key
// and move that item (if any) to the most recent position
// boolean ok indicates presence of a value
func (l *LRUCache[key, val]) Get(k key) (value val, ok bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock() // List manages its own locks

	existingNode, ok := l.values.Get(k)
	if ok {
		value = existingNode.value.Value
		l.list.ToFirst(existingNode)
		return
	}
	ok = false

	return
}

// Clear all values from the cache
func (l *LRUCache[key, val]) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.list.Clear()
	l.values.Clear()
}

// Put the value into the map with the supplied key
// replacing the value for the key as needed
// Make the item most recent
func (l *LRUCache[key, val]) Put(k key, v val) {
	l.mutex.RLock()
	defer l.mutex.RUnlock() // Map and list manage their own locks

	existingNode, present := l.values.Get(k)
	if present {
		existingNode.value.Value = v
		l.list.ToFirst(existingNode)
		return
	}

	// if we are at capacity, evict the least recent
	if l.values.Size() >= l.capacity {
		last, ok := l.list.RemoveLast()
		if ok {
			l.values.Remove(last.Key)
		}
	}

	// create the new node and make it the most recent
	newNode := CacheEntry[key, val]{
		Key:   k,
		Value: v,
	}

	l.list.AddFirst(&newNode)
	l.values.Put(k, l.list.first)
}
