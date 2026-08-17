package main

import "testing"

func TestHandle(t *testing.T) {
	if Handle(JoinCommand{"A"}) != "join:A" {
		t.Fatal()
	}
	if Handle(LeaveCommand{}) != "leave" {
		t.Fatal()
	}
	if Handle(42) != "unknown" {
		t.Fatal()
	}
}
