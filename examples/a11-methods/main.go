package main

import "fmt"

type Counter struct{ n int }

func (c Counter) Value() int { return c.n }

func (c *Counter) Inc() {
	if c == nil {
		return
	}
	c.n++
}

func main() {
	var c Counter
	c.Inc()
	c.Inc()
	fmt.Println(c.Value())
}
