package player_test

import (
	"testing"

	"github.com/willyliao/golearning/examples/a02-packages/player"
)

func TestDamage(t *testing.T) {
	p := player.New("Ada", 10)
	p.Damage(3)
	if p.HP() != 7 {
		t.Fatalf("HP=%d want 7", p.HP())
	}
}
