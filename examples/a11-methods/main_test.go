package main

import "testing"

func TestInc(t *testing.T) {
	var c Counter
	c.Inc()
	if c.Value() != 1 {
		t.Fatal(c.Value())
	}
}
