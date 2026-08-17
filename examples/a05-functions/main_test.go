package main

import "testing"

func TestFindPlayer(t *testing.T) {
	ps := []Player{{1, "A"}, {2, "B"}}
	p, ok := FindPlayer(ps, 2)
	if !ok || p.Name != "B" {
		t.Fatalf("got %+v ok=%v", p, ok)
	}
	_, ok = FindPlayer(ps, 9)
	if ok {
		t.Fatal("expected not found")
	}
}

func TestJoinFull(t *testing.T) {
	if err := Join(2, 2); err == nil {
		t.Fatal("expected error")
	}
}
