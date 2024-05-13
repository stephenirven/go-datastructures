package godatastructures

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestLRU(t *testing.T) {

	t.Parallel()

	tests := []struct {
		data     []MapEntry[string, int] // slice of ints to put in the source
		capacity int
	}{
		{data: []MapEntry[string, int]{
			{key: "one", value: 1},
			{key: "two", value: 2},
			{key: "three", value: 3},
			{key: "four", value: 4},
			{key: "five", value: 5},
		},
			capacity: 2,
		},
		{data: []MapEntry[string, int]{
			{key: "one", value: 1},
			{key: "two", value: 2},
			{key: "three", value: 3},
			{key: "four", value: 4},
			{key: "five", value: 5},
		},
			capacity: 5,
		},
		{data: []MapEntry[string, int]{
			{key: "one", value: 1},
			{key: "two", value: 2},
			{key: "three", value: 3},
			{key: "four", value: 4},
			{key: "five", value: 5},
		},
			capacity: 10,
		},
		{data: []MapEntry[string, int]{},
			capacity: 2,
		},
	}

	t.Log("Given the need to test LRU behaviour on sample data")
	{
		for i, test := range tests {
			t.Logf("\tTest: %d\t When testing source data %v ", i, test.data)
			{
				t.Logf("\t%d\t Testing empty LRU behaviour", i)

				l, ok := NewLRUCache[string, int](test.capacity)
				if !ok {
					t.Errorf("\t%d\t Create LRU with size %d failed", i, test.capacity)
				}

				if l.Cap() != test.capacity {
					t.Errorf("\t%d\t Capacity on new LRU expected %d : %d", i, test.capacity, l.Cap())
				}

				if l.Len() != 0 {
					t.Errorf("\t%d\t Len on new LRU was %d", i, l.Cap())
				}

				un, ok := l.Get("unknown_key")
				if ok {
					t.Errorf("\t%d\t Empty LRU contained unknown key %v", i, un)
				}

				if l.Len() != 0 {
					t.Errorf("\t%d\t Get on unknown key on empty map affected length %d", i, l.Len())
				}

				t.Logf("\t%d\t Testing on LRU behaviour on sample data", i)

				for _, v := range test.data {
					l.Put(v.key, v.value)
				}

				testLen := int(
					math.Min(
						float64(len(test.data)),
						float64(test.capacity),
					),
				)
				if l.Len() != testLen {
					t.Errorf("\t%d\t Len on LRU expected %d : %d", i, testLen, l.Len())
				}
				// trim overflow on test input
				if len(test.data) > testLen {
					test.data = test.data[len(test.data)-testLen:]
				}

				for _, val := range test.data {
					if !l.Contains(val.key) {
						t.Errorf("\t%d\t LRU expected to contain %v", i, val.key)
					}
				}

				if l.Cap() == l.Len() {

					last := l.list.last

					key := "NewItem_" + fmt.Sprintf("%d", rand.Intn(1000))

					l.Put(key, 1000)

					if !l.Contains(key) {
						t.Errorf("\t%d\t LRU expected to contain %v", i, key)
					}

					if l.Contains(last.value.Key) {
						t.Errorf("\t%d\t LRU expected item to be evicted: %v", i, last.value.Key)
					}

				}

			}
		}
	}
}
