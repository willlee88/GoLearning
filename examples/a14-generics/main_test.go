package main

import "testing"

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3}, func(n int) bool { return n > 1 })
	if len(got) != 2 || got[0] != 2 {
		t.Fatalf("%v", got)
	}
}
