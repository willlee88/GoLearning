package registry_test

import (
	"errors"
	"testing"

	registry "github.com/willyliao/golearning/examples/a08-player-registry"
)

func TestAddGetRemove(t *testing.T) {
	r := registry.New()
	if err := r.Add(registry.Player{ID: 1, Name: "Ada", Score: 3}); err != nil {
		t.Fatal(err)
	}
	if err := r.Add(registry.Player{ID: 1, Name: "Dup"}); !errors.Is(err, registry.ErrExists) {
		t.Fatalf("err=%v", err)
	}
	p, err := r.Get(1)
	if err != nil || p.Name != "Ada" {
		t.Fatalf("get %+v err=%v", p, err)
	}
	snap := r.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("snap len=%d", len(snap))
	}
	if err := r.Remove(1); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Get(1); !errors.Is(err, registry.ErrNotFound) {
		t.Fatalf("err=%v", err)
	}
}
