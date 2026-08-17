package main

import "fmt"

type Logger struct{}

func (Logger) Log(msg string) { fmt.Println("log:", msg) }

type Server struct {
	Logger
	Addr string
}

func main() {
	s := Server{Addr: ":8080"}
	s.Log("listening on " + s.Addr)
}
