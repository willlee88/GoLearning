package room

import (
	"errors"
	"fmt"
)

var (
	ErrRoomFull    = errors.New("room full")
	ErrNotFound    = errors.New("player not found")
	ErrInvalidPhase = errors.New("invalid phase")
)

type Phase int

const (
	PhaseLobby Phase = iota
	PhasePlaying
)

type Room struct {
	ID       string
	Capacity int
	Phase    Phase
	Players  map[string]struct{}
}

func New(id string, capacity int) (*Room, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}
	return &Room{
		ID:       id,
		Capacity: capacity,
		Phase:    PhaseLobby,
		Players:  make(map[string]struct{}),
	}, nil
}

func (r *Room) Join(name string) error {
	if r.Phase != PhaseLobby {
		return fmt.Errorf("join %s: %w", r.ID, ErrInvalidPhase)
	}
	if len(r.Players) >= r.Capacity {
		return fmt.Errorf("join %s: %w", r.ID, ErrRoomFull)
	}
	r.Players[name] = struct{}{}
	return nil
}

func (r *Room) Leave(name string) error {
	if _, ok := r.Players[name]; !ok {
		return fmt.Errorf("leave %s: %w", name, ErrNotFound)
	}
	delete(r.Players, name)
	return nil
}
