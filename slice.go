package casset

// SliceElement is a slice-backed implementation of IElement.
// It provides O(1) navigation within the slice instead of pointer chasing.
type SliceElement[T any] struct {
	values *[]T // Pointer to shared slice (shared among all SliceElements in same block)
	index  int  // Current position in the slice
	memory IMemory[T]

	// External links to regular linked list elements at boundaries
	prevExternal IElement[T] // Element before this slice block (when index == 0)
	nextExternal IElement[T] // Element after this slice block (when index == len-1)
}

// Correction of interface.
var _ IElement[any] = (*SliceElement[any])(nil)

// NewSliceElement creates a new slice-backed element block from a slice of values.
// The prevExt and nextExt are the external elements that link to the regular linked list.
func NewSliceElement[T any](values []T, prevExt, nextExt IElement[T], mem IMemory[T]) IElement[T] {
	if len(values) == 0 {
		return nil
	}

	// Create a copy of the slice to own it
	owned := make([]T, len(values))
	copy(owned, values)

	return &SliceElement[T]{
		values:       &owned,
		index:        0,
		memory:       mem,
		prevExternal: prevExt,
		nextExternal: nextExt,
	}
}

// newSliceElementAt creates a SliceElement pointing to a specific index in an existing slice.
func newSliceElementAt[T any](values *[]T, index int, prevExt, nextExt IElement[T], mem IMemory[T]) *SliceElement[T] {
	return &SliceElement[T]{
		values:       values,
		index:        index,
		memory:       mem,
		prevExternal: prevExt,
		nextExternal: nextExt,
	}
}

func (s *SliceElement[T]) cleanup() {
	s.values = nil
	s.memory = nil
	s.prevExternal = nil
	s.nextExternal = nil
	s.index = 0
}

func (s *SliceElement[T]) GetMemory() IMemory[T] {
	return s.memory
}

func (s *SliceElement[T]) SetMemory(m IMemory[T]) IElement[T] {
	s.memory = m
	return s
}

func (s *SliceElement[T]) GetValue() T {
	if s.values == nil || s.index < 0 || s.index >= len(*s.values) {
		return *new(T)
	}
	return (*s.values)[s.index]
}

func (s *SliceElement[T]) SetValue(v T) IElement[T] {
	if s.values != nil && s.index >= 0 && s.index < len(*s.values) {
		(*s.values)[s.index] = v
	}
	return s
}

// GetNextElement returns the next element - O(1) within slice.
func (s *SliceElement[T]) GetNextElement() IElement[T] {
	if s.values == nil {
		return nil
	}

	// If we can move within the slice, do so
	if s.index < len(*s.values)-1 {
		return newSliceElementAt(s.values, s.index+1, s.prevExternal, s.nextExternal, s.memory)
	}

	// At the end of slice, return external next
	return s.nextExternal
}

// GetPrevElement returns the previous element - O(1) within slice.
func (s *SliceElement[T]) GetPrevElement() IElement[T] {
	if s.values == nil {
		return nil
	}

	// If we can move within the slice, do so
	if s.index > 0 {
		return newSliceElementAt(s.values, s.index-1, s.prevExternal, s.nextExternal, s.memory)
	}

	// At the start of slice, return external prev
	return s.prevExternal
}

// SetNextElement sets the next external element (only affects boundary).
func (s *SliceElement[T]) SetNextElement(element IElement[T]) IElement[T] {
	// Only update external link if we're at the last position
	if s.values != nil && s.index == len(*s.values)-1 {
		s.nextExternal = element
	}
	return s
}

// SetPrevElement sets the previous external element (only affects boundary).
func (s *SliceElement[T]) SetPrevElement(element IElement[T]) IElement[T] {
	// Only update external link if we're at the first position
	if s.index == 0 {
		s.prevExternal = element
	}
	return s
}

// Delete removes the current element from the slice and shifts remaining elements.
func (s *SliceElement[T]) Delete() IElement[T] {
	if s.values == nil || len(*s.values) == 0 {
		return nil
	}

	// If this is the only element in the slice
	if len(*s.values) == 1 {
		// Reconnect external elements
		if s.prevExternal != nil {
			s.prevExternal.SetNextElement(s.nextExternal)
		}
		if s.nextExternal != nil {
			s.nextExternal.SetPrevElement(s.prevExternal)
		}

		// Update memory front/back if needed
		if s.memory != nil {
			if s.memory.GetFront() == s {
				s.memory.SetFront(s.nextExternal)
			}
			if s.memory.GetBack() == s {
				s.memory.SetBack(s.prevExternal)
			}
			s.memory.DecLen()
		}

		ret := s.nextExternal
		s.cleanup()
		return ret
	}

	// Remove element at current index by shifting
	vals := *s.values
	copy(vals[s.index:], vals[s.index+1:])
	vals = vals[:len(vals)-1]
	*s.values = vals

	// Update memory length
	if s.memory != nil {
		s.memory.DecLen()
	}

	// Return next element
	if s.index < len(*s.values) {
		// Still valid index after shift, return element at same index (which is now the "next" value)
		return newSliceElementAt(s.values, s.index, s.prevExternal, s.nextExternal, s.memory)
	}

	// We were at the end, return external next or back up one
	if s.nextExternal != nil {
		return s.nextExternal
	}

	// Stay at the new last element
	if len(*s.values) > 0 {
		return newSliceElementAt(s.values, len(*s.values)-1, s.prevExternal, s.nextExternal, s.memory)
	}

	return nil
}

// Next returns the next element if it exists, or creates a new one at the end.
func (s *SliceElement[T]) Next(v T) IElement[T] {
	if s.values == nil {
		return nil
	}

	// If not at the end, just return next
	if s.index < len(*s.values)-1 {
		return newSliceElementAt(s.values, s.index+1, s.prevExternal, s.nextExternal, s.memory)
	}

	// At the end - if there's an external next, return it
	if s.nextExternal != nil {
		return s.nextExternal
	}

	// Append to slice
	*s.values = append(*s.values, v)

	if s.memory != nil {
		s.memory.IncLen()
	}

	return newSliceElementAt(s.values, len(*s.values)-1, s.prevExternal, s.nextExternal, s.memory)
}

// Prev returns the previous element if it exists, or creates a new one at the start.
func (s *SliceElement[T]) Prev(v T) IElement[T] {
	if s.values == nil {
		return nil
	}

	// If not at the start, just return prev
	if s.index > 0 {
		return newSliceElementAt(s.values, s.index-1, s.prevExternal, s.nextExternal, s.memory)
	}

	// At the start - if there's an external prev, return it
	if s.prevExternal != nil {
		return s.prevExternal
	}

	// Prepend to slice (insert at beginning)
	newVals := make([]T, len(*s.values)+1)
	newVals[0] = v
	copy(newVals[1:], *s.values)
	*s.values = newVals

	// Adjust our index since we prepended
	s.index++

	if s.memory != nil {
		s.memory.IncLen()
	}

	return newSliceElementAt(s.values, 0, s.prevExternal, s.nextExternal, s.memory)
}

// Len returns the length of the underlying slice.
func (s *SliceElement[T]) Len() int {
	if s.values == nil {
		return 0
	}
	return len(*s.values)
}

// Index returns the current index within the slice.
func (s *SliceElement[T]) Index() int {
	return s.index
}

// At returns a SliceElement at the specified index - O(1) random access.
func (s *SliceElement[T]) At(index int) IElement[T] {
	if s.values == nil || index < 0 || index >= len(*s.values) {
		return nil
	}
	return newSliceElementAt(s.values, index, s.prevExternal, s.nextExternal, s.memory)
}

// First returns the first element in this slice block.
func (s *SliceElement[T]) First() IElement[T] {
	return s.At(0)
}

// Last returns the last element in this slice block.
func (s *SliceElement[T]) Last() IElement[T] {
	if s.values == nil {
		return nil
	}
	return s.At(len(*s.values) - 1)
}

// Values returns a copy of the underlying slice.
func (s *SliceElement[T]) Values() []T {
	if s.values == nil {
		return nil
	}
	result := make([]T, len(*s.values))
	copy(result, *s.values)
	return result
}
