// Fake WS clients for Arena Mini pressure testing.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	addr := flag.String("addr", "http://127.0.0.1:8080", "base http URL of arena-mini")
	n := flag.Int("n", 20, "number of clients")
	room := flag.String("room", "load", "room id")
	seconds := flag.Int("seconds", 10, "how long to run")
	flag.Parse()

	u, err := url.Parse(*addr)
	if err != nil {
		log.Fatal(err)
	}
	wsScheme := "ws"
	if u.Scheme == "https" {
		wsScheme = "wss"
	}
	base := fmt.Sprintf("%s://%s/ws", wsScheme, u.Host)

	var (
		ok    atomic.Int64
		fail  atomic.Int64
		msgs  atomic.Int64
		wg    sync.WaitGroup
		stop  = time.After(time.Duration(*seconds) * time.Second)
		start = time.Now()
	)

	for i := 0; i < *n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("bot-%d", i)
			wsURL := fmt.Sprintf("%s?room=%s&name=%s", base, url.QueryEscape(*room), url.QueryEscape(name))
			cfg, err := websocket.NewConfig(wsURL, *addr)
			if err != nil {
				fail.Add(1)
				return
			}
			conn, err := websocket.DialConfig(cfg)
			if err != nil {
				fail.Add(1)
				return
			}
			defer conn.Close()
			ok.Add(1)

			// ready then spam input
			_ = websocket.Message.Send(conn, `{"v":1,"type":"ready","payload":"1"}`)
			deadline := time.Now().Add(time.Duration(*seconds) * time.Second)
			for time.Now().Before(deadline) {
				_ = websocket.Message.Send(conn, `{"v":1,"type":"input","payload":"1,0"}`)
				msgs.Add(1)
				var raw string
				_ = conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
				if err := websocket.Message.Receive(conn, &raw); err == nil {
					msgs.Add(1)
				}
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	<-stop
	// give clients a moment; they exit on their deadline
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("clients_ok=%d fail=%d msg_ops≈%d elapsed=%s\n", ok.Load(), fail.Load(), msgs.Load(), elapsed)
}
