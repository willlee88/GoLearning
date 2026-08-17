package main

import "fmt"

type Notifier interface {
	Notify(msg string) error
}

type LogNotifier struct{}

func (LogNotifier) Notify(msg string) error {
	fmt.Println("notify:", msg)
	return nil
}

func Broadcast(n Notifier, msg string) error {
	if n == nil {
		return fmt.Errorf("notifier nil")
	}
	return n.Notify(msg)
}

func main() {
	_ = Broadcast(LogNotifier{}, "room started")

	var p *LogNotifier = nil
	var n Notifier = p
	fmt.Println("n == nil?", n == nil) // false: typed nil inside interface
	fmt.Printf("dynamic: %T\n", n)
}
