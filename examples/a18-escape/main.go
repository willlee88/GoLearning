package main

import "fmt"

//go:noinline
func answer() *int {
	x := 42
	return &x // typically escapes to heap
}

func main() {
	p := answer()
	fmt.Println(*p)
}
