package main

import (
	"encoding/json"
	"fmt"
)

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Player struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Pos  Vec2   `json:"pos"`
	secret string // not exported → omitted from JSON
}

func main() {
	p := Player{ID: 1, Name: "Ada", Pos: Vec2{X: 3, Y: 4}, secret: "x"}
	b, _ := json.Marshal(p)
	fmt.Println(string(b))

	var q Player
	_ = json.Unmarshal(b, &q)
	fmt.Printf("%+v\n", q)
}
