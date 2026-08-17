package store_test

import (
	"context"
	"errors"
	"testing"

	store "github.com/willyliao/golearning/examples/i02-sql-memory"
)

func TestMemoryRepo(t *testing.T) {
	r := store.NewMemoryRepo()
	ctx := context.Background()
	if err := r.Upsert(ctx, store.Player{ID: "1", Name: "Ada", Score: 10}); err != nil {
		t.Fatal(err)
	}
	p, err := r.Get(ctx, "1")
	if err != nil || p.Name != "Ada" {
		t.Fatalf("%+v %v", p, err)
	}
	_, err = r.Get(ctx, "nope")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatal(err)
	}
}
