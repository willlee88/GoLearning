package main

import (
	"encoding/json"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	p := Player{ID: 7, Name: "Lin", Pos: Vec2{1, 2}}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var q Player
	if err := json.Unmarshal(b, &q); err != nil {
		t.Fatal(err)
	}
	if q.ID != 7 || q.Name != "Lin" || q.Pos.X != 1 {
		t.Fatalf("%+v", q)
	}
}
