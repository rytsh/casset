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

func (l *Len) Sub(n int64) ILen {
	l.value.Sub(l.value, big.NewInt(n))

	return l
}

func (l *Len) Add(n int64) ILen {
	l.value.Add(l.value, big.NewInt(n))

	return l
}

func (l *Len) Cmp(y int64) int {
	return l.value.Cmp(big.NewInt(y))
}
