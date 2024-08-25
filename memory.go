package casset

import "iter"

// Memory is main struct of linked list.
type Memory[T any] struct {
	Front    IElement[T]
	Back     IElement[T]
	Elements map[string]IElement[T]

	len ILen
}

// NewMemory return new empty memory. Before use, you must call Init method.
func NewMemory[T any]() IMemory[T] {
	return &Memory[T]{
		len: NewLen(),
	}
}

func (m *Memory[T]) Hold(f func(h map[string]IElement[T])) {
	if m.Elements == nil {
		m.Elements = make(map[string]IElement[T])
	}

	f(m.Elements)
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
	return m.Front
}

func (m *Memory[T]) SetFront(e IElement[T]) {
	m.Front = e
}

func (m *Memory[T]) GetBack() IElement[T] {
	return m.Back
}

func (m *Memory[T]) SetBack(e IElement[T]) {
	m.Back = e
}

func (m *Memory[T]) Init(e IElement[T]) IMemory[T] {
	element := e.Clone().SetMemory(m)

	m.Front = element
	m.Back = element

	m.Elements = make(map[string]IElement[T])

	m.len = NewLen().Add(1)

	return m
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
		front = m.Front
	}

	if e2 == nil {
		back = m.Back
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
