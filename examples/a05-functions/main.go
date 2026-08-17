package main

import (
	"errors"
	"fmt"
)

type Player struct {
	ID   int64
	Name string
}

func FindPlayer(ps []Player, id int64) (Player, bool) {
	for _, p := range ps {
		if p.ID == id {
			return p, true
		}
	}
	return Player{}, false
}

func Join(capacity int, count int) error {
	if count >= capacity {
		return errors.New("room full")
	}
	return nil
}

func main() {
	ps := []Player{{1, "Ada"}, {2, "Lin"}}
	if p, ok := FindPlayer(ps, 2); ok {
		fmt.Println("found", p.Name)
	}
	if err := Join(2, 2); err != nil {
		fmt.Println("join:", err)
	}
}
