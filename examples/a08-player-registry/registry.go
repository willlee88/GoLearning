package registry

import (
	"encoding/json"
	"errors"
	"sync"
)

type PlayerID int64

type Player struct {
	ID    PlayerID `json:"id"`
	Name  string   `json:"name"`
	Score int      `json:"score"`
}

// Registry keeps players by ID. Safe for concurrent use.
type Registry struct {
	mu   sync.RWMutex
	byID map[PlayerID]*Player
}

func New() *Registry {
	return &Registry{byID: make(map[PlayerID]*Player)}
}

var (
	ErrExists   = errors.New("player exists")
	ErrNotFound = errors.New("player not found")
)

func (r *Registry) Add(p Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[p.ID]; ok {
		return ErrExists
	}
	cp := p
	r.byID[p.ID] = &cp
	return nil
}

func (r *Registry) Get(id PlayerID) (Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.byID[id]
	if !ok {
		return Player{}, ErrNotFound
	}
	return *p, nil
}

func (r *Registry) Remove(id PlayerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[id]; !ok {
		return ErrNotFound
	}
	delete(r.byID, id)
	return nil
}

// Snapshot returns a deep-ish copy for JSON / broadcast without holding locks outside.
func (r *Registry) Snapshot() []Player {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Player, 0, len(r.byID))
	for _, p := range r.byID {
		out = append(out, *p)
	}
	return out
}

func (r *Registry) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Snapshot())
}
