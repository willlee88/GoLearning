package main

import "testing"

func TestBroadcast(t *testing.T) {
	if err := Broadcast(LogNotifier{}, "x"); err != nil {
		t.Fatal(err)
	}
	if err := Broadcast(nil, "x"); err == nil {
		t.Fatal("expected error")
	}
}
