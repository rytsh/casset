<img src="_assets/casset.png" alt="casset_logo" width="400"/>

Casset is double linked endless memory library.

Always generate new space automatically.

Memory hold length, front, back and current location  
Element hold value belong memory address and next, previos elements address.

## Installation

```sh
go get github.com/rytsh/casset
```

## Usage

```go
l := casset.NewMemory("My First Element").Current
l.Next("second element").Next(3.14).Next(struct{ v string }{v: "4th element"})

for e := l.GetMemory().GetFront(); e != nil; e = e.GetNextElement() {
    fmt.Println(e.GetValue())
}
```

## License

[MIT](./LICENSE)
