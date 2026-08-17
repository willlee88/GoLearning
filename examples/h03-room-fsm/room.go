package room

import (
	"errors"
	"fmt"
)

type Phase int

const (
	PhaseLobby Phase = iota
	PhasePlaying
	PhaseEnded
)

var (
	ErrInvalidPhase = errors.New("invalid phase")
	ErrRoomFull     = errors.New("room full")
	ErrNotReady     = errors.New("not all ready")
	ErrTooFew       = errors.New("need at least 2 players")
	ErrExists       = errors.New("already joined")
)

type Room struct {
	Capacity int
	Phase    Phase
	Members  map[string]bool // name -> ready
}

func New(cap int) *Room {
	return &Room{Capacity: cap, Members: map[string]bool{}}
}

func (r *Room) Join(name string) error {
	if r.Phase != PhaseLobby {
		return fmt.Errorf("join: %w", ErrInvalidPhase)
	}
	if len(r.Members) >= r.Capacity {
		return ErrRoomFull
	}
	if _, ok := r.Members[name]; ok {
		return ErrExists
	}
	r.Members[name] = false
	return nil
}

func (r *Room) Ready(name string) error {
	if r.Phase != PhaseLobby {
		return fmt.Errorf("ready: %w", ErrInvalidPhase)
	}
	if _, ok := r.Members[name]; !ok {
		return errors.New("not in room")
	}
	r.Members[name] = true
	return r.tryStart()
}

func (r *Room) tryStart() error {
	if len(r.Members) < 2 {
		return nil // not an error; just wait
	}
	for _, ready := range r.Members {
		if !ready {
			return nil
		}
	}
	r.Phase = PhasePlaying
	return nil
}

func (r *Room) End() error {
	if r.Phase != PhasePlaying {
		return ErrInvalidPhase
	}
	r.Phase = PhaseEnded
	return nil
}
