// Package broken intentionally races for teaching.
// Run: go test -race ./broken
package broken

// Lobby counts players without synchronization — DO NOT copy this pattern.
type Lobby struct {
	Count int
}

func (l *Lobby) Join() {
	// data race if called concurrently
	l.Count++
}
