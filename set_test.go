package godatastructures

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSet(t *testing.T) {

	t.Parallel()

	tests := []struct {
		source      []int // slice of ints to put in the source
		uniqueCount int   // count of unique ints in the slice
	}{
		{source: []int{1, 2, 3, 4, 5}, uniqueCount: 5},
		{source: []int{10, 24, 6, 24, 3, 2}, uniqueCount: 5},
		{source: []int{10}, uniqueCount: 1},
		{source: []int{}, uniqueCount: 0},
		{source: []int{1, 2, 1, 2, 1, 1, 1}, uniqueCount: 2},
		{source: []int{1, 1, 1, 1, 1}, uniqueCount: 1},
		{source: []int{1, 2, 3, 4, 5, 5}, uniqueCount: 5},
	}

	t.Log("Given the need to test Set behaviour on sample data")
	{
		for i, test := range tests {
			t.Logf("\tTest: %d\t When testing source data %v ", i, test.source)
			{
				s := NewSet[int]()

				t.Logf("\t%d\t Testing empty set behaviour", i)

				if s.Size() != 0 {
					t.Errorf("\t%d\t Size on empty Set should be %d : %d", i, 0, s.Size())
				}
				s.Remove(1)
				if s.Size() != 0 {
					t.Errorf("\t%d\t Size after remove on empty Set should be %d : %d", i, 0, s.Size())
				}

				for i := range test.source {
					if s.Contains(test.source[i]) {
						t.Errorf("\t%d\t Empty Set contained item %d", i, test.source[i])
					}
				}

				es := s.Intersection(NewSet[int]())
				if es.Size() != 0 {
					t.Errorf("\t%d\t Intersection of empty Sets should be size 0 %d : %d", i, 0, es.Size())
				}
				es = s.Union(NewSet[int]())
				if es.Size() != 0 {
					t.Errorf("\t%d\t Union of empty Sets should be size 0 %d : %d", i, 0, es.Size())
				}
				es = s.Difference(NewSet[int]())
				if es.Size() != 0 {
					t.Errorf("\t%d\t Difference of empty Sets should be size 0 %d : %d", i, 0, es.Size())
				}

				if len(s.Slice()) != 0 {
					t.Errorf("\t%d\t Slice of empty Sets should be empty %d : %d", i, 0, es.Size())
				}

				s.AddSlice(test.source)

				if s.Size() != test.uniqueCount {
					t.Errorf("\t%d\t Size on populated set should be %d : %d", i, test.uniqueCount, s.Size())
				}

				for i := range test.source {
					if !s.Contains(test.source[i]) {
						t.Errorf("\t%d\t Set from slice should contain all slice elements %d : %d", i, 0, es.Size())
					}
				}

				sl := s.Slice()
				if len(sl) != int(s.Size()) {
					t.Errorf("\t%d\t Slice from set should have the same length %d : %d", i, s.Size(), len(sl))
				}

				if len(test.source) > 0 {
					size := s.Size()
					s.Remove(test.source[0])
					if s.Size() != size-1 {
						t.Errorf("\t%d\t Remove item %d should reduce size %d : %d", i, test.source[0], size-1, s.Size())
						t.Errorf("%v", s.Slice())
						t.Errorf("%#v", s.m.KeySet())

					}
					if s.Contains(test.source[0]) {
						t.Errorf("\t%d\t Remove item should not be present: %d", i, test.source[0])
					}

					s.Add(test.source[0])
					if s.Size() != size {
						t.Errorf("\t%d\t Add item %d should increase size %d : %d", i, test.source[0], size, s.Size())
						t.Errorf("%v", s.Slice())
					}
					if !s.Contains(test.source[0]) {
						t.Errorf("\t%d\t Added item should be present: %d", i, test.source[0])
					}

				}

				t.Logf("\tTest: %d\t Testing union behaviour %v ", i, test.source)
				{
					selfUnion := s.Union(s)

					if selfUnion.Size() != test.uniqueCount {
						t.Errorf("\t%d\t Size on union of set with self should be %d : %d", i, test.uniqueCount, s.Size())
					}

					for i := range test.source {
						if !selfUnion.Contains(test.source[i]) {
							t.Errorf("\t%d\t Self union should contain all slice elements %d : %d", i, 0, es.Size())
						}
					}

					// sort on comparison
					if !cmp.Equal(s.Slice(), selfUnion.Slice(), cmpopts.SortSlices(func(a, b int) bool { return a < b })) {
						t.Errorf("\t%d\t Slice comparison of union on self should be equal to self %v : %v", i, s.Slice(), selfUnion.Slice())
					}

					ns := NewSet[int]()
					ns.Add(-10000)

					union := s.Union(ns)

					if union.Size() != test.uniqueCount+1 {
						t.Errorf("\t%d\t Size on union of set with single element set should be %d : %d", i, test.uniqueCount+1, union.Size())
					}

					for i := range test.source {
						if !union.Contains(test.source[i]) {
							t.Errorf("\t%d\t Union should contain all slice elements %d : %d", i, 0, es.Size())
						}
					}

					if !union.Contains(-10000) {
						t.Errorf("\t%d\t Union should contain element from second set", i)
					}

				}
				{
					selfDiff := s.Difference(s)

					if selfDiff.Size() != 0 {
						t.Errorf("\t%d\t Size on difference of set with self should be %d : %d", i, 0, selfDiff.Size())
					}

					if test.uniqueCount > 1 {
						ds := NewSet[int]()
						ds.Add(test.source[0])
						t.Logf("Diffing %v and %v", s.Slice(), ds.Slice())

						diff := s.Difference(ds)
						t.Logf("Diff: %v", diff.Slice())
						if diff.Size() != test.uniqueCount-1 {
							t.Errorf("\t%d\t Size on difference of set with first element of set should be %d : %d", i, test.uniqueCount-1, diff.Size())
						}

						union := ds.Union(diff)
						if !cmp.Equal(s.Slice(), union.Slice(), cmpopts.SortSlices(func(a, b int) bool { return a < b })) {
							t.Errorf("\t%d\t Union of first and diff should be equal to original %v : %v", i, s.Slice(), union.Slice())
						}
					}

				}
				{
					selfIntersect := s.Intersection(s)
					if selfIntersect.Size() != s.Size() {
						t.Errorf("\t%d\t Size on intersect of set with self should be %d : %d", i, s.Size(), selfIntersect.Size())
					}

					for n := range test.source {
						if !selfIntersect.Contains(test.source[n]) {
							t.Errorf("\t%d\t Self intersect should contain all slice elements %d : %d", i, 0, test.source[n])
						}
					}

					if len(test.source) > 0 {
						is := NewSet[int]()
						is.Add(test.source[0])

						intersect := s.Intersection(is)

						if intersect.Size() != 1 {
							t.Errorf("\t%d\t Size on intersect of set with first element of set should be 1 : %d", i, intersect.Size())
						}

						if !intersect.Contains(test.source[0]) {
							t.Errorf("\t%d\t Intersect should contain element from second set", i)
						}
					}

				}

			}
		}
	}

}
