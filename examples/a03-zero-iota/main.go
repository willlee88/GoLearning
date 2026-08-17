package main

import "fmt"

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
		return fmt.Sprintf("Phase(%d)", int(p))
	}
}

func main() {
	var phase Phase // zero value == PhaseLobby
	fmt.Println("zero phase:", phase)

	var n int
	var s string
	var ok bool
	fmt.Printf("zero int=%d string=%q bool=%v\n", n, s, ok)

	phase = PhasePlaying
	fmt.Println("now:", phase)
}
