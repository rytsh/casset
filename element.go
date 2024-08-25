// Package casset help you to create memory on double linked list.
package casset

import "math/big"

// Element is an struct of double-linked list.
type Element[T any] struct {
	nextElement IElement[T]
	prevElement IElement[T]

	// Know about belong memory
	memory IMemory[T]
	value  T
}

// Correction of interface.
var _ IElement[any] = (*Element[any])(nil)

func NewElement[T any](v T) IElement[T] {
	return &Element[T]{
		value: v,
	}
}

func (e *Element[T]) cleanup() {
	e.memory = nil
	e.nextElement = nil
	e.prevElement = nil
	e.value = *new(T)
}

func (e *Element[T]) GetMemory() IMemory[T] {
	return e.memory
}

func (e *Element[T]) SetMemory(m IMemory[T]) IElement[T] {
	e.memory = m

	return e
}

func (e *Element[T]) GetValue() T {
	return e.value
}

func (e *Element[T]) GetNextElement() IElement[T] {
	if e.nextElement != nil {
		return e.nextElement
	}

	return nil
}

func (e *Element[T]) GetPrevElement() IElement[T] {
	if e.prevElement != nil {
		return e.prevElement
	}

	return nil
}

func (e *Element[T]) SetValue(v T) IElement[T] {
	e.value = v

	return e
}

// SetNextElement set next element.
func (e *Element[T]) SetNextElement(element IElement[T]) IElement[T] {
	e.nextElement = element

	return e
}

// SetPrevElement set previous element.
func (e *Element[T]) SetPrevElement(element IElement[T]) IElement[T] {
	e.prevElement = element

	return e
}

// Delete this element, reconnect prev and next if exist.
// When deleting current element, it will set current element to last or next element.
func (e *Element[T]) Delete() IElement[T] {
	if e.memory == nil {
		if e.prevElement != nil {
			e.prevElement.SetNextElement(e.nextElement)
		}

		if e.nextElement != nil {
			e.nextElement.SetPrevElement(e.prevElement)
		}

		// set return
		ret := e.nextElement

		e.cleanup()

		return ret
	}

	// just one element
	if e.memory.GetLen().Cmp(big.NewInt(1)) == 0 {
		e.memory.Clear()

		e.cleanup()

		return nil
	}

	ret := e.nextElement

	switch e {
	case e.memory.GetFront():
		e.memory.SetFront(e.nextElement)
		e.nextElement.SetPrevElement(nil)
	case e.memory.GetBack():
		e.memory.SetBack(e.prevElement)
		e.prevElement.SetNextElement(nil)
	default:
		e.nextElement.SetPrevElement(e.prevElement)
		e.prevElement.SetNextElement(e.nextElement)
	}

	e.memory.GetLen().Set(func(i *big.Int) *big.Int {
		return i.Sub(i, big.NewInt(1))
	})

	e.cleanup()

	return ret
}

// Next generate new element with argument and return new element.
func (e *Element[T]) Next(v T) IElement[T] {
	if e.nextElement == nil {
		e.nextElement = &Element[T]{
			memory:      e.memory,
			prevElement: e,
			value:       v,
		}

		if e.memory != nil {
			e.memory.SetBack(e.nextElement)
			e.memory.GetLen().Set(func(i *big.Int) *big.Int {
				return i.Add(i, big.NewInt(1))
			})
		}
	}

	return e.nextElement
}

// Prev generate new element with argument and return new element.
func (e *Element[T]) Prev(v T) IElement[T] {
	if e.prevElement == nil {
		e.prevElement = &Element[T]{
			memory:      e.memory,
			nextElement: e,
			value:       v,
		}

		if e.memory != nil {
			e.memory.SetFront(e.prevElement)
			e.memory.GetLen().Set(func(i *big.Int) *big.Int {
				return i.Add(i, big.NewInt(1))
			})
		}
	}

	return e.prevElement
}
