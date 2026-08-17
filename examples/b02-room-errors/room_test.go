package room_test

import (
	"errors"
	"testing"

	"github.com/willyliao/golearning/examples/b02-room-errors"
)

func TestJoinFullWrapped(t *testing.T) {
	r, err := room.New("r1", 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Join("a"); err != nil {
		t.Fatal(err)
	}
	err = r.Join("b")
	if !errors.Is(err, room.ErrRoomFull) {
		t.Fatalf("got %v", err)
	}
}

func TestInvalidPhase(t *testing.T) {
	r, _ := room.New("r1", 2)
	r.Phase = room.PhasePlaying
	err := r.Join("a")
	if !errors.Is(err, room.ErrInvalidPhase) {
		t.Fatalf("got %v", err)
	}
}
