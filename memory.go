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
	if e1 == nil && e2 == nil {
		return
	}

	if e1 != nil && e1.GetMemory() != m {
		return
	}

	if e2 != nil && e2.GetMemory() != m {
		return
	}

	front := e1
	back := e2

	var (
		fixFront   bool
		fixBack    bool
		fixCurrent bool
	)

	if e1 == nil {
		front = m.Front
		fixFront = true
	}

	if e2 == nil {
		back = m.Back
		fixBack = true
	}

	frontPrev := front.GetPrevElement()
	// frontNext := front.GetNextElement()
	// backPrev := back.GetPrevElement()
	backNext := back.GetNextElement()

	current := front
	count := NewLen(0)

	deleteElement := func(current IElement) IElement {
		switch current {
		case m.Current:
			fixCurrent = true
		case m.Front:
			fixFront = true
		case m.Back:
			fixBack = true
		}

		current.SetPrevElement(nil)
		current.SetMemory(nil)
		current.SetValue(nil)

		current = current.GetNextElement()

		if current != nil {
			current.SetPrevElement(nil)
		}

		count.Add(1)

		return current
	}

	for current != back {
		current = deleteElement(current)

		// reverse delete
		if current == nil && back != nil {
			current = m.Front
		}
	}

	// delete back
	if current != nil {
		current = deleteElement(current)
	}

	if current == nil || count.IsZero() {
		return
	}

	m.len.SubLen(count)

	if m.len.IsZero() {
		return
	}

	// fix memory front, back, current
	if fixFront {
		m.Front = backNext
		m.Front.SetPrevElement(nil)
	}

	if fixBack {
		m.Back = frontPrev
		m.Back.SetNextElement(nil)
	}

	if fixCurrent {
		m.Current = frontPrev
	}

	if !fixBack && !fixFront {
		frontPrev.SetNextElement(backNext)
		backNext.SetPrevElement(frontPrev)
	}
}
