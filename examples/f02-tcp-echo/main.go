package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9000", "listen or dial address")
	client := flag.Bool("client", false, "run as client")
	flag.Parse()

	if *client {
		runClient(*addr)
		return
	}
	runServer(*addr)
}

func runServer(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	log.Printf("echo listening on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(60 * time.Second))
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("read: %v", err)
			}
			return
		}
		if _, err := io.WriteString(conn, line); err != nil {
			return
		}
	}
}

func runClient(addr string) {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Fprintln(conn, "hello-tcp")
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(buf[:n]))
}
