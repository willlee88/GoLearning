package main

import (
	"fmt"

	"github.com/willyliao/golearning/examples/a02-packages/player"
)

func main() {
	p := player.New("Ada", 100)
	p.Damage(15)
	fmt.Printf("%s hp=%d\n", p.Name, p.HP())
}
