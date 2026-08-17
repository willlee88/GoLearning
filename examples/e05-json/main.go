package main

import (
	"encoding/json"
	"fmt"
)

type Player struct {
	Name   string `json:"name"`
	HP     int    `json:"hp"`
	secret string // unexported: omitted
}

func main() {
	p := Player{Name: "Ada", HP: 10, secret: "x"}
	b, _ := json.Marshal(p)
	fmt.Println(string(b))

	var q Player
	_ = json.Unmarshal([]byte(`{"name":"Lin","hp":7,"extra":true}`), &q)
	fmt.Printf("%+v\n", q)
}
