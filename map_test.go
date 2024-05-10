package godatastructures

import (
	"strconv"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {

	t.Parallel()

	tests := []struct {
		source []MapEntry[string, int] // slice of ints to put in the source
	}{
		{source: []MapEntry[string, int]{
			{key: "one", value: 1},
			{key: "two", value: 2},
			{key: "three", value: 3},
			{key: "four", value: 4},
			{key: "five", value: 5},
		},
		},
		{source: []MapEntry[string, int]{
			{key: "one", value: 1},
			{key: "two", value: 2},
			{key: "three", value: 3},
		},
		},
		{source: []MapEntry[string, int]{
			{key: "one", value: 1},
		},
		},
		{source: []MapEntry[string, int]{}},
	}

	t.Log("Given the need to test Map behaviour on sample data")
	{
		for i, test := range tests {
			t.Logf("\tTest: %d\t When testing source data %v ", i, test.source)
			{

				t.Logf("\t%d\t Testing empty map behaviour", i)

				m := NewMap[string, int](10)

				v, ok := m.Get("non_existent_1")
				if ok {
					t.Errorf("\t%d\t Get on empty map returned a value : %d", i, v)
				}

				ok = m.Remove("non_existent_2")
				if ok {
					t.Errorf("\t%d\t Remove on empty map returned ok", i)
				}

				if m.ContainsKey("non_existent_3") {
					t.Errorf("\t%d\t ContainsKey on empty map returned true", i)
				}

				if m.ContainsValue(-1000) {
					t.Errorf("\t%d\t ContainsValue on empty map returned true", i)
				}
				if len(m.Values()) != 0 {
					t.Errorf("\t%d\t Values on empty map were not empty: %v", i, m.Values())
				}

				if m.KeySet().Size() != 0 {
					t.Errorf("\t%d\t Values on empty map were not empty: %v", i, m.Values())
				}

				if m.Size() != 0 {
					t.Errorf("\t%d\t Operations on empty map affected size: %d", i, m.Size())
				}

				t.Logf("\t%d\t Testing sample set behaviour", i)

				m.PutAll(test.source)

				if m.Size() != uint64(len(test.source)) {
					t.Errorf("\t%d\t PutAll on sample data should be correct length %d : %d", i, len(test.source), m.size)
				}

				for v := range test.source {
					if !m.ContainsKey(test.source[v].key) {
						t.Errorf("\t%d\t Map of sample data should contain key %v", i, test.source[v].key)
					}
					if !m.ContainsValue(test.source[v].value) {
						t.Errorf("\t%d\t Map of sample data should contain value %v", i, test.source[v].value)
					}
				}

				{
					ks := m.KeySet()
					if ks.Size() != uint64(len(test.source)) {
						t.Errorf("\t%d\t KeySet was unexpected size %d : %d", i, ks.Size(), len(test.source))
					}

					for v := range test.source {
						if !ks.Contains(test.source[v].key) {
							t.Errorf("\t%d\t KeySet should contain %v", i, test.source[v].key)
						}
					}
				}

				{
					vs := m.Values()
					if len(vs) != len(test.source) {
						t.Errorf("\t%d\t Value set was unexpected size %d : %d", i, len(vs), len(test.source))
					}

					// n^2 is ok for small test sizes
					values := m.Values()
					for _, sourceVal := range test.source {
						found := false
						for _, val := range values {
							if val == sourceVal.value {
								found = true
							}
						}
						if !found {
							t.Errorf("\t%d\t Values did not contain expected value %v : %v", i, sourceVal, values)
						}
					}

				}
				{
					m.Put("anotherKey", 210)
					if m.Size() != uint64(len(test.source)+1) {
						t.Errorf("\t%d\t Add did not affect size: %d", i, m.Size())
					}

					ks := m.KeySet()
					if !ks.Contains("anotherKey") {
						t.Errorf("\t%d\t Add did not affect KeySet: %d", i, m.Size())
					}

					v, ok := m.Get("anotherKey")
					if !ok {
						t.Errorf("\t%d\t Get on existing key returned not ok", i)
					}
					if v != 210 {
						t.Errorf("\t%d\t Get on existing key returned wrong value: %d : %d", i, 210, v)
					}

					// modify test data for testing purposes
					test.source = append(test.source, MapEntry[string, int]{"anotherKey", 1000})

					m.Put("anotherKey", 1000)
					if m.Size() != uint64(len(test.source)) {
						t.Errorf("\t%d\t Put on existing key affected size: %d", i, m.Size())
					}

					ks = m.KeySet()
					if !ks.Contains("anotherKey") {
						t.Errorf("\t%d\t Add did not affect KeySet: %d", i, m.Size())
					}

					v, ok = m.Get("anotherKey")
					if !ok {
						t.Errorf("\t%d\t Get on existing key returned not ok", i)
					}
					if v != 1000 {
						t.Errorf("\t%d\t Get on existing key returned wrong value: %d : %d", i, 210, v)
					}

					if !m.KeySet().Contains("anotherKey") {
						t.Errorf("\t%d\t KeySet did not contain expected key", i)
					}

					// n^2 is ok for small test sizes
					values := m.Values()
					for _, sourceVal := range test.source {
						found := false
						for _, val := range values {
							if val == sourceVal.value {
								found = true
							}
						}
						if !found {
							t.Errorf("\t%d\t Values did not contain expected value %v : %v", i, sourceVal, values)
						}
					}

				}

			}
		}
	}
}

func TestMapConcurrent(t *testing.T) {
	t.Parallel()

	t.Log("Given the need to test Map behaviour concurrently")
	{
		numAttempts := 500

		t.Logf("Testing Put with %v concurrent", numAttempts)
		m := NewMap[string, int](10)

		wg := sync.WaitGroup{}
		wg.Add(numAttempts)

		for i := range numAttempts {
			i := i // prevent capture - pre go 1.22 feature
			go func() {
				defer wg.Done()
				m.Put(strconv.Itoa(i), i)
			}()
		}

		wg.Wait()

		values := m.Values()
		if len(values) != numAttempts {
			t.Errorf("\t Values did not contain correct number of elements %v : %v :  %v", numAttempts, len(values), (values))
		}

		t.Logf("Testing Put and Remove with %v of each concurrently", numAttempts)

		wg.Add(numAttempts * 2)

		for i := range numAttempts {
			i := i // prevent capture - pre go 1.22 feature
			go func() {
				defer wg.Done()
				m.Remove(strconv.Itoa(i))
			}()
			go func() {
				defer wg.Done()
				m.Put(strconv.Itoa(i+numAttempts), (i + numAttempts))
			}()
		}

		wg.Wait()
		values = m.Values()
		t.Logf("Size: %v", m.Size())

		if len(values) != numAttempts {
			t.Errorf("\t After concurrent add and remove, values did not contain number of elements %v : %v - %v", numAttempts, len(values), values)
		}

		t.Log("Testing the internal lists are consistent")

		totalItems := 0
		for b := range m.buckets {
			itemsInBucket := 0
			m.buckets[b].Do(func(e *MapEntry[string, int]) {
				itemsInBucket++
			})
			totalItems += itemsInBucket

			if itemsInBucket != m.buckets[b].size {
				t.Errorf("\t Bucket %v did not contain the expected number of items : %v :  %v", b, m.buckets[b].size, itemsInBucket)
			}
		}
		if totalItems != int(m.Size()) {
			t.Errorf("\t Map did not contain the expected number of items : %v :  %v", m.Size(), totalItems)
		}

	}

}
