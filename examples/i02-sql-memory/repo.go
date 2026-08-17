package store

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type Player struct {
	ID    string
	Name  string
	Score int
}

// PlayerRepo is the persistence port (swap memory for sql later).
type PlayerRepo interface {
	Upsert(ctx context.Context, p Player) error
	Get(ctx context.Context, id string) (Player, error)
}

type MemoryRepo struct {
	mu   sync.RWMutex
	data map[string]Player
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{data: map[string]Player{}}
}

func (m *MemoryRepo) Upsert(ctx context.Context, p Player) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	m.data[p.ID] = p
	m.mu.Unlock()
	return nil
}

func (m *MemoryRepo) Get(ctx context.Context, id string) (Player, error) {
	if err := ctx.Err(); err != nil {
		return Player{}, err
	}
	m.mu.RLock()
	p, ok := m.data[id]
	m.mu.RUnlock()
	if !ok {
		return Player{}, ErrNotFound
	}
	return p, nil
}
