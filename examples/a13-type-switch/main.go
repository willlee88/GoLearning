package main

import "fmt"

type JoinCommand struct{ Name string }
type MoveCommand struct{ X, Y int }
type LeaveCommand struct{}

func Handle(cmd any) string {
	switch c := cmd.(type) {
	case JoinCommand:
		return "join:" + c.Name
	case MoveCommand:
		return fmt.Sprintf("move:%d,%d", c.X, c.Y)
	case LeaveCommand:
		return "leave"
	default:
		return "unknown"
	}
}

func main() {
	fmt.Println(Handle(JoinCommand{"Ada"}))
	fmt.Println(Handle(MoveCommand{1, 2}))
	fmt.Println(Handle(LeaveCommand{}))
}
