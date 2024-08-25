package casset

import "math/big"

type Len struct {
	value *big.Int
}

var _ ILen = &Len{}

func NewLen() *Len {
	return &Len{
		value: big.NewInt(0),
	}
}

func (l *Len) String() string {
	return l.value.String()
}

func (l *Len) Value() big.Int {
	return *l.value
}

func (l *Len) Set(f func(*big.Int) *big.Int) ILen {
	l.value = f(l.value)

	return l
}

func (l *Len) Cmp(y *big.Int) int {
	return l.value.Cmp(y)
}
