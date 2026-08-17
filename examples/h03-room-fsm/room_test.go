package room_test

import (
	"errors"
	"testing"

	"github.com/willyliao/golearning/examples/h03-room-fsm"
)

func TestReadyStartsWhenTwoReady(t *testing.T) {
	r := room.New(4)
	_ = r.Join("a")
	_ = r.Join("b")
	if err := r.Ready("a"); err != nil {
		t.Fatal(err)
	}
	if r.Phase != room.PhaseLobby {
		t.Fatal("should still lobby")
	}
	if err := r.Ready("b"); err != nil {
		t.Fatal(err)
	}
	if r.Phase != room.PhasePlaying {
		t.Fatalf("phase=%v", r.Phase)
	}
}

func TestJoinWhilePlaying(t *testing.T) {
	r := room.New(4)
	_ = r.Join("a")
	_ = r.Join("b")
	_ = r.Ready("a")
	_ = r.Ready("b")
	err := r.Join("c")
	if !errors.Is(err, room.ErrInvalidPhase) {
		t.Fatalf("%v", err)
	}
}
