package godatastructures

import (
	"cmp"
	"slices"
	"testing"
)

// Example comparison functions
func SortAscendingInt(i, j int) int {
	return cmp.Compare(i, j)
}

func SortDescendingInt(i, j int) int {
	return cmp.Compare(j, i)
}

func TestHeap(t *testing.T) {

	t.Parallel()

	tests := []struct {
		data        []int // initial data to populate heap
		extra       []int // extra data to insert during reads
		compareFunc func(i, j int) int
	}{
		{
			data:        []int{2, 1, 4, 3, 5},
			extra:       []int{7, 2, 10, 30, 0},
			compareFunc: SortAscendingInt,
		},
		{
			data:        []int{2, 1, 4, 4, 5},
			extra:       []int{7, 2, 10, 30, 0},
			compareFunc: SortDescendingInt,
		},
		{
			data:        []int{2},
			extra:       []int{7, 2, 10, 1, 0, 5},
			compareFunc: SortDescendingInt,
		},
		{
			data:        []int{24},
			extra:       []int{},
			compareFunc: SortDescendingInt,
		},
		{
			data:        []int{},
			extra:       []int{7, 2, 10, 30, 0},
			compareFunc: SortDescendingInt,
		},
	}
	t.Log("Given the need to test Heap behaviour on sample data")
	{
		for i, test := range tests {
			t.Logf("\tTest: %d\t When testing source data %v ", i, test.data)
			{

				t.Logf("\t%d\t Testing empty Heap behaviour", i)

				h := NewHeap(10, test.compareFunc)

				v, ok := h.Get()
				if ok {
					t.Errorf("\t%d\t Get on empty heap returned a value: %v", i, v)
				}

				v, ok = h.Peek()
				if ok {
					t.Errorf("\t%d\t Peek on empty heap returned a value: %v", i, v)
				}

				t.Logf("\t%d\t Putting sample data into heap", i)

				for _, n := range test.data {
					h.Put(n)
				}

				// we sort the test source after adding using the supplied func
				// to simulate heap behaviour for testing
				slices.SortFunc(test.data, test.compareFunc)

				if h.Size() != len(test.data) {
					t.Errorf("\t%d\t Size of heap was unexpected %v : %v", i, len(test.data), h.Size())
				}

				v, ok = h.Peek()
				if len(test.data) > 0 && !ok {
					t.Errorf("\t%d\t Peek on populated heap returned no value", i)
				}

				if len(test.data) > 0 && v != test.data[0] {
					t.Errorf("\t%d\t Peek on populated list expected %v : %v", i, test.data[0], v)
				}

				t.Logf("\t%d\t Getting sample data and putting extra data into heap", i)

				ok = true
				for len(test.data) > 0 {

					v, ok = h.Get()
					testVal := test.data[0]
					test.data = test.data[1:]

					if !ok {
						t.Errorf("\t%d\t Get populated heap expected %v", i, testVal)
					}

					if v != testVal {
						t.Errorf("\t%d\t Get expected %v : %v", i, testVal, v)
					}

					// if there's remaining data to add, add some
					if len(test.extra) > 0 {
						additional := test.extra[0]
						test.extra = test.extra[1:]

						h.Put(additional) // add the additional data to the heap
						{                 // add the data into the source slice
							test.data = append(test.data, additional)
							slices.SortFunc(test.data, test.compareFunc)
						}
					}
				}

				if h.Size() != 0 {
					t.Errorf("\t%d\t Size after Get on all values expected 0 : %v", i, h.Size())
				}
				if len(test.data) > 0 {
					t.Errorf("\t%d\t Heap did not return all values from source : %v", i, test.data)
				}

				t.Logf("\t%d\t Testing on remaining extra data", i)

				for _, val := range test.extra {
					h.Put(val)
					test.data = append(test.data, val)
				}
				slices.SortFunc(test.data, test.compareFunc)

				if h.Size() != len(test.data) {
					t.Errorf("\t%d\t Size of heap was unexpected %v : %v", i, len(test.data), h.Size())
				}

				v, ok = h.Peek()
				if len(test.data) > 0 && !ok {
					t.Errorf("\t%d\t Peek on populated heap returned no value", i)
				}

				if len(test.data) > 0 && v != test.data[0] {
					t.Errorf("\t%d\t Peek on populated list expected %v : %v", i, test.data[0], v)
				}

				for len(test.data) > 0 {

					v, ok = h.Get()
					testVal := test.data[0]
					test.data = test.data[1:]

					if !ok {
						t.Errorf("\t%d\t Get populated heap expected %v", i, testVal)
					}

					if v != testVal {
						t.Errorf("\t%d\t Get expected %v : %v", i, testVal, v)
					}
				}

				if h.Size() != 0 {
					t.Errorf("\t%d\t Size of consumed heap was unexpected 0 : %v", i, h.Size())
				}

				v, ok = h.Peek()
				if len(test.data) > 0 && !ok {
					t.Errorf("\t%d\t Peek on empty heap returned a value : %v", i, v)
				}

			}
		}
	}
}
