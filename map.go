package godatastructures

import (
	"fmt"
	"hash/maphash"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"
)

// load factor is a limit to the ratio between size and capacity
// To avoid excessive locking complexity, the capacity will grow
// if the loadfactor is reached and a put is called - regardless
// of whether it creates a new key

const loadFactor uint64 = 5 // the limit to the average list length

// Generic hashmap with mutex and growth behaviour
type Map[key comparable, val comparable] struct {
	buckets  []DoublyLinkedList[*MapEntry[key, val]]
	capacity uint64
	mutex    sync.RWMutex
	resize   *semaphore.Weighted // Indicates resize in progress
	size     uint64
	seed     maphash.Seed
}

// constructor
func NewMap[key comparable, val comparable](capacity uint64) *Map[key, val] {

	m := Map[key, val]{
		buckets:  make([]DoublyLinkedList[*MapEntry[key, val]], capacity),
		capacity: capacity,
		resize:   semaphore.NewWeighted(1),
		seed:     maphash.MakeSeed(),
	}

	return &m
}

// Get the size
func (m *Map[key, val]) Size() uint64 {
	return m.size
}

// Values for cache entries
type MapEntry[key comparable, val any] struct {
	key   key
	value val
}

// Get the value associated with the supplied key
// Boolean ok indicates presence of key
func (m *Map[key, val]) Get(k key) (value val, ok bool) {

	// read lock on struct
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// generate the numeric index within the backing slice
	var h maphash.Hash
	h.SetSeed(m.seed)
	h.Write([]byte(fmt.Sprintf("%v", k)))
	idx := h.Sum64() % uint64(m.capacity)

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

func (m *Map[key, val]) Put(k key, v val) {

	// If we need a resize and one isn't already in progress, set it running
	if (m.size / m.capacity) > loadFactor {
		if m.resize.TryAcquire(1) {
			go m.grow()
		}
	}

	// read lock while we search
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// get the numeric index
	var h maphash.Hash
	h.SetSeed(m.seed)
	f := []byte(fmt.Sprintf("%v", k))
	h.Write(f)
	idx := h.Sum64() % m.capacity

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
		atomic.AddUint64(&m.size, uint64(1))
	}

}

func (m *Map[key, val]) grow() {
	defer m.resize.Release(1) // release the semaphor when done, to allow grow to run again
	start := time.Now()

	// We need a lock on this to prevent access during resize
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if (m.size / m.capacity) < loadFactor {
		return // drop out if we don't need to resize
	}
	fmt.Printf("growing size (%d), capacity (%d), rate (%d)\n", m.size, m.capacity, (m.size / m.capacity))

	newCapacity := newCapacity(m.capacity)

	newBuckets := make([]DoublyLinkedList[*MapEntry[key, val]], newCapacity)

	// build a new set of populated buckets
	for b := range m.buckets {
		var h maphash.Hash
		h.SetSeed(m.seed)
		m.buckets[b].Do(func(e *MapEntry[key, val]) {
			h.Reset()
			h.Write([]byte(fmt.Sprintf("%v", e.key)))
			idx := h.Sum64() % uint64(m.capacity)
			newBuckets[idx].AddFirst(e)
		})
	}
	m.buckets = newBuckets
	m.capacity = newCapacity

	elapsed := time.Since(start)
	fmt.Printf("resize to %d took %d us\n", newCapacity, elapsed.Microseconds())

}

// some magic numbers for plausibly smooth growth
// https://go.googlesource.com/go/+/2dda92ff6f9f07eeb110ecbf0fc2d7a0ddd27f9d
func newCapacity(currentCapacity uint64) uint64 {

	if currentCapacity < 256 {
		return 2 * currentCapacity
	}
	if currentCapacity < 512 {
		return uint64(1.63 * float64(currentCapacity))
	}
	if currentCapacity < 1024 {
		return uint64(1.44 * float64(currentCapacity))
	}
	if currentCapacity < 2048 {
		return uint64(1.35 * float64(currentCapacity))
	}
	return uint64(1.30 * float64(currentCapacity))

}
