package casset

const maxUint = ^uint(0)

// Len is just hold positive value and made it with IElement.
type Len struct {
	current IElement
}

var _ ILen = &Len{}

// NewLen creates a new Len with the given value.
func NewLen(val uint) *Len {
	return &Len{
		current: &Element{Value: val},
	}
}

// IsZero returns true if the len is zero.
func (l *Len) IsZero() bool {
	if l.current.GetPrevElement() != nil || l.current.GetValue().(uint) != 0 {
		return false
	}

	return true
}

// GetValueCurrent returns the current value of the len not all len value.
func (l *Len) GetValueCurrent() uint {
	return l.current.GetValue().(uint)
}

// GetElement returns the current element.
// Use this method for read-only operations.
func (l *Len) GetElement() IElement {
	return l.current
}

// SubLen subtracts the given len from the current len.
func (l *Len) SubLen(l2 ILen) ILen {
	for l2Element := l2.GetElement(); l2Element != nil; l2Element = l2Element.GetPrevElement() {
		v := l2Element.GetValue().(uint)
		l.Sub(v)

		if l.IsZero() {
			break
		}
	}

	return l
}

// AddLen adds the given len to the current len.
func (l *Len) AddLen(l2 ILen) ILen {
	for l2Element := l2.GetElement(); l2Element != nil; l2Element = l2Element.GetPrevElement() {
		v := l2Element.GetValue().(uint)
		l.Add(v)
	}

	return l
}

// Sub subtracts the given value from the len.
func (l *Len) Sub(val uint) ILen {
	if val >= l.current.GetValue().(uint) {
		if l.current.GetPrevElement() != nil {
			val -= l.current.GetValue().(uint)
			l.current = l.current.Delete()
			l.current.SetValue(l.current.GetValue().(uint) - val)

			return l
		}

		l.current.SetValue(uint(0))

		return l
	}

	l.current.SetValue(l.current.GetValue().(uint) - val)

	return l
}

// Add adds the given value to the len.
func (l *Len) Add(val uint) ILen {
	if maxUint-l.current.GetValue().(uint) >= val {
		l.current.SetValue(l.current.GetValue().(uint) + val)

		return l
	}

	val -= maxUint - l.current.GetValue().(uint)
	l.current.SetValue(maxUint)

	l.current = l.current.Next(val)

	return l
}

// Cmp compares x and y on current element and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (l *Len) Cmp(y uint) int {
	if l.current.GetValue().(uint) > y {
		return 1
	}

	if l.current.GetValue().(uint) < y {
		return -1
	}

	return 0
}
