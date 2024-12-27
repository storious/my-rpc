package main

import "fmt"

type Number interface {
	~int | ~float64 | ~float32 | ~int32 | ~int64
}

type Adder[T Number] interface {
	Add(T, T) T
}

type AddInt struct {
}

func (a AddInt) Add(x, y int) int {
	return x + y
}

type Float float32
type AddFloat struct {
}

func (a AddFloat) Add(x, y Float) Float {
	return x + y
}

func main() {
	var a Adder[int] = AddInt{}
	fmt.Println(a.Add(1, 2))
	var b Adder[Float] = AddFloat{}
	fmt.Println(b.Add(1.2, 2))
}
