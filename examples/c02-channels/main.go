package main

import "fmt"

func main() {
	ch := make(chan string, 2)
	ch <- "ping"
	ch <- "pong"
	close(ch)

	for v := range ch {
		fmt.Println(v)
	}
}
