package broken_test

import (
	"sync"
	"testing"

	"github.com/willyliao/golearning/examples/c06-race-lab/broken"
)

// This test is expected to FAIL under -race.
func TestBrokenJoinRace(t *testing.T) {
	var l broken.Lobby
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Join()
		}()
	}
	wg.Wait()
	if l.Count != 100 {
		// may or may not fail without race detector
		t.Logf("count=%d (nondeterministic without proper sync)", l.Count)
	}
}
