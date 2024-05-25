package godatastructures

import (
	"math/rand"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var concurrency = 5000 // must be divisible by 4
func TestLockingList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		source []int
	}{
		{source: []int{1, 2, 3, 4, 5}},       // !!! Non-negative input ints only !!!
		{source: []int{10, 24, 6, 24, 3, 2}}, // negatives used internally to the tests and to detect non present values
		{source: []int{10}},
		{source: []int{}},
	}

	t.Log("Given the need to test linked list behaviour on sample data")
	{
		for i, test := range tests {
			t.Logf("\tTest: %d\t When testing source data %v ", i, test.source)
			{

				t.Logf("\t%d\t Testing empty list behaviour", i)

				l := NewLList[int]()

				v, ok := l.RemoveFirst()
				if ok {
					t.Errorf("\t%d\t RemoveFirst on empty list should not return a value - returned %t : %d", i, ok, v)
				}

				v, ok = l.RemoveLast()
				if ok {
					t.Errorf("\t%d\t RemoveLast on empty list should not return a value - returned %t : %d", i, ok, v)
				}

				v, ok = l.PeekFirst()
				if ok {
					t.Errorf("\t%d\t PeekFirst on empty list should not return a value - returned %t : %d", i, ok, v)
				}

				v, ok = l.PeekLast()
				if ok {
					t.Errorf("\t%d\tPeekLast on empty list should not return a value - returned %t : %d", i, ok, v)
				}

				n, ok := l.FindFirst(123)
				if ok {
					t.Errorf("\t%d\t FindFirst on empty list should not return a value - returned %t : %#v", i, ok, n)
				}

				n, ok = l.FindLast(123)
				if ok {
					t.Errorf("\t%d\t FindLast on empty list should not return a value - returned %t : %#v", i, ok, n)
				}

				ok = l.Contains(123)
				if ok {
					t.Errorf("\t%d\t Contains on empty list should not return a value - returned %t", i, ok)
				}

				t.Logf("\t%d\t Testing loading a list from a slice", i)

				l.FromSlice(test.source)

				if len(test.source) != l.Size() {
					t.Errorf("\t%d\t Length should be %d : %d", i, len(test.source), l.Size())
				}

				ts := l.Slice()

				if !cmp.Equal(test.source, ts) {
					t.Errorf("\t%d\t To / from slice should be equal %v : %v", i, test.source, ts)
				}

				slices.Reverse(test.source)
				tsr := l.ReverseSlice()

				if !cmp.Equal(test.source, tsr) {
					t.Errorf("\t%d\t To / from slice reversed should be equal %v : %v", i, test.source, tsr)
				}

				// bring the test data back to original order
				slices.Reverse(test.source)

				if l.First() != nil && l.First().Value() != test.source[0] {
					t.Errorf("\t%d\t First of list is wrong %v : %v", i, test.source[0], l.First().Value())
				}

				peekFirst, ok := l.PeekFirst()

				if l.Size() > 0 && !ok {
					t.Errorf("\t%d\t Peek() returned no value", i)
				}

				if l.Size() == 0 && ok {
					t.Errorf("\t%d\t Peek() returned a value on an empty list", i)
				}

				if l.First() != nil && peekFirst != test.source[0] {
					t.Errorf("\t%d\t PeekFirst of list is wrong %v : %v", i, test.source[0], peekFirst)
				}

				if l.Last() != nil && l.Last().Value() != test.source[len(test.source)-1] {
					t.Errorf("\t%d\t Last of list is wrong %v : %v", i, test.source[len(test.source)-1], l.Last().Value())
				}

				peekLast, ok := l.PeekLast()

				if l.Size() > 0 && !ok {
					t.Errorf("\t%d\t PeekLast() returned no value", i)
				}

				if l.Size() == 0 && ok {
					t.Errorf("\t%d\t PeekLast() returned a value for empty list : %v", i, peekLast)
				}

				if l.Size() > 0 && peekLast != test.source[len(test.source)-1] {
					t.Errorf("\t%d\t PeekLast of list is wrong %v : %v", i, test.source[len(test.source)-1], peekLast)
				}

				t.Logf("\t%d\t Testing adding to start of list", i)

				l.AddFirst(-1)
				peekFirst, ok = l.PeekFirst()

				if !ok || peekFirst != -1 {
					t.Errorf("\t%d\t PeekFirst of list is wrong %v : %v", i, -1, peekFirst)
				}

				l.AddFirst(-2)
				peekFirst, ok = l.PeekFirst()

				if !ok || peekFirst != -2 {
					t.Errorf("\t%d\t PeekFirst of list is wrong %v : %v", i, -2, peekFirst)
				}

				l.AddFirst(-3)
				peekFirst, ok = l.PeekFirst()

				if !ok || peekFirst != -3 {
					t.Errorf("\t%d\t PeekFirst of list is wrong %v : %v", i, -3, peekFirst)
				}

				if l.Size() != len(test.source)+3 {
					t.Errorf("\t%d\t Size of list after AddFirst is wrong %v : %v", i, len(test.source)+3, l.Size())
				}

				t.Logf("\t%d\t Testing adding to end of list", i)

				l.AddLast(1e5)
				peekLast, ok = l.PeekLast()

				if !ok || peekLast != 1e5 {
					t.Errorf("\t%d\t PeekLast of list is wrong %v : %v", i, 1e5, peekLast)
				}

				l.AddLast(1e6)
				peekLast, ok = l.PeekLast()

				if !ok || peekLast != 1e6 {
					t.Errorf("\t%d\t PeekLast of list is wrong %v : %v", i, 1e6, peekLast)
				}

				l.AddLast(1e7)
				peekLast, ok = l.PeekLast()

				if !ok || peekLast != 1e7 {
					t.Errorf("\t%d\t PeekLast of list is wrong %v : %v", i, 1e7, peekLast)
				}

				if l.Size() != len(test.source)+6 {
					t.Errorf("\t%d\t Size of list after addfirst is wrong %v : %v", i, len(test.source)+6, l.Size())
				}

				t.Logf("\t%d\t Testing removing from start of list", i)

				v, ok = l.RemoveFirst()
				if !ok || v != -3 {
					t.Errorf("\t%d\t RemoveFirst should return an item (%t) %v : %v", i, ok, -3, v)
				}

				fn := l.First()

				if fn.Value() != -2 {
					t.Errorf("\t%d\t First item should have value %v : %v", i, -2, v)
				}

				if fn.Prev() != nil {
					t.Errorf("\t%d\t First item should not have a prev node : %v", i, fn.Prev())
				}

				if l.Size() != len(test.source)+5 {
					t.Errorf("\t%d\t Size of list after RemoveFirst is wrong %v : %v", i, len(test.source)+5, l.Size())
				}

				t.Logf("\t%d\t Testing removing from end of list", i)

				v, ok = l.RemoveLast()
				if !ok || v != 1e7 {
					t.Errorf("\t%d\t RemoveLast should return an item (%t) %v : %v", i, ok, 1e7, v)
				}

				if l.Last().Value() != 1e6 {
					t.Errorf("\t%d\t First item should have value %v : %v", i, -2, v)
				}

				if l.Last().Next() != nil {
					t.Errorf("\t%d\t Last item should not have a next node : %v", i, l.Last().Next())
				}

				if l.Size() != len(test.source)+4 {
					t.Errorf("\t%d\t Size of list after addfirst is wrong %v : %v", i, len(test.source)+4, l.Size())
				}

				t.Logf("\t%d\t Testing Find", i)

				sl, ok := l.FindLast(1e5)
				if !ok || sl.Value() != 1e5 {
					t.Errorf("\t%d\t Second last value (%t) should be %v : %v", i, ok, 1e5, sl.Value())
				}

				sl2, ok := l.FindFirst(1e5)
				if !ok || sl2.Value() != 1e5 {
					t.Errorf("\t%d\t Second last value (%t) should be %v : %v", i, ok, 1e5, sl2.Value())
				}

				if sl2 != sl {
					t.Errorf("\t%d\t FindFirst and FindLast should return the same node %v : %v", i, &sl, &sl2)
				}

				if sl.Next() != l.Last() || l.Last().Next() != nil {
					t.Errorf("\t%d\t Item after the found item should be the last %v : %v", i, sl, l.Last())
				}

				t.Logf("\t%d\t Testing Contains", i)

				if !l.Contains(l.First().Value()) || !l.Contains(l.Last().Value()) || !l.Contains(sl.Value()) {
					t.Errorf("\t%d\t List should contain first, last and second last item", i)
				}

				t.Logf("\t%d\t Testing find last on non present value", i)

				npv := -1e6 - rand.Intn(200) // non present value
				nf, ok := l.FindLast(npv)

				if ok {
					t.Errorf("\t%d\t Non present value found last %v : %v", i, npv, nf.value)
				}

				t.Logf("\t%d\t Testing find first on non present value", i)

				nf, ok = l.FindFirst(npv)
				if ok {
					t.Errorf("\t%d\t Non present value found first %v : %v", i, npv, nf.value)
				}

				t.Logf("\t%d\t Testing adding a node to the start", i)
				newFirst := NewLNode(-4)

				l.AddBefore(l.First(), newFirst)

				if l.First().Value() != -4 {
					t.Errorf("\t%d\t AddBefore first should cause the node to be added as first", i)
				}
				if v, ok = l.PeekFirst(); !ok {
					t.Errorf("\t%d\t AddBefore first should cause the node to be added as first", i)
				}
				if v != -4 {
					t.Errorf("\t%d\t AddBefore first should cause the node to be added as first", i)
				}

				if l.First().Prev() != nil {
					t.Errorf("\t%d\t AddBefore first should have nil prev", i)
				}

				if l.First().Next().Value() != -2 {
					t.Errorf("\t%d\t AddBefore first should have correct next", i)
				}

				if l.First().Next().Prev() != newFirst {
					t.Errorf("\t%d\t AddBefore first should have next-prev correct", i)
				}

				if l.Size() != len(test.source)+5 {
					t.Errorf("\t%d\t Size of list after AddBefore first is wrong %v : %v", i, len(test.source)+5, l.size)
				}

				t.Logf("\t%d\t Testing adding before second node", i)
				newSecond := NewLNode(-3)

				l.AddBefore(l.First().Next(), newSecond)

				if l.First().Next().Value() != -3 {
					t.Errorf("\t%d\t AddBefore second should cause the node to be added as second", i)
				}

				if l.First().Next() != newSecond {
					t.Errorf("\t%d\t AddBefore second should cause the node to be added as second", i)
				}

				if newSecond.Prev() != l.First() {
					t.Errorf("\t%d\t AddBefore second should have first as prev", i)
				}

				if newSecond.Next().Value() != -2 {
					t.Errorf("\t%d\t AddBefore second should have correct next", i)
				}

				if newSecond.Next().Prev() != newSecond {
					t.Errorf("\t%d\t AddBefore second should have next-prev correct", i)
				}

				if newSecond.Prev().Next() != newSecond {
					t.Errorf("\t%d\t AddBefore second should have next-prev correct", i)
				}

				if l.Size() != len(test.source)+6 {
					t.Errorf("\t%d\t Size of list after AddBefore second is wrong %v : %v", i, len(test.source)+6, l.Size())
				}

				t.Logf("\t%d\t Testing adding node before last", i)
				newSecondLast := NewLNode(500000)

				l.AddBefore(l.Last(), newSecondLast)

				if l.Last().Prev().Value() != 500000 {
					t.Errorf("\t%d\t AddBefore last should cause the node to be added as second last", i)
				}

				if l.Last().Prev() != newSecondLast {
					t.Errorf("\t%d\t AddBefore last should cause the node to be added as second last", i)
				}

				if newSecondLast.Next() != l.Last() {
					t.Errorf("\t%d\t AddBefore second should have first as prev", i)
				}

				if newSecondLast.Next().Value() != 1e6 {
					t.Errorf("\t%d\t AddBefore second should have correct next", i)
				}

				if newSecondLast.Prev().Next() != newSecondLast {
					t.Errorf("\t%d\t AddBefore second should have prev-next correct", i)
				}

				if newSecondLast.Next().Prev() != newSecondLast {
					t.Errorf("\t%d\t AddBefore second should have prev-next correct", i)
				}

				if l.Size() != len(test.source)+7 {
					t.Errorf("\t%d\t Size of list after AddBefore last is wrong %v : %v", i, len(test.source)+7, l.Size())
				}

				t.Logf("\t%d\t Testing adding after last node", i)
				newLast := NewLNode(10000000)
				newLast.next = newFirst
				newLast.prev = newFirst

				l.AddAfter(l.Last(), newLast)

				if l.Last().Value() != 10000000 {
					t.Errorf("\t%d\t AddAfter last should cause the node to be added as last", i)
				}

				if l.Last().Next() != nil {
					t.Errorf("\t%d\t AddAfter last should have nil next", i)
				}

				if l.Last().Prev().Value() != 1000000 {
					t.Errorf("\t%d\t AddAfter last should have correct prev", i)
				}

				if l.Last().Prev().Next() != newLast {
					t.Errorf("\t%d\t AddAfter last should have prev-next correct", i)
				}

				if l.Size() != len(test.source)+8 {
					t.Errorf("\t%d\t Size of list after AddLast last is wrong %v : %v", i, len(test.source)+8, l.Size())
				}

				t.Logf("\t%d\t Testing Unlink on first node", i)

				f1 := l.First()
				l.Unlink(f1)

				if f1.Next() != nil || f1.Prev() != nil {
					t.Errorf("\t%d\t Unlink first should remove prev/next", i)
				}

				if l.First().Prev() != nil {
					t.Errorf("\t%d\t Unlink first should ensure first has no prev", i)
				}

				if l.First().Value() != -3 {
					t.Errorf("\t%d\t Value of first after unlink first is wrong %v : %v", i, -3, l.First().Value())
				}

				if l.Size() != len(test.source)+7 {
					t.Errorf("\t%d\t Size of list after Unlink first is wrong %v : %v", i, len(test.source)+7, l.Size())
				}

				t.Logf("\t%d\t Testing Unlink on last node", i)

				l1 := l.Last()
				l.Unlink(l1)

				if l1.Next() != nil || l1.Prev() != nil {
					t.Errorf("\t%d\t Unlink last should remove prev/next", i)
				}

				if l.Last().Next() != nil {
					t.Errorf("\t%d\t Unlink last should ensure last has no next", i)
				}

				if l.Last().Value() != 1e6 {
					t.Errorf("\t%d\t Value of last after unlink last is wrong %v : %v", i, 1e6, l.Last().Value())
				}

				if l.Size() != len(test.source)+6 {
					t.Errorf("\t%d\t length of list after Unlink last is wrong %v : %v", i, len(test.source)+6, l.Size())
				}

				t.Logf("\t%d\t Testing UnLink on node in list body", i)

				m1 := l.Last().Prev()
				l.Unlink(m1)

				if m1.Next() != nil || m1.Prev() != nil {
					t.Errorf("\t%d\t Unlink second last should remove prev/next", i)
				}

				if l.Last().Next() != nil {
					t.Errorf("\t%d\t Unlink second last should ensure last has no next", i)
				}

				if l.Last().Value() != 1e6 {
					t.Errorf("\t%d\t Value of last after unlink last is wrong %v : %v", i, 1e6, l.Last().Value())
				}

				if l.Last().Prev() == nil {
					t.Errorf("\t%d\t Value of last prev after unlink second last is not set", i)
				}

				if l.Last().Prev().Next() != l.Last() {
					t.Errorf("\t%d\t Value of last prev next after unlink second last is not correct %#v : %#v", i, l.Last(), l.Last().Prev().Next())
				}

				if l.Size() != len(test.source)+5 {
					t.Errorf("\t%d\t Size of list after Unlink last is wrong %v : %v", i, len(test.source)+5, l.Size())
				}

				t.Logf("\t%d\t Testing Slice behaviour", i)

				forwardSlice := l.Slice()
				reverseSlice := l.ReverseSlice()

				if l.Size() != len(forwardSlice) || l.Size() != len(reverseSlice) {
					t.Errorf("\t%d\t Size of slices does not match size %v : %v : %v", i, l.size, len(forwardSlice), len(reverseSlice))
				}

				var sourceData = []int{-3, -2, -1}
				sourceData = append(sourceData, test.source...)
				sourceData = append(sourceData, 1e5)
				sourceData = append(sourceData, 1e6)

				if l.Size() != len(sourceData) {
					t.Errorf("\t%d\t Size of source data does not match size %v : %v", i, l.Size(), len(sourceData))
				}

				if !cmp.Equal(forwardSlice, sourceData) {
					t.Errorf("\t%d\t Forward slice is not equal to source data %#v : %v ", i, sourceData, forwardSlice)
				}

				slices.Reverse(sourceData)

				if !cmp.Equal(reverseSlice, sourceData) {
					t.Errorf("\t%d\t Reverse slice is not equal to reversed source data %#v : %v ", i, sourceData, forwardSlice)
				}

				t.Logf("\t%d\t Testing moving the first node", i)
				{
					moveFirst := l.Slice()
					l.ToFirst(l.First())

					if !cmp.Equal(moveFirst, l.Slice()) {
						t.Errorf("\t%d\t Moving first node to first shouldn't change list %#v : %v ", i, moveFirst, l.Slice())
					}

					moveFirst = append(moveFirst, moveFirst[0])[1:]
					l.ToLast(l.First())

					if !cmp.Equal(moveFirst, l.Slice()) {
						t.Errorf("\t%d\t Moving first node to last should reflect in the shouldn't change list %#v : %v ", i, moveFirst, l.Slice())
					}

				}

				t.Logf("\t%d\t Testing moving the last node", i)
				{
					moveLast := l.Slice()
					l.ToLast(l.Last())

					if !cmp.Equal(moveLast, l.Slice()) {
						t.Errorf("\t%d\t Moving last node to last shouldn't change list %#v : %v ", i, moveLast, l.Slice())
					}

					moveLast = append([]int{moveLast[len(moveLast)-1]}, moveLast...)
					moveLast = moveLast[:len(moveLast)-1]
					l.ToFirst(l.Last())

					if !cmp.Equal(moveLast, l.Slice()) {
						t.Errorf("\t%d\t Moving last item to first should reflect in the shouldn't change list %#v : %v ", i, moveLast, l.Slice())
					}

				}

				if l.Size() > 2 { // if there are more than 2 nodes
					t.Logf("\t%d\t Testing moving a node from the list body", i)
					moveMid := l.Slice()
					l.ToLast(l.Last().Prev())

					midVal := moveMid[len(moveMid)-2]
					moveMid = append(moveMid[0:len(moveMid)-2], moveMid[len(moveMid)-1:]...)
					moveMid = append(moveMid, midVal)
					if !cmp.Equal(moveMid, l.Slice()) {
						t.Errorf("\t%d\t Moving last item to last shouldn't change list %#v : %v ", i, moveMid, l.Slice())
					}

					lastVal := moveMid[len(moveMid)-1]
					moveMid = append([]int{moveMid[len(moveMid)-2]}, moveMid[0:len(moveMid)-2]...)
					moveMid = append(moveMid, lastVal)

					l.ToFirst(l.Last().Prev())

					if !cmp.Equal(moveMid, l.Slice()) {
						t.Errorf("\t%d\t Moving last item to first should reflect in the shouldn't change list %#v : %v ", i, moveMid, l.Slice())
					}

				}

				t.Logf("\t%d\t Testing the Clear functionality", i)

				l.Clear()
				if l.Size() != 0 {
					t.Errorf("\t%d\t After clear, list sould have size 0 : %v", i, l.Size())
				}

				if l.First() != nil {
					t.Errorf("\t%d\t After clear, list sould have no first node : %#v", i, l.First())
				}

				if l.Last() != nil {
					t.Errorf("\t%d\t After clear, list sould have no last node : %#v", i, l.Last())
				}

			}
		}
	}
}

func TestLockingListSingleItem(t *testing.T) {

	t.Log("Given the need to test Unlink on single item list")
	{
		l := NewLList[int]()
		l.AddFirst(1)

		l.Unlink(l.First())

		if l.Size() != 0 {
			t.Errorf("Expected size of linked list to be 0 : %v", l.Size())
		}

		l.AddFirst(10)

		l.Unlink(l.Last())

		if l.Size() != 0 {
			t.Errorf("Expected size of linked list to be 0 : %v", l.Size())
		}

	}

}
func TestListConcurrentFirstLast(t *testing.T) {
	t.Parallel()
	t.Log("Given the need to test concurrent add/remove first/last to the list")
	{
		l := NewLList[int]()

		t.Logf("Testing AddFirst and AddLast with %d concurrent", concurrency)

		wg := sync.WaitGroup{}
		wg.Add(concurrency)
		for i := range concurrency {
			i := i // prevent loop capture - pre Go 1.22 feature
			if i%2 == 0 {
				go func() {
					defer wg.Done()
					l.AddFirst(i)
				}()
			} else {
				go func() {
					defer wg.Done()
					l.AddLast(i)
				}()
			}
		}
		wg.Wait()

		if l.Size() != concurrency {
			t.Errorf("\t Size was expected to be %d : %v", concurrency, l.Size())
		}

		for i := range concurrency {
			if !l.Contains(i) {
				t.Errorf("Expected list to contain value: %v", i)
			}
		}

		t.Logf("Testing AddFirst and AddLast with %d concurrent", concurrency)

		t.Logf("Testing RemoveFirst and RemoveLast with %d concurrent", concurrency)
		wg.Add(concurrency)
		for i := range concurrency {
			i := i // prevent capture
			mod := i % 4
			if mod == 0 {
				go func() {
					defer wg.Done()
					l.RemoveFirst()
				}()
			} else if mod == 1 {
				go func() {
					defer wg.Done()
					l.RemoveLast()
				}()
			} else if mod == 2 {
				go func() {
					defer wg.Done()
					l.AddFirst(i)
				}()
			} else if mod == 3 {
				go func() {
					defer wg.Done()
					l.AddLast(i)
				}()
			}
		}
		wg.Wait()

		if l.Size() != concurrency {
			t.Errorf("\t Size was expected to be 0 : %v", l.Size())
		}

	}
}
func TestListConcurrentBeforeAfter(t *testing.T) {
	t.Parallel()
	t.Log("Given the need to test concurrent add/remove before/after the first item in the list")
	{
		l := NewLList[int]()

		l.AddFirst(100000) // single entry to allow AddBefore / AddAfter

		t.Log("Testing AddBefore and AddAfter concurrently")
		wg := sync.WaitGroup{}
		wg.Add(concurrency)
		for i := range concurrency {
			i := i // prevent capture pre Go 1.22
			if i%2 == 0 {
				go func() {
					defer wg.Done()
					l.AddBefore(l.First(), NewLNode(i))

				}()
			} else {
				go func() {
					defer wg.Done()
					l.AddAfter(l.First(), NewLNode(i))
				}()
			}
		}
		wg.Wait()

		if l.Size() != concurrency+1 {
			t.Errorf("\t Size was expected to be %d : %d", concurrency+1, l.Size())
		}

		t.Logf("Testing ToFirst / ToLast with %d concurrent", concurrency)

		wg.Add(concurrency)
		for i := range concurrency {
			i := i // prevent capture
			if i%2 == 0 {
				go func() {
					defer wg.Done()
					time.Sleep(time.Duration(rand.Intn(concurrency)) * time.Microsecond) // random sleep to increase unpredictability
					l.ToLast(l.First())
				}()
			} else {
				go func() {
					defer wg.Done()
					time.Sleep(time.Duration(rand.Intn(concurrency)) * time.Microsecond)
					l.ToFirst(l.Last())
				}()
			}
		}
		wg.Wait()

		if l.Size() != concurrency+1 {
			t.Errorf("\t Size was expected to be %d : %v", concurrency+1, l.Size())
		}

		for i := range concurrency {
			if !l.Contains(i) {
				t.Errorf("Expected list to contain value: %v", i)
			}
		}
		if !l.Contains(100000) {
			t.Errorf("Expected list to contain value: %v", 100000)
		}

	}
}
