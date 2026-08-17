package main

import "testing"

func TestHealPointer(t *testing.T) {
	p := Player{HP: 1}
	healPointer(&p, 2)
	if p.HP != 3 {
		t.Fatalf("HP=%d", p.HP)
	}
	healValue(p, 10)
	if p.HP != 3 {
		t.Fatalf("value heal should not mutate, HP=%d", p.HP)
	}
}
