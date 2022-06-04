package casset

type IElement interface {
	// New generate new element with value and memory.
	New(interface{}, IMemory) IElement
	// GetMemory return memory of this element.
	GetMemory() IMemory
	// SetMemory set memory of this element.
	SetMemory(m IMemory) IElement
	// GetValue return value of this element.
	GetValue() interface{}
	// GetNextElement return next element of this element.
	GetNextElement() IElement
	// GetPrevElement return prev element of this element.
	GetPrevElement() IElement
	// SetValue set value of this element.
	SetValue(v interface{})
	// SetNextElement set next element.
	SetNextElement(IElement) IElement
	// SetPrevElement set prev element.
	SetPrevElement(IElement) IElement
	// Delete this element, reconnect prev and next if exist.
	Delete() IElement
	// Next generate new element with argument and return new element.
	Next(interface{}) IElement
	// Prev generate new element with argument and return new element.
	Prev(interface{}) IElement
}

type IMemory interface {
	Init(e IElement) IMemory
	// Remove remove range of elements.
	Remove(e1, e2 IElement)
	GetLen() ILen
	GetFront() IElement
	SetFront(e IElement)
	GetBack() IElement
	SetBack(e IElement)
	GetCurrent() IElement
	SetCurrent(e IElement)
}

type ILen interface {
	IsZero() bool
	GetValueCurrent() uint
	GetElement() IElement
	SubLen(ILen) ILen
	AddLen(ILen) ILen
	Sub(uint) ILen
	Add(uint) ILen
	Cmp(uint) int
}
