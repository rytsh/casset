// Package casset help you to create memory on double linked list.
package casset

// Element is an struct of duble-linked list.
type Element struct {
	NextElement IElement
	PrevElement IElement

	// Know about belong Memory
	Memory IMemory
	Value  interface{}
}

// Correction of interface.
var _ IElement = &Element{}

func NewElement(v interface{}) *Element {
	return &Element{
		Value: v,
	}
}

// New generate new element with value and memory.
func (e Element) New(v interface{}, m IMemory) IElement {
	return &Element{
		Value:  v,
		Memory: m,
	}
}

func (e *Element) cleanup() {
	e.Memory = nil
	e.NextElement = nil
	e.PrevElement = nil
	e.Value = nil
}

func (e *Element) GetMemory() IMemory {
	return e.Memory
}

func (e *Element) SetMemory(m IMemory) IElement {
	e.Memory = m

	return e
}

func (e *Element) GetValue() interface{} {
	return e.Value
}

func (e *Element) GetNextElement() IElement {
	if e.NextElement != nil {
		return e.NextElement
	}

	return nil
}

func (e *Element) GetPrevElement() IElement {
	if e.PrevElement != nil {
		return e.PrevElement
	}

	return nil
}

func (e *Element) SetValue(v interface{}) {
	e.Value = v
}

// SetNextElement set next element.
func (e *Element) SetNextElement(element IElement) IElement {
	e.NextElement = element

	return e
}

// SetPrevElement set previous element.
func (e *Element) SetPrevElement(element IElement) IElement {
	e.PrevElement = element

	return e
}

// Delete this element, reconnect prev and next if exist.
// When deleting current element, it will set current element to last or next element.
func (e *Element) Delete() IElement {
	if e.Memory == nil {
		if e.PrevElement != nil {
			e.PrevElement.SetNextElement(e.NextElement)
		}

		if e.NextElement != nil {
			e.NextElement.SetPrevElement(e.PrevElement)
		}

		// set return
		ret := e.NextElement
		if ret == nil {
			ret = e.PrevElement
		}

		// celanup
		e.cleanup()

		return ret
	}

	// just one element
	if e.Memory.GetLen().Cmp(1) == 0 {
		e.Memory.SetFront(nil)
		e.Memory.SetBack(nil)
		e.Memory.SetCurrent(nil)
		e.Memory.GetLen().Set(0)

		e.cleanup()

		return nil
	}

	ret := e.NextElement

	switch e {
	case e.Memory.GetFront():
		if e.Memory.GetCurrent() == e.Memory.GetFront() {
			e.Memory.SetCurrent(e.NextElement)
		}

		e.Memory.SetFront(e.NextElement)
		e.NextElement.SetPrevElement(nil)
	case e.Memory.GetBack():
		if e.Memory.GetCurrent() == e.Memory.GetBack() {
			e.Memory.SetCurrent(e.PrevElement)
		}

		e.Memory.SetBack(e.PrevElement)
		e.PrevElement.SetNextElement(nil)

		// next element is nil so we need to set return element to prev element
		ret = e.PrevElement
	default:
		if e.Memory.GetCurrent() == e {
			e.Memory.SetCurrent(e.NextElement)
		}

		e.NextElement.SetPrevElement(e.PrevElement)
		e.PrevElement.SetNextElement(e.NextElement)
	}

	e.Memory.GetLen().Sub(1)

	e.cleanup()

	if ret == nil {
		return nil
	}

	return ret
}

// Next generate new element with argument and return new element.
func (e *Element) Next(v interface{}) IElement {
	if e.NextElement == nil {
		e.NextElement = &Element{
			Memory:      e.Memory,
			PrevElement: e,
			Value:       v,
		}

		if e.Memory != nil {
			e.Memory.SetBack(e.NextElement)
			e.Memory.GetLen().Add(uint(1))
		}
	}

	return e.NextElement
}

// Prev generate new element with argument and return new element.
func (e *Element) Prev(v interface{}) IElement {
	if e.PrevElement == nil {
		e.PrevElement = &Element{
			Memory:      e.Memory,
			NextElement: e,
			Value:       v,
		}

		if e.Memory != nil {
			e.Memory.SetFront(e.PrevElement)
			e.Memory.GetLen().Add(uint(1))
		}
	}

	return e.PrevElement
}
