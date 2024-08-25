package casset

import (
	"iter"
	"math/big"
)

type IElement[T any] interface {
	// GetMemory return memory of this element.
	GetMemory() IMemory[T]
	// SetMemory set memory of this element.
	SetMemory(m IMemory[T]) IElement[T]
	// GetValue return value of this element.
	GetValue() T
	// GetNextElement return next element of this element.
	GetNextElement() IElement[T]
	// GetPrevElement return prev element of this element.
	GetPrevElement() IElement[T]
	// SetValue set value of this element.
	SetValue(v T) IElement[T]
	// SetNextElement set next element.
	SetNextElement(IElement[T]) IElement[T]
	// SetPrevElement set prev element.
	SetPrevElement(IElement[T]) IElement[T]
	// Delete this element, reconnect prev and next if exist.
	Delete() IElement[T]
	// Next return next element if exist or generate new element with argument and return new element.
	Next(T) IElement[T]
	// Prev return prev element if exist or generate new element with argument and return new element.
	Prev(T) IElement[T]
}

type IMemory[T any] interface {
	// Clear remove all elements.
	Clear() IMemory[T]
	// RemoveRange remove range of elements including e1 and e2.
	// If e1 is nil, start from front.
	// If e2 is nil, end at back.
	// Both e1 and e2 are nil, remove all.
	RemoveRange(e1, e2 IElement[T])
	GetLen() ILen
	GetFront() IElement[T]
	SetFront(e IElement[T])
	GetBack() IElement[T]
	SetBack(e IElement[T])
	Range() iter.Seq[IElement[T]]
}

type ILen interface {
	Value() big.Int
	Set(func(*big.Int) *big.Int) ILen
	// Cmp compares x and y on current element and returns:
	//
	//   -1 if x <  y
	//    0 if x == y
	//   +1 if x >  y
	Cmp(y *big.Int) int
}
