package main

import (
	"sync"
	"testing"
)

func TestIncConcurrent(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	const n = 500
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if c.Value() != n {
		t.Fatalf("got %d", c.Value())
	}
}
