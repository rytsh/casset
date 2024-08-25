package casset_test

import (
	"fmt"

	"github.com/rytsh/casset"
)

func Example() {
	// create a memory with a first value
	// value could be anything
	l := casset.NewMemory[any]().Init(casset.NewElement[any]("My First Element")).GetFront()
	l.Next("second element").Next(3.14).Next(struct{ v string }{v: "4th element"})

	// for e := l.GetMemory().GetFront(); e != nil; e = e.GetNextElement() {
	// 	fmt.Println(e.GetValue())
	// }

	for e := range l.GetMemory().Range() {
		fmt.Println(e.GetValue())
	}

	// Output:
	// My First Element
	// second element
	// 3.14
	// {4th element}
}

func ExampleElement() {
	e := casset.NewElement(1234)
	e.Next(5678).Next(91011)

	for v := e; v != nil; v = v.GetNextElement() {
		fmt.Println(v.GetValue())
	}

	// Output:
	// 1234
	// 5678
	// 91011
}
