package main

import (
	"flag"
	"fmt"
)

func main() {
	n := flag.Int("n", 1, "times to greet")
	name := flag.String("name", "world", "name")
	flag.Parse()
	for i := 0; i < *n; i++ {
		fmt.Println("hello", *name)
	}
}
