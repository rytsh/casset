<img src="_assets/casset.png" alt="casset_logo" width="400"/>

[![License](https://img.shields.io/github/license/rytsh/casset?color=red&style=flat-square)](https://raw.githubusercontent.com/rytsh/casset/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/rytsh_casset?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=rytsh_casset)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/rytsh/casset/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/rytsh/casset/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/rytsh/casset?style=flat-square)](https://goreportcard.com/report/github.com/rytsh/casset)
[![Go PKG](https://raw.githubusercontent.com/rakunlabs/.github/main/assets/badges/gopkg.svg)](https://pkg.go.dev/github.com/rytsh/casset)

Casset is double linked endless memory library.

Always generate new space automatically.

Memory hold length, front, back and current location  
Element hold value belong memory address and next, previous elements address.

## Installation

```sh
go get github.com/rytsh/casset
```

## Usage

```go
// create a memory with a first value
// value could be anything
l := casset.NewMemory[any]().GetFront().SetValue("My First Element")
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
```
