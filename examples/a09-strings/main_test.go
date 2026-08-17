package main

import (
	"testing"
	"unicode/utf8"
)

func TestRuneCount(t *testing.T) {
	if utf8.RuneCountInString("你好") != 2 {
		t.Fatal()
	}
}
