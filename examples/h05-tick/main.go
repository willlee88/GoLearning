package main

import (
	"fmt"
	"time"
)

// Minimal tick loop demo: accumulate "inputs", apply each tick.
func main() {
	inbox := make(chan int, 8)
	go func() {
		for i := 1; i <= 5; i++ {
			inbox <- i
			time.Sleep(30 * time.Millisecond)
		}
		close(inbox)
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	var buf []int
	var sum int
	ticks := 0
	open := true

	for open || len(buf) > 0 {
		select {
		case v, ok := <-inbox:
			if !ok {
				open = false
				inbox = nil
				continue
			}
			buf = append(buf, v)
		case <-ticker.C:
			ticks++
			for _, v := range buf {
				sum += v
			}
			fmt.Printf("tick=%d applied=%v sum=%d\n", ticks, buf, sum)
			buf = buf[:0]
			if !open && len(buf) == 0 && ticks >= 3 {
				return
			}
		}
	}
}
