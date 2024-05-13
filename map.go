package godatastructures

import (
	"fmt"
	"hash/maphash"
	"sync"
	"sync/atomic"
	"unsafe"

	"golang.org/x/sync/semaphore"
)

// The seed for the hashing function
var seed = maphash.MakeSeed()

// load factor is a limit to the ratio between size and capacity
// a resize is triggered (in background routine) if the load factor is exceeded
// to prevent the average list length growing too large
const loadFactor int = 5

// Generic hashmap with mutex and growth behaviour
type Map[key comparable, val comparable] struct {
	buckets  []DoublyLinkedList[*MapEntry[key, val]]
	capacity int
	mutex    sync.RWMutex
	resize   *semaphore.Weighted // Indicates resize in progress
	size     int
}

// constructor
func NewMap[key comparable, val comparable](capacity int) *Map[key, val] {

	m := Map[key, val]{
		buckets:  make([]DoublyLinkedList[*MapEntry[key, val]], capacity),
		capacity: capacity,
		resize:   semaphore.NewWeighted(1),
	}

	return &m
}

// Returns the number of key-value mappings in this map.
func (m *Map[key, val]) Size() int {
	if unsafe.Sizeof(m.size) == 8 {
		return int(atomic.LoadInt64((*int64)(unsafe.Pointer(&m.size))))
	} else {
		return int(atomic.LoadInt32((*int32)(unsafe.Pointer(&m.size))))
	}
}

// Structure to store key-value mappings in the map
type MapEntry[key comparable, val any] struct {
	key   key
	value val
}

// Returns the value to which the specified key is mapped,
// or the zero value if not present. Boolean ok indicates
// whether the value was present
func (m *Map[key, val]) Get(k key) (value val, ok bool) {

	// read lock on struct
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// generate the numeric index within the backing slice
	idx := hash(k, m.capacity)

	// read lock the bucket while searching
	m.buckets[idx].mutex.RLock()
	defer m.buckets[idx].mutex.RUnlock()

	// search the list, starting with the most recently added
	entry, ok := m.buckets[idx].FindFirstFunc(func(v *MapEntry[key, val]) bool {
		return v.key == k
	})
	if ok { // item was found
		value = entry.value.value
		return
	}
	ok = false // item was not found
	return

}

// Associates the specified value with the specified key in this map.
func (m *Map[key, val]) Put(k key, v val) {

	// read lock while we search
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// If we need a resize and one isn't already in progress, set it running
	if (m.Size() / m.capacity) > loadFactor {
		if m.resize.TryAcquire(1) {
			go m.grow()
		}
	}

	// get the numeric index
	idx := hash(k, m.capacity)

	// Search for the key, starting with the most recently added
	entry, ok := m.buckets[idx].FindFirstFunc(func(v *MapEntry[key, val]) bool {
		return v.key == k
	})
	if ok {
		entry.value.value = v
		return
	} else {
		entry := MapEntry[key, val]{key: k, value: v}
		m.buckets[idx].AddFirst(&entry)
		if unsafe.Sizeof(m.size) == 8 {
			atomic.AddInt64((*int64)(unsafe.Pointer(&m.size)), 1)
		} else {
			atomic.AddInt32((*int32)(unsafe.Pointer(&m.size)), 1)
		}
	}

}

// Stores all the provided key/value pairs in the map, replacing any
// keys already present
func (m *Map[key, val]) PutAll(entries []MapEntry[key, val]) {
	for e := range entries {
		m.Put(entries[e].key, entries[e].value)
	}
}

// Removes the mapping for the specified key from this map if present.
func (m *Map[key, val]) Remove(k key) (ok bool) {

	// read lock while we search
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// get the numeric index
	idx := hash(k, m.capacity)

	// Search for the key, starting with the most recently added
	entry, ok := m.buckets[idx].FindFirstFunc(func(v *MapEntry[key, val]) bool {
		return v.key == k
	})
	if ok {
		m.buckets[idx].Unlink(entry)
		if unsafe.Sizeof(m.size) == 8 {
			atomic.AddInt64((*int64)(unsafe.Pointer(&m.size)), -1)
		} else {
			atomic.AddInt32((*int32)(unsafe.Pointer(&m.size)), -1)
		}
		return
	}

	return
}

// Removes all of the mappings from this map.
func (m *Map[key, val]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.buckets = make([]DoublyLinkedList[*MapEntry[key, val]], m.capacity)
	if unsafe.Sizeof(m.size) == 8 {
		atomic.StoreInt64((*int64)(unsafe.Pointer(&m.size)), 0)
	} else {
		atomic.StoreInt32((*int32)(unsafe.Pointer(&m.size)), 0)
	}

	m.size = 0

}

// Returns true if this map contains a mapping for the specified key
func (m *Map[key, val]) ContainsKey(k key) (ok bool) {
	// read lock while we search
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// get the numeric index
	idx := hash(k, m.capacity)

	// Search for the key, starting with the most recently added
	_, ok = m.buckets[idx].FindFirstFunc(func(v *MapEntry[key, val]) bool {
		return v.key == k
	})
	return

}

// Returns true if this map maps one or more keys to the specified value
func (m *Map[key, val]) ContainsValue(v val) bool {

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for b := range m.buckets {
		_, ok := m.buckets[b].FindFirstFunc(
			func(e *MapEntry[key, val]) bool {
				return e.value == v
			})
		if ok {
			return true
		}
	}
	return false
}

// Returns a set of the keys in the map
func (m *Map[key, val]) KeySet() *Set[key] {
	s := NewSet[key]()

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for b := range m.buckets {
		m.buckets[b].Do(
			func(e *MapEntry[key, val]) {
				s.Add(e.key)
			})
	}

	return s
}

// Returns a slice of the values in the map
func (m *Map[key, val]) Values() (values []val) {

	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for b := range m.buckets {
		m.buckets[b].Do(
			func(e *MapEntry[key, val]) {
				values = append(values, e.value)
			})
	}
	return
}

// Used internally to grow the backing array
func (m *Map[key, val]) grow() {
	defer m.resize.Release(1) // release the semaphor when done, to allow grow to run again

	// We need a lock on this to prevent access during resize
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if (m.size / m.capacity) < loadFactor {
		return // drop out if we don't need to resize
	}

	newCapacity := newCapacity(m.capacity)

	newBuckets := make([]DoublyLinkedList[*MapEntry[key, val]], newCapacity)

	// build a new set of populated buckets
	for b := range m.buckets {
		m.buckets[b].Do(func(e *MapEntry[key, val]) {
			idx := hash(e.key, newCapacity)
			newBuckets[idx].AddFirst(e)
		})
	}
	m.buckets = newBuckets
	m.capacity = newCapacity

}

// Used internally to hash a key to a int index
func hash[key comparable](k key, capacity int) int {
	var h maphash.Hash
	h.SetSeed(seed)
	h.Write([]byte(fmt.Sprintf("%v", k)))
	return int(h.Sum64() % uint64(capacity))
}

// Used internally to decide (smoothed) growth rate. Figures based on
// https://go.googlesource.com/go/+/2dda92ff6f9f07eeb110ecbf0fc2d7a0ddd27f9d
func newCapacity(currentCapacity int) int {

	if currentCapacity < 256 {
		return 2 * currentCapacity
	}
	if currentCapacity < 512 {
		return int(1.63 * float64(currentCapacity))
	}
	if currentCapacity < 1024 {
		return int(1.44 * float64(currentCapacity))
	}
	if currentCapacity < 2048 {
		return int(1.35 * float64(currentCapacity))
	}
	return int(1.30 * float64(currentCapacity))
}
