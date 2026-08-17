package snap

import (
	"encoding/json"
	"testing"
)

type State struct {
	Tick    int      `json:"tick"`
	Players []string `json:"players"`
}

func snapshot() []byte {
	b, _ := json.Marshal(State{Tick: 1, Players: []string{"a", "b", "c", "d"}})
	return b
}

func BenchmarkSnapshot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = snapshot()
	}
}
