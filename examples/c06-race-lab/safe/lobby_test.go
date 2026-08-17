package safe_test

import (
	"sync"
	"testing"

	"github.com/willyliao/golearning/examples/c06-race-lab/safe"
)

func TestSafeJoin(t *testing.T) {
	var l safe.Lobby
	var wg sync.WaitGroup
	const n = 200
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			l.Join()
		}()
	}
	wg.Wait()
	if l.Count != n {
		t.Fatalf("count=%d", l.Count)
	}
}
