package casset

import (
	"iter"
	"math/big"
)

// Memory is main struct of linked list.
type Memory[T any] struct {
	front IElement[T]
	back  IElement[T]

	len ILen
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

	m.len = NewLen().Set(func(_ *big.Int) *big.Int {
		return big.NewInt(1)
	})

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

func (m *Memory[T]) GetLen() ILen {
	return m.len
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
