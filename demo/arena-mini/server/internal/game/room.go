package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Phase int

const (
	PhaseLobby Phase = iota
	PhasePlaying
	PhaseEnded
)

func (p Phase) String() string {
	switch p {
	case PhaseLobby:
		return "lobby"
	case PhasePlaying:
		return "playing"
	case PhaseEnded:
		return "ended"
	default:
		return fmt.Sprintf("phase(%d)", int(p))
	}
}

var (
	ErrInvalidPhase = errors.New("invalid phase")
	ErrRoomFull     = errors.New("room full")
	ErrBadInput     = errors.New("bad input")
	ErrNotInRoom    = errors.New("not in room")
	ErrExists       = errors.New("already joined")
)

type Input struct {
	Player string
	DX, DY int
}

type Player struct {
	Name   string `json:"name"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Ready  bool   `json:"ready"`
	Score  int    `json:"score"`
	lastDX int
	lastDY int
	iFrame int // ticks until can score again
}

// Room is pure game state + rules (no network I/O).
type Room struct {
	mu sync.Mutex

	ID           string
	Capacity     int
	Width        int
	Height       int
	Speed        int
	HitRadius    int // collision distance
	HitCooldown   int // i-frames after a scoring bump
	ScoreToWin   int
	EndHoldTicks int // how long to show Ended before lobby reset
	MaxTicks     int // optional match timeout (0 = off)

	Phase    Phase
	Tick     int
	EndTicks int
	Winner   string
	Players  map[string]*Player
	inbox    []Input
	spawnN   int
}

type Config struct {
	ID           string
	Capacity     int
	Width        int
	Height       int
	Speed        int
	HitRadius    int
	HitCooldown   int
	ScoreToWin   int
	EndHoldTicks int
	MaxTicks     int
}

func NewRoom(cfg Config) *Room {
	if cfg.Capacity <= 0 {
		cfg.Capacity = 4
	}
	if cfg.Width <= 0 {
		cfg.Width = 400
	}
	if cfg.Height <= 0 {
		cfg.Height = 240
	}
	if cfg.Speed <= 0 {
		cfg.Speed = 8
	}
	if cfg.HitRadius <= 0 {
		cfg.HitRadius = 28
	}
	if cfg.HitCooldown <= 0 {
		cfg.HitCooldown = 15 // 0.75s at 20Hz
	}
	if cfg.ScoreToWin <= 0 {
		cfg.ScoreToWin = 3
	}
	if cfg.EndHoldTicks <= 0 {
		cfg.EndHoldTicks = 60 // 3s at 20Hz
	}
	if cfg.MaxTicks <= 0 {
		cfg.MaxTicks = 20 * 60 // 60s at 20Hz
	}
	return &Room{
		ID:           cfg.ID,
		Capacity:     cfg.Capacity,
		Width:        cfg.Width,
		Height:       cfg.Height,
		Speed:        cfg.Speed,
		HitRadius:    cfg.HitRadius,
		HitCooldown:   cfg.HitCooldown,
		ScoreToWin:   cfg.ScoreToWin,
		EndHoldTicks: cfg.EndHoldTicks,
		MaxTicks:     cfg.MaxTicks,
		Phase:        PhaseLobby,
		Players:      map[string]*Player{},
	}
}

func (r *Room) Join(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Phase != PhaseLobby {
		return fmt.Errorf("join: %w", ErrInvalidPhase)
	}
	if _, ok := r.Players[name]; ok {
		return ErrExists
	}
	if len(r.Players) >= r.Capacity {
		return ErrRoomFull
	}
	x, y := r.spawnPosLocked()
	r.Players[name] = &Player{Name: name, X: x, Y: y, Ready: false}
	return nil
}

func (r *Room) spawnPosLocked() (x, y int) {
	x = 40 + (r.spawnN%4)*80
	y = 40 + (r.spawnN/4)*50
	r.spawnN++
	if x > r.Width {
		x = r.Width / 2
	}
	if y > r.Height {
		y = r.Height / 2
	}
	return x, y
}

func (r *Room) Leave(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Players, name)
	if r.Phase == PhasePlaying && len(r.Players) < 2 {
		r.endLocked("forfeit")
	}
}

func (r *Room) Ready(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Phase != PhaseLobby {
		return fmt.Errorf("ready: %w", ErrInvalidPhase)
	}
	p, ok := r.Players[name]
	if !ok {
		return ErrNotInRoom
	}
	p.Ready = true
	r.tryStartLocked()
	return nil
}

func (r *Room) tryStartLocked() {
	if len(r.Players) < 2 {
		return
	}
	for _, p := range r.Players {
		if !p.Ready {
			return
		}
	}
	r.Phase = PhasePlaying
	r.Tick = 0
	r.EndTicks = 0
	r.Winner = ""
	r.inbox = r.inbox[:0]
	for _, p := range r.Players {
		p.Score = 0
		p.iFrame = 0
		p.lastDX, p.lastDY = 0, 0
	}
}

func (r *Room) PhaseSnapshot() Phase {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Phase
}

func (r *Room) PlayerCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.Players)
}

func (r *Room) PushInput(in Input) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Phase != PhasePlaying {
		return fmt.Errorf("input: %w", ErrInvalidPhase)
	}
	if abs(in.DX) > 1 || abs(in.DY) > 1 {
		return fmt.Errorf("%w: dx,dy in -1..1", ErrBadInput)
	}
	if _, ok := r.Players[in.Player]; !ok {
		return ErrNotInRoom
	}
	r.inbox = append(r.inbox, in)
	return nil
}

// TickOnce advances simulation. Safe in Lobby (no-op), Playing, or Ended (countdown→Lobby).
func (r *Room) TickOnce() {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch r.Phase {
	case PhaseLobby:
		return
	case PhaseEnded:
		r.EndTicks++
		if r.EndTicks >= r.EndHoldTicks {
			r.resetLobbyLocked()
		}
		return
	case PhasePlaying:
		// continue below
	default:
		return
	}

	r.Tick++
	// decay i-frames
	for _, p := range r.Players {
		if p.iFrame > 0 {
			p.iFrame--
		}
		p.lastDX, p.lastDY = 0, 0
	}

	last := map[string]Input{}
	for _, in := range r.inbox {
		last[in.Player] = in
	}
	r.inbox = r.inbox[:0]
	for name, in := range last {
		p := r.Players[name]
		p.lastDX, p.lastDY = in.DX, in.DY
		p.X = clamp(p.X+in.DX*r.Speed, 0, r.Width)
		p.Y = clamp(p.Y+in.DY*r.Speed, 0, r.Height)
	}

	r.resolveCollisionsLocked()
	r.checkWinLocked()
}

func (r *Room) resolveCollisionsLocked() {
	names := make([]string, 0, len(r.Players))
	for n := range r.Players {
		names = append(names, n)
	}
	hr2 := r.HitRadius * r.HitRadius
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			a := r.Players[names[i]]
			b := r.Players[names[j]]
			if a.iFrame > 0 || b.iFrame > 0 {
				continue
			}
			dx := a.X - b.X
			dy := a.Y - b.Y
			if dx*dx+dy*dy > hr2 {
				continue
			}
			// Award point to the more "active" bumper this tick; tie-break by name.
			am := abs(a.lastDX) + abs(a.lastDY)
			bm := abs(b.lastDX) + abs(b.lastDY)
			if am > bm || (am == bm && a.Name < b.Name) {
				a.Score++
			} else {
				b.Score++
			}
			a.iFrame = r.HitCooldown
			b.iFrame = r.HitCooldown
			// small separation so they don't stick
			if a.X <= b.X {
				a.X = clamp(a.X-4, 0, r.Width)
				b.X = clamp(b.X+4, 0, r.Width)
			} else {
				a.X = clamp(a.X+4, 0, r.Width)
				b.X = clamp(b.X-4, 0, r.Width)
			}
		}
	}
}

func (r *Room) checkWinLocked() {
	for _, p := range r.Players {
		if p.Score >= r.ScoreToWin {
			r.endLocked(p.Name)
			return
		}
	}
	if r.MaxTicks > 0 && r.Tick >= r.MaxTicks {
		// highest score wins; tie => empty winner
		best, bestScore := "", -1
		tie := false
		for _, p := range r.Players {
			if p.Score > bestScore {
				best, bestScore = p.Name, p.Score
				tie = false
			} else if p.Score == bestScore {
				tie = true
			}
		}
		if tie || bestScore < 0 {
			r.endLocked("")
		} else {
			r.endLocked(best)
		}
	}
}

func (r *Room) endLocked(winner string) {
	r.Phase = PhaseEnded
	r.Winner = winner
	r.EndTicks = 0
	r.inbox = r.inbox[:0]
}

func (r *Room) resetLobbyLocked() {
	r.Phase = PhaseLobby
	r.Tick = 0
	r.EndTicks = 0
	r.Winner = ""
	r.inbox = r.inbox[:0]
	r.spawnN = 0
	for _, p := range r.Players {
		p.Ready = false
		p.Score = 0
		p.iFrame = 0
		p.lastDX, p.lastDY = 0, 0
		p.X, p.Y = r.spawnPosLocked()
	}
}

type Snapshot struct {
	Room       string         `json:"room"`
	Phase      string         `json:"phase"`
	Tick       int            `json:"tick"`
	Width      int            `json:"width"`
	Height     int            `json:"height"`
	HitRadius  int            `json:"hit_radius"`
	ScoreToWin int            `json:"score_to_win"`
	Winner     string         `json:"winner,omitempty"`
	Players    []PlayerPublic `json:"players"`
}

type PlayerPublic struct {
	Name  string `json:"name"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Ready bool   `json:"ready"`
	Score int    `json:"score"`
}

func (r *Room) Snapshot() Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	players := make([]PlayerPublic, 0, len(r.Players))
	for _, p := range r.Players {
		players = append(players, PlayerPublic{
			Name: p.Name, X: p.X, Y: p.Y, Ready: p.Ready, Score: p.Score,
		})
	}
	return Snapshot{
		Room:       r.ID,
		Phase:      r.Phase.String(),
		Tick:       r.Tick,
		Width:      r.Width,
		Height:     r.Height,
		HitRadius:  r.HitRadius,
		ScoreToWin: r.ScoreToWin,
		Winner:     r.Winner,
		Players:    players,
	}
}

func (r *Room) SnapshotJSON() string {
	b, _ := json.Marshal(r.Snapshot())
	return string(b)
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
