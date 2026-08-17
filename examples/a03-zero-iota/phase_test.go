package main

import "testing"

func TestPhaseZeroIsLobby(t *testing.T) {
	var p Phase
	if p != PhaseLobby {
		t.Fatalf("zero phase=%v want lobby", p)
	}
	if p.String() != "lobby" {
		t.Fatalf("String=%q", p.String())
	}
}
