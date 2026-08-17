package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	h := newHub()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "rooms": h.stats()})
	})
	mux.Handle("GET /ws", websocket.Handler(h.handleWS))
	web := env("WEB_DIR", "web")
	mux.Handle("GET /", http.FileServer(http.Dir(web)))

	addr := env("ADDR", ":8090")
	log.Printf("f08-ws-room on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

type hub struct {
	mu    sync.RWMutex
	rooms map[string]*room
}

func newHub() *hub { return &hub{rooms: map[string]*room{}} }

func (h *hub) stats() map[string]int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make(map[string]int, len(h.rooms))
	for id, r := range h.rooms {
		out[id] = r.size()
	}
	return out
}

func (h *hub) get(id string) *room {
	h.mu.Lock()
	defer h.mu.Unlock()
	if r, ok := h.rooms[id]; ok {
		return r
	}
	r := &room{id: id, members: map[*websocket.Conn]string{}}
	h.rooms[id] = r
	return r
}

type envelope struct {
	V       int    `json:"v"`
	Type    string `json:"type"`
	From    string `json:"from,omitempty"`
	Room    string `json:"room,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type room struct {
	id      string
	mu      sync.Mutex
	members map[*websocket.Conn]string
}

func (r *room) size() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.members)
}

func (h *hub) handleWS(conn *websocket.Conn) {
	req := conn.Request()
	roomID := req.URL.Query().Get("room")
	if roomID == "" {
		roomID = "lobby"
	}
	name := req.URL.Query().Get("name")
	if name == "" {
		name = "anon"
	}
	r := h.get(roomID)
	r.mu.Lock()
	r.members[conn] = name
	r.mu.Unlock()

	r.broadcast(envelope{V: 1, Type: "system", Room: roomID, Payload: name + " joined"})
	defer func() {
		r.mu.Lock()
		delete(r.members, conn)
		r.mu.Unlock()
		_ = conn.Close()
		r.broadcast(envelope{V: 1, Type: "system", Room: roomID, Payload: name + " left"})
	}()

	for {
		var raw string
		if err := websocket.Message.Receive(conn, &raw); err != nil {
			return
		}
		var env envelope
		if err := json.Unmarshal([]byte(raw), &env); err != nil {
			env = envelope{V: 1, Type: "chat", Payload: raw}
		}
		if env.V == 0 {
			env.V = 1
		}
		if env.Type == "" {
			env.Type = "chat"
		}
		if env.Type == "ping" {
			_ = websocket.Message.Send(conn, mustJSON(envelope{V: 1, Type: "pong", Payload: time.Now().UTC().Format(time.RFC3339)}))
			continue
		}
		env.From = name
		env.Room = roomID
		r.broadcast(env)
	}
}

func (r *room) broadcast(env envelope) {
	data := mustJSON(env)
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.members {
		_ = websocket.Message.Send(c, data)
	}
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
