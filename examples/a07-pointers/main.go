package main

import "fmt"

type Player struct {
	HP int
}

func healValue(p Player, n int) {
	p.HP += n
}

func healPointer(p *Player, n int) {
	if p == nil {
		return
	}
	p.HP += n
}

func main() {
	p := Player{HP: 10}
	healValue(p, 5)
	fmt.Println("after value heal:", p.HP) // 10

	healPointer(&p, 5)
	fmt.Println("after pointer heal:", p.HP) // 15
}
