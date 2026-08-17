package player

// Player is an exported type. hp is intentionally unexported.
type Player struct {
	Name string
	hp   int
}

func New(name string, hp int) Player {
	return Player{Name: name, hp: hp}
}

func (p Player) HP() int {
	return p.hp
}

func (p *Player) Damage(n int) {
	if p == nil {
		return
	}
	p.hp -= n
	if p.hp < 0 {
		p.hp = 0
	}
}
