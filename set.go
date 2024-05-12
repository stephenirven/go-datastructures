package godatastructures

// This class implements a set interface
type Set[val comparable] struct {
	m *Map[val, struct{}]
}

// constructor
func NewSet[val comparable]() *Set[val] {
	s := Set[val]{
		m: NewMap[val, struct{}](1),
	}
	return &s
}

// Returns the number of elements in this set
func (s *Set[val]) Size() int {
	return s.m.Size()
}

func (s *Set[val]) Add(v val) {
	s.m.Put(v, struct{}{})
}

func (s *Set[val]) AddSet(t *Set[val]) {
	for b := range t.m.buckets {
		t.m.buckets[b].Do(
			func(m *MapEntry[val, struct{}]) {
				s.Add(m.key)
			})
	}
}

func (s *Set[val]) AddSlice(t []val) {
	for v := range t {
		s.Add(t[v])
	}
}

func (s *Set[val]) Remove(v val) {
	s.m.Remove(v)
}

func (s *Set[val]) RemoveSlice(values []val) {
	for i := range values {
		s.Remove(values[i])
	}
}

func (s *Set[val]) Contains(v val) bool {
	return s.m.ContainsKey(v)
}

func (s *Set[val]) Union(t *Set[val]) (union *Set[val]) {
	union = NewSet[val]()
	union.AddSet(s)
	union.AddSet(t)
	return
}

func (s *Set[val]) Intersection(t *Set[val]) (intersection *Set[val]) {
	intersection = NewSet[val]()
	for b := range s.m.buckets {
		s.m.buckets[b].Do(
			func(m *MapEntry[val, struct{}]) {
				if t.Contains(m.key) {
					intersection.Add(m.key)
				}
			})
	}
	for b := range t.m.buckets {
		t.m.buckets[b].Do(
			func(m *MapEntry[val, struct{}]) {
				if s.Contains(m.key) {
					intersection.Add(m.key)
				}
			})
	}
	return
}

func (s *Set[val]) Difference(t *Set[val]) (difference *Set[val]) {
	difference = NewSet[val]()
	difference.AddSet(s)
	for b := range t.m.buckets {
		t.m.buckets[b].Do(
			func(m *MapEntry[val, struct{}]) {
				difference.Remove(m.key)
			})
	}

	return
}

// func (s *Set[val]) SubSet(t *Set[val]) bool {
// }

func (s *Set[val]) Slice() []val {
	sl := make([]val, s.m.Size())

	idx := 0
	for b := range s.m.buckets {
		s.m.buckets[b].Do(
			func(m *MapEntry[val, struct{}]) {
				sl[idx] = m.key
				idx++
			})
	}
	return sl
}
