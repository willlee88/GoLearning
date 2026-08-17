package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSecretOmitted(t *testing.T) {
	b, err := json.Marshal(Player{Name: "A", HP: 1, secret: "nope"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "nope") {
		t.Fatalf("secret leaked: %s", b)
	}
}
