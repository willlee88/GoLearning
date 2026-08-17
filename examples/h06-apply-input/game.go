package game

import (
	"errors"
	"fmt"
)

var (
	ErrBadInput     = errors.New("bad input")
	ErrInvalidPhase = errors.New("invalid phase")
)

type Phase int

const (
	PhaseLobby Phase = iota
	PhasePlaying
)

type Input struct {
	Player string
	DX, DY int
}

type Player struct {
	Name string
	X, Y int
}

type Room struct {
	Phase   Phase
	Width   int
	Height  int
	Speed   int
	Players map[string]*Player
	inbox   []Input
}

func NewRoom(w, h, speed int) *Room {
	return &Room{
		Width:   w,
		Height:  h,
		Speed:   speed,
		Players: map[string]*Player{},
	}
}

func (r *Room) Add(name string, x, y int) {
	r.Players[name] = &Player{Name: name, X: x, Y: y}
}

func (r *Room) Start() { r.Phase = PhasePlaying }

func (r *Room) PushInput(in Input) error {
	if r.Phase != PhasePlaying {
		return ErrInvalidPhase
	}
	if abs(in.DX) > 1 || abs(in.DY) > 1 {
		return fmt.Errorf("%w: dx dy must be -1..1", ErrBadInput)
	}
	if _, ok := r.Players[in.Player]; !ok {
		return errors.New("unknown player")
	}
	r.inbox = append(r.inbox, in)
	return nil
}

// Tick applies queued inputs once (last input per player wins).
func (r *Room) Tick() {
	if r.Phase != PhasePlaying {
		return
	}
	last := map[string]Input{}
	for _, in := range r.inbox {
		last[in.Player] = in
	}
	r.inbox = r.inbox[:0]
	for name, in := range last {
		p := r.Players[name]
		p.X = clamp(p.X+in.DX*r.Speed, 0, r.Width)
		p.Y = clamp(p.Y+in.DY*r.Speed, 0, r.Height)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
