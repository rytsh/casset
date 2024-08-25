// Package casset help you to create memory on double linked list.
package casset

// Element is an struct of double-linked list.
type Element[T any] struct {
	NextElement IElement[T]
	PrevElement IElement[T]

	// Know about belong Memory
	Memory IMemory[T]
	Value  T
}

// Correction of interface.
var _ IElement[any] = (*Element[any])(nil)

func NewElement[T any](v T) IElement[T] {
	return &Element[T]{
		Value: v,
	}
}

func (e *Element[T]) Clone() IElement[T] {
	return &Element[T]{
		Value: e.Value,
	}
}

func (e *Element[T]) cleanup() {
	e.Memory = nil
	e.NextElement = nil
	e.PrevElement = nil
}

func (e *Element[T]) GetMemory() IMemory[T] {
	return e.Memory
}

func (e *Element[T]) SetMemory(m IMemory[T]) IElement[T] {
	e.Memory = m

	return e
}

func (e *Element[T]) GetValue() T {
	return e.Value
}

func (e *Element[T]) GetNextElement() IElement[T] {
	if e.NextElement != nil {
		return e.NextElement
	}

	return nil
}

func (e *Element[T]) GetPrevElement() IElement[T] {
	if e.PrevElement != nil {
		return e.PrevElement
	}

	return nil
}

func (e *Element[T]) SetValue(v T) IElement[T] {
	e.Value = v

	return e
}

// SetNextElement set next element.
func (e *Element[T]) SetNextElement(element IElement[T]) IElement[T] {
	e.NextElement = element

	return e
}

// SetPrevElement set previous element.
func (e *Element[T]) SetPrevElement(element IElement[T]) IElement[T] {
	e.PrevElement = element

	return e
}

// Delete this element, reconnect prev and next if exist.
// When deleting current element, it will set current element to last or next element.
func (e *Element[T]) Delete() IElement[T] {
	if e.Memory == nil {
		if e.PrevElement != nil {
			e.PrevElement.SetNextElement(e.NextElement)
		}

		if e.NextElement != nil {
			e.NextElement.SetPrevElement(e.PrevElement)
		}

		// set return
		ret := e.NextElement

		e.cleanup()

		return ret
	}

	// just one element
	if e.Memory.GetLen().Cmp(1) == 0 {
		e.Memory.SetFront(nil)
		e.Memory.SetBack(nil)
		e.Memory.GetLen().Sub(1)

		e.cleanup()

		return nil
	}

	ret := e.NextElement

	switch e {
	case e.Memory.GetFront():
		e.Memory.SetFront(e.NextElement)
		e.NextElement.SetPrevElement(nil)
	case e.Memory.GetBack():
		e.Memory.SetBack(e.PrevElement)
		e.PrevElement.SetNextElement(nil)
	default:
		e.NextElement.SetPrevElement(e.PrevElement)
		e.PrevElement.SetNextElement(e.NextElement)
	}

	e.Memory.GetLen().Sub(1)

	e.cleanup()

	return ret
}

// Next generate new element with argument and return new element.
func (e *Element[T]) Next(v T) IElement[T] {
	if e.NextElement == nil {
		e.NextElement = &Element[T]{
			Memory:      e.Memory,
			PrevElement: e,
			Value:       v,
		}

		if e.Memory != nil {
			e.Memory.SetBack(e.NextElement)
			e.Memory.GetLen().Add(1)
		}
	}

	return e.NextElement
}

// Prev generate new element with argument and return new element.
func (e *Element[T]) Prev(v T) IElement[T] {
	if e.PrevElement == nil {
		e.PrevElement = &Element[T]{
			Memory:      e.Memory,
			NextElement: e,
			Value:       v,
		}

		if e.Memory != nil {
			e.Memory.SetFront(e.PrevElement)
			e.Memory.GetLen().Add(1)
		}
	}

	return e.PrevElement
}
