package casset

import (
	"iter"
)

// Memory is main struct of linked list.
type Memory[T any] struct {
	front IElement[T]
	back  IElement[T]

	len uint64
}

// NewMemory return new empty memory. Before use, you must call Init method.
func NewMemory[T any]() IMemory[T] {
	return (&Memory[T]{}).Clear()
}

func (m *Memory[T]) Clear() IMemory[T] {
	element := new(Element[T])
	element.SetMemory(m)

	m.front = element
	m.back = element

	m.len = 1

	return m
}

func (m *Memory[T]) Range() iter.Seq[IElement[T]] {
	return func(yield func(IElement[T]) bool) {
		for e := m.GetFront(); e != nil; e = e.GetNextElement() {
			if !yield(e) {
				return
			}
		}
	}
}

func (m *Memory[T]) Len() uint64 {
	return m.len
}

func (m *Memory[T]) IncLen() {
	m.len++
}

func (m *Memory[T]) DecLen() {
	if m.len > 0 {
		m.len--
	}
}

func (m *Memory[T]) GetFront() IElement[T] {
	return m.front
}

func (m *Memory[T]) SetFront(e IElement[T]) {
	m.front = e
}

func (m *Memory[T]) GetBack() IElement[T] {
	return m.back
}

func (m *Memory[T]) SetBack(e IElement[T]) {
	m.back = e
}

// ToSliceElements converts a range of elements to a slice-backed element block.
// The original elements are removed and replaced with a SliceElement.
// Returns the first SliceElement of the converted range.
func (m *Memory[T]) ToSliceElements(e1, e2 IElement[T]) IElement[T] {
	if e1 != nil && e1.GetMemory() != m {
		return nil
	}

	if e2 != nil && e2.GetMemory() != m {
		return nil
	}

	front := e1
	back := e2

	if e1 == nil {
		front = m.front
	}

	if e2 == nil {
		back = m.back
	}

	// Collect values and count elements
	var values []T
	current := front
	for {
		values = append(values, current.GetValue())
		if current == back {
			break
		}
		next := current.GetNextElement()
		if next == nil {
			break
		}
		current = next
	}

	if len(values) == 0 {
		return nil
	}

	// Get external links (elements outside the range)
	prevExt := front.GetPrevElement()
	nextExt := back.GetNextElement()

	// Create the slice element
	sliceElem := NewSliceElement(values, prevExt, nextExt, m).(*SliceElement[T])

	// Get the last slice element for linking
	lastSliceElem := sliceElem.Last().(*SliceElement[T])

	// Update external links to point to slice element
	if prevExt != nil {
		prevExt.SetNextElement(sliceElem)
	} else {
		m.front = sliceElem
	}

	if nextExt != nil {
		nextExt.SetPrevElement(lastSliceElem)
	} else {
		m.back = lastSliceElem
	}

	// Clean up original elements (don't use Delete() as it modifies length)
	// Just disconnect them
	current = front
	for {
		next := current.GetNextElement()

		// Disconnect without affecting memory length
		current.SetMemory(nil)
		current.SetNextElement(nil)
		current.SetPrevElement(nil)

		if current == back {
			break
		}
		if next == nil {
			break
		}
		current = next
	}

	return sliceElem
}

// Remove remove range of elements.
// If elements not inside of memory, nothing change.
func (m *Memory[T]) RemoveRange(e1, e2 IElement[T]) {
	if e1 != nil && e1.GetMemory() != m {
		return
	}

	if e2 != nil && e2.GetMemory() != m {
		return
	}

	front := e1
	back := e2

	if e1 == nil {
		front = m.front
	}

	if e2 == nil {
		back = m.back
	}

	current := front

	for current != back {
		current = current.Delete()

		if current == nil {
			current = m.GetFront()
		}
	}

	// delete back
	if current == back && current != nil {
		current.Delete()
	}
}
