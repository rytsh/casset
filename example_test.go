package casset_test

import (
	"fmt"

	"github.com/rytsh/casset"
)

func Example() {
	// create a memory with a first value
	// value could be anything
	l := casset.NewMemory("My First Element").Current
	l.Next("second element").Next(3.14).Next(struct{ v string }{v: "4th element"})

	for e := l.GetMemory().GetFront(); e != nil; e = e.GetNextElement() {
		fmt.Println(e.GetValue())
	}

	// Output:
	// My First Element
	// second element
	// 3.14
	// {4th element}
}
