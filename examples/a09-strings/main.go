package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "你好Go"
	fmt.Println("bytes:", len(s))
	fmt.Println("runes:", utf8.RuneCountInString(s))
	for i, r := range s {
		fmt.Printf("offset=%d rune=%c\n", i, r)
	}
}
