package main

import (
	"fmt"
	"time"
)

func main() {
	inbox := make(chan string, 1)
	inbox <- "cmd:ready"

	timeout := time.After(50 * time.Millisecond)
	select {
	case msg := <-inbox:
		fmt.Println("got", msg)
	case <-timeout:
		fmt.Println("timeout")
	}
}
