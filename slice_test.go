package casset

import (
	"testing"
)

func TestSliceElement_BasicNavigation(t *testing.T) {
	// Create a slice element with values 1, 2, 3, 4, 5
	values := []int{1, 2, 3, 4, 5}
	elem := NewSliceElement(values, nil, nil, nil)

	// Test GetValue at first position
	if got := elem.GetValue(); got != 1 {
		t.Errorf("GetValue() = %v, want 1", got)
	}

	// Navigate forward
	next := elem.GetNextElement()
	if next == nil {
		t.Fatal("GetNextElement() returned nil")
	}
	if got := next.GetValue(); got != 2 {
		t.Errorf("GetNextElement().GetValue() = %v, want 2", got)
	}

	// Navigate to end
	current := elem
	for i := 0; i < 4; i++ {
		current = current.GetNextElement()
	}
	if got := current.GetValue(); got != 5 {
		t.Errorf("After 4 nexts, GetValue() = %v, want 5", got)
	}

	// Next at end should return nil (no external link)
	if got := current.GetNextElement(); got != nil {
		t.Errorf("GetNextElement() at end = %v, want nil", got)
	}

	// Navigate backward
	prev := current.GetPrevElement()
	if prev == nil {
		t.Fatal("GetPrevElement() returned nil")
	}
	if got := prev.GetValue(); got != 4 {
		t.Errorf("GetPrevElement().GetValue() = %v, want 4", got)
	}
}

func TestSliceElement_SetValue(t *testing.T) {
	values := []int{1, 2, 3}
	elem := NewSliceElement(values, nil, nil, nil)

	// Set value at first position
	elem.SetValue(100)
	if got := elem.GetValue(); got != 100 {
		t.Errorf("After SetValue(100), GetValue() = %v, want 100", got)
	}

	// Navigate and set
	next := elem.GetNextElement()
	next.SetValue(200)
	if got := next.GetValue(); got != 200 {
		t.Errorf("After SetValue(200), GetValue() = %v, want 200", got)
	}

	// Verify the slice is shared - original elem's next should also see the change
	if got := elem.GetNextElement().GetValue(); got != 200 {
		t.Errorf("Shared slice: elem.GetNextElement().GetValue() = %v, want 200", got)
	}
}

func TestSliceElement_At(t *testing.T) {
	values := []int{10, 20, 30, 40, 50}
	elem := NewSliceElement(values, nil, nil, nil).(*SliceElement[int])

	// Test At() for O(1) random access
	if got := elem.At(2).GetValue(); got != 30 {
		t.Errorf("At(2).GetValue() = %v, want 30", got)
	}

	if got := elem.At(4).GetValue(); got != 50 {
		t.Errorf("At(4).GetValue() = %v, want 50", got)
	}

	// Out of bounds
	if got := elem.At(-1); got != nil {
		t.Errorf("At(-1) = %v, want nil", got)
	}
	if got := elem.At(5); got != nil {
		t.Errorf("At(5) = %v, want nil", got)
	}
}

func TestSliceElement_Delete(t *testing.T) {
	mem := NewMemory[int]()
	mem.GetFront().SetValue(0)
	mem.GetFront().Next(1).Next(2).Next(3).Next(4)

	// Memory has: 0 -> 1 -> 2 -> 3 -> 4
	// Length should be 5
	if mem.Len() != 5 {
		t.Errorf("Initial len = %v, want 5", mem.Len())
	}

	// Convert middle elements (1, 2, 3) to slice
	start := mem.GetFront().GetNextElement()       // element with value 1
	end := start.GetNextElement().GetNextElement() // element with value 3

	sliceElem := mem.ToSliceElements(start, end)
	if sliceElem == nil {
		t.Fatal("ToSliceElements returned nil")
	}

	// Verify slice has correct values
	se := sliceElem.(*SliceElement[int])
	if se.Len() != 3 {
		t.Errorf("SliceElement len = %v, want 3", se.Len())
	}

	// Delete middle element (value 2) from slice
	middle := se.At(1) // value 2
	next := middle.Delete()

	// After delete, slice should have [1, 3]
	if se.Len() != 2 {
		t.Errorf("After delete, SliceElement len = %v, want 2", se.Len())
	}

	// Next should point to value 3 (which shifted to index 1)
	if next != nil && next.GetValue() != 3 {
		t.Errorf("After delete, next.GetValue() = %v, want 3", next.GetValue())
	}

	// Memory length should decrease
	if mem.Len() != 4 {
		t.Errorf("After delete, memory len = %v, want 4", mem.Len())
	}
}

func TestSliceElement_DeleteAll(t *testing.T) {
	values := []int{1}
	elem := NewSliceElement(values, nil, nil, nil)

	// Delete the only element
	next := elem.Delete()
	if next != nil {
		t.Errorf("After deleting only element, next = %v, want nil", next)
	}
}

func TestSliceElement_ExternalLinks(t *testing.T) {
	mem := NewMemory[int]()
	mem.GetFront().SetValue(0)
	mem.GetFront().Next(1).Next(2).Next(3).Next(4).Next(5)

	// Memory has: 0 -> 1 -> 2 -> 3 -> 4 -> 5
	// Convert 2, 3, 4 to slice
	elem2 := mem.GetFront().GetNextElement().GetNextElement()
	elem4 := elem2.GetNextElement().GetNextElement()

	sliceElem := mem.ToSliceElements(elem2, elem4)
	if sliceElem == nil {
		t.Fatal("ToSliceElements returned nil")
	}

	// Verify front is still 0
	if mem.GetFront().GetValue() != 0 {
		t.Errorf("Front value = %v, want 0", mem.GetFront().GetValue())
	}

	// Navigate from front through slice to back
	// 0 -> 1 -> [slice: 2,3,4] -> 5
	current := mem.GetFront()
	expected := []int{0, 1, 2, 3, 4, 5}
	for i, exp := range expected {
		if current == nil {
			t.Fatalf("current is nil at index %d", i)
		}
		if current.GetValue() != exp {
			t.Errorf("At index %d, value = %v, want %v", i, current.GetValue(), exp)
		}
		current = current.GetNextElement()
	}

	// Navigate backward
	current = mem.GetBack()
	for i := len(expected) - 1; i >= 0; i-- {
		if current == nil {
			t.Fatalf("current is nil at reverse index %d", i)
		}
		if current.GetValue() != expected[i] {
			t.Errorf("At reverse index %d, value = %v, want %v", i, current.GetValue(), expected[i])
		}
		current = current.GetPrevElement()
	}
}

func TestSliceElement_Next(t *testing.T) {
	values := []int{1, 2}
	elem := NewSliceElement(values, nil, nil, nil).(*SliceElement[int])

	// Go to last element
	last := elem.At(1)

	// Call Next to append
	newElem := last.Next(3)
	if newElem.GetValue() != 3 {
		t.Errorf("Next(3).GetValue() = %v, want 3", newElem.GetValue())
	}

	// Verify slice grew
	if elem.Len() != 3 {
		t.Errorf("After Next, len = %v, want 3", elem.Len())
	}
}

func TestSliceElement_Prev(t *testing.T) {
	values := []int{2, 3}
	elem := NewSliceElement(values, nil, nil, nil).(*SliceElement[int])

	// Call Prev to prepend
	newElem := elem.Prev(1)
	if newElem.GetValue() != 1 {
		t.Errorf("Prev(1).GetValue() = %v, want 1", newElem.GetValue())
	}

	// Verify slice grew
	if elem.Len() != 3 {
		t.Errorf("After Prev, len = %v, want 3", elem.Len())
	}

	// Verify order: 1, 2, 3
	se := newElem.(*SliceElement[int])
	vals := se.Values()
	expected := []int{1, 2, 3}
	for i, exp := range expected {
		if vals[i] != exp {
			t.Errorf("Values()[%d] = %v, want %v", i, vals[i], exp)
		}
	}
}

func TestSliceElement_Values(t *testing.T) {
	values := []int{1, 2, 3}
	elem := NewSliceElement(values, nil, nil, nil).(*SliceElement[int])

	// Get copy of values
	got := elem.Values()

	// Verify it's a copy
	got[0] = 100
	if elem.GetValue() == 100 {
		t.Error("Values() should return a copy, not the original slice")
	}
}

func TestMemory_ToSliceElements_Nil(t *testing.T) {
	mem := NewMemory[int]()
	mem.GetFront().SetValue(1)
	mem.GetFront().Next(2).Next(3)

	// Convert all to slice (nil, nil means front to back)
	sliceElem := mem.ToSliceElements(nil, nil)
	if sliceElem == nil {
		t.Fatal("ToSliceElements(nil, nil) returned nil")
	}

	se := sliceElem.(*SliceElement[int])
	if se.Len() != 3 {
		t.Errorf("SliceElement len = %v, want 3", se.Len())
	}

	// Verify all values
	expected := []int{1, 2, 3}
	for i, exp := range expected {
		elem := se.At(i)
		if elem.GetValue() != exp {
			t.Errorf("At(%d).GetValue() = %v, want %v", i, elem.GetValue(), exp)
		}
	}
}

func TestMemory_ToSliceElements_WrongMemory(t *testing.T) {
	mem1 := NewMemory[int]()
	mem1.GetFront().SetValue(1)

	mem2 := NewMemory[int]()
	mem2.GetFront().SetValue(2)

	// Try to convert element from different memory
	result := mem1.ToSliceElements(mem2.GetFront(), nil)
	if result != nil {
		t.Error("ToSliceElements with element from different memory should return nil")
	}
}

func TestSliceElement_Performance(t *testing.T) {
	// This test verifies that SliceElement navigation is fast
	// by doing many next/prev operations

	size := 10000
	values := make([]int, size)
	for i := 0; i < size; i++ {
		values[i] = i
	}

	elem := NewSliceElement(values, nil, nil, nil).(*SliceElement[int])

	// Navigate forward and backward many times - should be O(1) each
	current := elem.At(size / 2)
	for i := 0; i < 1000; i++ {
		current = current.GetNextElement()
		current = current.GetPrevElement()
	}

	// Verify we're still at the same position
	if current.GetValue() != size/2 {
		t.Errorf("After navigation, value = %v, want %v", current.GetValue(), size/2)
	}

	// Test random access - should be O(1)
	for i := 0; i < 1000; i++ {
		idx := i % size
		val := elem.At(idx).GetValue()
		if val != idx {
			t.Errorf("At(%d).GetValue() = %v, want %v", idx, val, idx)
		}
	}
}

func TestSliceElement_FirstLast(t *testing.T) {
	values := []int{10, 20, 30, 40}
	elem := NewSliceElement(values, nil, nil, nil).(*SliceElement[int])

	if elem.First().GetValue() != 10 {
		t.Errorf("First().GetValue() = %v, want 10", elem.First().GetValue())
	}

	if elem.Last().GetValue() != 40 {
		t.Errorf("Last().GetValue() = %v, want 40", elem.Last().GetValue())
	}
}

func TestSliceElement_Index(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	elem := NewSliceElement(values, nil, nil, nil).(*SliceElement[int])

	if elem.Index() != 0 {
		t.Errorf("Initial Index() = %v, want 0", elem.Index())
	}

	next := elem.GetNextElement().(*SliceElement[int])
	if next.Index() != 1 {
		t.Errorf("After GetNextElement(), Index() = %v, want 1", next.Index())
	}

	at3 := elem.At(3).(*SliceElement[int])
	if at3.Index() != 3 {
		t.Errorf("At(3).Index() = %v, want 3", at3.Index())
	}
}
