package casset

// Memory is main struct of linked list.
type Memory struct {
	Front   IElement
	Back    IElement
	Current IElement
	len     ILen
}

func NewMemory(v interface{}) *Memory {
	m := new(Memory).Init(NewElement(v))
	return m.(*Memory)
}

func (m *Memory) GetLen() ILen {
	return m.len
}

func (m *Memory) GetFront() IElement {
	return m.Front
}

func (m *Memory) SetFront(e IElement) {
	m.Front = e
}

func (m *Memory) GetBack() IElement {
	return m.Back
}

func (m *Memory) SetBack(e IElement) {
	m.Back = e
}

func (m *Memory) GetCurrent() IElement {
	return m.Current
}

func (m *Memory) SetCurrent(e IElement) {
	m.Current = e
}

func (m *Memory) Init(e IElement) IMemory {
	element := e.New(e.GetValue(), m)

	m.Front = element
	m.Back = element
	m.Current = element

	m.len = NewLen(1)

	return m
}

// Remove remove range of elements.
// If elements not inside of memory, nothing change.
func (m *Memory) Remove(e1, e2 IElement) {
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

	for current != back && current != nil {
		current = current.Delete()
	}

	// delete back
	if current == back && current != nil {
		current.Delete()
	}
}
