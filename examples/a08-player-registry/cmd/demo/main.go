package main

import (
	"fmt"

	"github.com/willyliao/golearning/examples/a08-player-registry"
)

func main() {
	r := registry.New()
	_ = r.Add(registry.Player{ID: 1, Name: "Ada", Score: 10})
	_ = r.Add(registry.Player{ID: 2, Name: "Lin", Score: 7})
	b, _ := r.MarshalJSON()
	fmt.Println(string(b))
}
