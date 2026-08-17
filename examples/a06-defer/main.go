package main

import (
	"fmt"
	"time"
)

func order() {
	defer fmt.Println("first deferred (runs last)")
	defer fmt.Println("second deferred (runs first among defers)")
	fmt.Println("body")
}

func timed(name string, fn func()) {
	start := time.Now()
	defer func() {
		fmt.Printf("%s took %s\n", name, time.Since(start))
	}()
	fn()
}

func main() {
	order()
	timed("sleep", func() { time.Sleep(20 * time.Millisecond) })
}
