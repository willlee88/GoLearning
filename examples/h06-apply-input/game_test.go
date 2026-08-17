package game_test

import (
	"errors"
	"testing"

	"github.com/willyliao/golearning/examples/h06-apply-input"
)

func TestApplyInputAuthoritative(t *testing.T) {
	r := game.NewRoom(100, 100, 5)
	r.Add("Ada", 10, 10)
	r.Start()
	if err := r.PushInput(game.Input{Player: "Ada", DX: 1, DY: 0}); err != nil {
		t.Fatal(err)
	}
	r.Tick()
	if r.Players["Ada"].X != 15 {
		t.Fatalf("x=%d", r.Players["Ada"].X)
	}
}

func TestRejectBadInput(t *testing.T) {
	r := game.NewRoom(100, 100, 5)
	r.Add("Ada", 0, 0)
	r.Start()
	err := r.PushInput(game.Input{Player: "Ada", DX: 5, DY: 0})
	if !errors.Is(err, game.ErrBadInput) {
		t.Fatalf("%v", err)
	}
}

func TestInputInLobbyFails(t *testing.T) {
	r := game.NewRoom(100, 100, 5)
	r.Add("Ada", 0, 0)
	err := r.PushInput(game.Input{Player: "Ada", DX: 1, DY: 0})
	if !errors.Is(err, game.ErrInvalidPhase) {
		t.Fatalf("%v", err)
	}
}
