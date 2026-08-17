package game_test

import (
	"errors"
	"testing"

	"github.com/willyliao/golearning/demo/arena-mini/internal/game"
)

func TestLobbyReadyStartAndMove(t *testing.T) {
	r := game.NewRoom(game.Config{ID: "t", Capacity: 4, Width: 100, Height: 100, Speed: 5})
	if err := r.Join("a"); err != nil {
		t.Fatal(err)
	}
	if err := r.Join("b"); err != nil {
		t.Fatal(err)
	}
	_ = r.Ready("a")
	if r.PhaseSnapshot() != game.PhaseLobby {
		t.Fatal("need both ready")
	}
	_ = r.Ready("b")
	if r.PhaseSnapshot() != game.PhasePlaying {
		t.Fatalf("phase=%v", r.PhaseSnapshot())
	}

	snap := r.Snapshot()
	var ax int
	for _, p := range snap.Players {
		if p.Name == "a" {
			ax = p.X
		}
	}
	if err := r.PushInput(game.Input{Player: "a", DX: 1, DY: 0}); err != nil {
		t.Fatal(err)
	}
	r.TickOnce()
	snap = r.Snapshot()
	for _, p := range snap.Players {
		if p.Name == "a" && p.X != ax+5 {
			t.Fatalf("x=%d want %d", p.X, ax+5)
		}
	}
}

func TestRejectClientStyleTeleport(t *testing.T) {
	r := game.NewRoom(game.Config{ID: "t", Speed: 5})
	_ = r.Join("a")
	_ = r.Join("b")
	_ = r.Ready("a")
	_ = r.Ready("b")
	err := r.PushInput(game.Input{Player: "a", DX: 9, DY: 0})
	if !errors.Is(err, game.ErrBadInput) {
		t.Fatalf("%v", err)
	}
}

func TestCollisionScoresAndWin(t *testing.T) {
	r := game.NewRoom(game.Config{
		ID: "t", Width: 200, Height: 200, Speed: 0,
		HitRadius: 50, HitCooldown: 5, ScoreToWin: 2, MaxTicks: 0,
	})
	_ = r.Join("a")
	_ = r.Join("b")
	// force positions close
	// start game
	_ = r.Ready("a")
	_ = r.Ready("b")
	// place on top of each other via inputs with speed 0 won't move — use Tick with manual?
	// Speed 0 means no move; set by pushing and we need positions equal.
	// Re-create with speed and many ticks toward each other is fragile.
	// Use ScoreToWin and collide by joining spawn close: spawn is 40,40 and 120,40 — not close.
	// Directly test by many ticks with high speed toward each other.
	r = game.NewRoom(game.Config{
		ID: "t2", Width: 400, Height: 240, Speed: 20,
		HitRadius: 30, HitCooldown: 2, ScoreToWin: 1, EndHoldTicks: 3, MaxTicks: 0,
	})
	_ = r.Join("ann")
	_ = r.Join("bob")
	_ = r.Ready("ann")
	_ = r.Ready("bob")
	// ann at left spawn, bob to the right — move ann right, bob left
	for i := 0; i < 20 && r.PhaseSnapshot() == game.PhasePlaying; i++ {
		_ = r.PushInput(game.Input{Player: "ann", DX: 1, DY: 0})
		_ = r.PushInput(game.Input{Player: "bob", DX: -1, DY: 0})
		r.TickOnce()
	}
	if r.PhaseSnapshot() != game.PhaseEnded {
		t.Fatalf("expected ended after collision score, phase=%v snap=%+v", r.PhaseSnapshot(), r.Snapshot())
	}
	// hold then lobby
	r.TickOnce()
	r.TickOnce()
	r.TickOnce()
	if r.PhaseSnapshot() != game.PhaseLobby {
		t.Fatalf("expected lobby after hold, phase=%v", r.PhaseSnapshot())
	}
	for _, p := range r.Snapshot().Players {
		if p.Ready || p.Score != 0 {
			t.Fatalf("lobby reset incomplete: %+v", p)
		}
	}
}

func TestTimeoutHighestScoreWins(t *testing.T) {
	r := game.NewRoom(game.Config{
		ID: "t", Width: 100, Height: 100, Speed: 0,
		HitRadius: 1, ScoreToWin: 99, MaxTicks: 3, EndHoldTicks: 100,
	})
	_ = r.Join("a")
	_ = r.Join("b")
	_ = r.Ready("a")
	_ = r.Ready("b")
	// no collision possible with tiny radius and speed 0; timeout
	r.TickOnce()
	r.TickOnce()
	r.TickOnce()
	if r.PhaseSnapshot() != game.PhaseEnded {
		t.Fatalf("phase=%v", r.PhaseSnapshot())
	}
}
