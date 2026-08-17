package main

import (
	"errors"
	"fmt"
)

var ErrBusy = errors.New("queue full")

func trySend(ch chan<- int, v int) error {
	select {
	case ch <- v:
		return nil
	default:
		return ErrBusy
	}
}

func main() {
	ch := make(chan int, 2)
	for i := 0; i < 5; i++ {
		if err := trySend(ch, i); err != nil {
			fmt.Println("drop", i, err)
		} else {
			fmt.Println("enqueued", i)
		}
	}
}
