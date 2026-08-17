package safe

import "sync"

// Lobby is safe for concurrent Join.
type Lobby struct {
	mu    sync.Mutex
	Count int
}

func (l *Lobby) Join() {
	l.mu.Lock()
	l.Count++
	l.mu.Unlock()
}
