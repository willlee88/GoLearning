package hub

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/willyliao/golearning/demo/arena-mini/internal/game"
	"github.com/willyliao/golearning/demo/arena-mini/internal/metrics"
	"golang.org/x/net/websocket"
)

const tickInterval = 50 * time.Millisecond // 20 Hz

// Hub manages rooms and websocket sessions.
type Hub struct {
	mu      sync.RWMutex
	rooms   map[string]*roomRuntime
	metrics *metrics.Metrics
	log     *slog.Logger
	closed  atomic.Bool
}

func New(m *metrics.Metrics, log *slog.Logger) *Hub {
	if log == nil {
		log = slog.Default()
	}
	if m == nil {
		m = metrics.New()
	}
	return &Hub{
		rooms:   make(map[string]*roomRuntime),
		metrics: m,
		log:     log,
	}
}

func (h *Hub) Metrics() *metrics.Metrics { return h.metrics }

func (h *Hub) RoomCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms)
}

func (h *Hub) RoomStats() map[string]any {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make(map[string]any, len(h.rooms))
	for id, rt := range h.rooms {
		out[id] = map[string]any{
			"players": rt.game.PlayerCount(),
			"phase":   rt.game.PhaseSnapshot().String(),
		}
	}
	return out
}

// Close stops all room tickers (call during graceful shutdown).
func (h *Hub) Close() {
	h.closed.Store(true)
	h.mu.Lock()
	defer h.mu.Unlock()
	for id, rt := range h.rooms {
		rt.stop()
		delete(h.rooms, id)
	}
	h.metrics.Rooms.Store(0)
	h.log.Info("hub closed")
}

type roomRuntime struct {
	game *game.Room
	hub  *Hub

	mu      sync.Mutex
	members map[*websocket.Conn]string

	tickOnce sync.Once
	stopTick chan struct{}
	stopped  atomic.Bool
}

func (h *Hub) getOrCreate(id string) *roomRuntime {
	h.mu.Lock()
	defer h.mu.Unlock()
	if r, ok := h.rooms[id]; ok {
		return r
	}
	r := &roomRuntime{
		game:     game.NewRoom(game.Config{ID: id}),
		hub:      h,
		members:  make(map[*websocket.Conn]string),
		stopTick: make(chan struct{}),
	}
	h.rooms[id] = r
	h.metrics.Rooms.Store(int64(len(h.rooms)))
	h.log.Info("room created", "room", id)
	return r
}

func (rt *roomRuntime) stop() {
	if rt.stopped.Swap(true) {
		return
	}
	select {
	case <-rt.stopTick:
	default:
		close(rt.stopTick)
	}
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	if h.closed.Load() {
		http.Error(w, "shutting down", http.StatusServiceUnavailable)
		return
	}
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		roomID = "lobby"
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "anon"
	}
	server := websocket.Server{
		Handler: func(conn *websocket.Conn) {
			h.serveConn(roomID, name, conn)
		},
	}
	server.ServeHTTP(w, r)
}

type Envelope struct {
	V       int    `json:"v"`
	Type    string `json:"type"`
	From    string `json:"from,omitempty"`
	Room    string `json:"room,omitempty"`
	Payload string `json:"payload,omitempty"`
}

func (h *Hub) serveConn(roomID, name string, conn *websocket.Conn) {
	if h.closed.Load() {
		_ = conn.Close()
		return
	}
	rt := h.getOrCreate(roomID)
	h.metrics.Connections.Add(1)
	defer h.metrics.Connections.Add(-1)

	if err := rt.game.Join(name); err != nil {
		h.metrics.Errors.Add(1)
		_ = websocket.Message.Send(conn, mustJSON(Envelope{
			V: 1, Type: "error", Room: roomID, Payload: err.Error(),
		}))
		_ = conn.Close()
		return
	}

	rt.mu.Lock()
	for c, n := range rt.members {
		if n == name {
			delete(rt.members, c)
			_ = c.Close()
		}
	}
	rt.members[conn] = name
	rt.mu.Unlock()

	h.log.Info("join", "room", roomID, "player", name)
	rt.broadcast(Envelope{V: 1, Type: "system", Room: roomID, Payload: name + " joined"})
	rt.broadcastState()

	defer func() {
		rt.game.Leave(name)
		rt.mu.Lock()
		delete(rt.members, conn)
		empty := len(rt.members) == 0
		rt.mu.Unlock()
		_ = conn.Close()
		rt.broadcast(Envelope{V: 1, Type: "system", Room: roomID, Payload: name + " left"})
		rt.broadcastState()
		if empty {
			h.removeRoom(roomID, rt)
		}
		h.log.Info("leave", "room", roomID, "player", name)
	}()

	for {
		if h.closed.Load() {
			rt.send(conn, Envelope{V: 1, Type: "system", Payload: "server_closing", Room: roomID})
			return
		}
		var raw string
		if err := websocket.Message.Receive(conn, &raw); err != nil {
			return
		}
		h.metrics.MessagesIn.Add(1)

		var msg Envelope
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			h.metrics.Errors.Add(1)
			rt.send(conn, Envelope{V: 1, Type: "error", Payload: "bad json", Room: roomID})
			continue
		}
		if msg.V == 0 {
			msg.V = 1
		}
		if msg.V != 1 {
			h.metrics.Errors.Add(1)
			rt.send(conn, Envelope{V: 1, Type: "error", Payload: "unsupported version", Room: roomID})
			continue
		}

		switch msg.Type {
		case "ping":
			rt.send(conn, Envelope{V: 1, Type: "pong", Payload: time.Now().UTC().Format(time.RFC3339), Room: roomID})
		case "chat":
			msg.From = name
			msg.Room = roomID
			rt.broadcast(msg)
		case "ready":
			if err := rt.game.Ready(name); err != nil {
				h.metrics.Errors.Add(1)
				rt.send(conn, Envelope{V: 1, Type: "error", Payload: err.Error(), Room: roomID})
				continue
			}
			rt.broadcast(Envelope{V: 1, Type: "system", From: name, Room: roomID, Payload: "ready"})
			rt.broadcastState()
			if rt.game.PhaseSnapshot() == game.PhasePlaying {
				rt.ensureTicker()
			}
		case "input":
			dx, dy, err := parseInput(msg.Payload)
			if err != nil {
				h.metrics.Errors.Add(1)
				rt.send(conn, Envelope{V: 1, Type: "error", Payload: err.Error(), Room: roomID})
				continue
			}
			if err := rt.game.PushInput(game.Input{Player: name, DX: dx, DY: dy}); err != nil {
				h.metrics.Errors.Add(1)
				rt.send(conn, Envelope{V: 1, Type: "error", Payload: err.Error(), Room: roomID})
			}
		case "pos":
			h.metrics.Errors.Add(1)
			rt.send(conn, Envelope{V: 1, Type: "error", Payload: "use type=input (authoritative); pos is disabled", Room: roomID})
		default:
			h.metrics.Errors.Add(1)
			rt.send(conn, Envelope{V: 1, Type: "error", Payload: "unknown type " + msg.Type, Room: roomID})
		}
	}
}

func (h *Hub) removeRoom(id string, rt *roomRuntime) {
	h.mu.Lock()
	defer h.mu.Unlock()
	cur, ok := h.rooms[id]
	if !ok || cur != rt {
		return
	}
	rt.stop()
	delete(h.rooms, id)
	h.metrics.Rooms.Store(int64(len(h.rooms)))
	h.log.Info("room removed", "room", id)
}

func (rt *roomRuntime) ensureTicker() {
	rt.tickOnce.Do(func() {
		go rt.tickLoop()
	})
}

func (rt *roomRuntime) tickLoop() {
	t := time.NewTicker(tickInterval)
	defer t.Stop()
	for {
		select {
		case <-rt.stopTick:
			return
		case <-t.C:
			if rt.hub.closed.Load() {
				return
			}
			phase := rt.game.PhaseSnapshot()
			// Lobby: idle. Playing: simulate. Ended: countdown → Lobby.
			if phase == game.PhaseLobby {
				continue
			}
			rt.game.TickOnce()
			if phase == game.PhasePlaying {
				rt.hub.metrics.InputsApplied.Add(1)
			}
			rt.broadcastState()
		}
	}
}

func (rt *roomRuntime) broadcastState() {
	rt.broadcast(Envelope{
		V:       1,
		Type:    "state",
		Room:    rt.game.ID,
		Payload: rt.game.SnapshotJSON(),
	})
}

func (rt *roomRuntime) broadcast(msg Envelope) {
	data := mustJSON(msg)
	rt.mu.Lock()
	defer rt.mu.Unlock()
	for c := range rt.members {
		if err := websocket.Message.Send(c, data); err != nil {
			rt.hub.log.Debug("send error", "err", err)
		} else {
			rt.hub.metrics.MessagesOut.Add(1)
		}
	}
}

func (rt *roomRuntime) send(conn *websocket.Conn, msg Envelope) {
	if err := websocket.Message.Send(conn, mustJSON(msg)); err == nil {
		rt.hub.metrics.MessagesOut.Add(1)
	}
}

func parseInput(payload string) (dx, dy int, err error) {
	payload = strings.TrimSpace(payload)
	if strings.Contains(payload, "=") {
		parts := strings.Split(payload, ",")
		for _, p := range parts {
			kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
			if len(kv) != 2 {
				continue
			}
			v, e := strconv.Atoi(strings.TrimSpace(kv[1]))
			if e != nil {
				return 0, 0, e
			}
			switch strings.TrimSpace(kv[0]) {
			case "dx":
				dx = v
			case "dy":
				dy = v
			}
		}
		return dx, dy, nil
	}
	parts := strings.Split(payload, ",")
	if len(parts) != 2 {
		return 0, 0, errBad("want dx,dy")
	}
	dx, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	dy, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	return dx, dy, err
}

type simpleError string

func (e simpleError) Error() string { return string(e) }

func errBad(s string) error { return simpleError(s) }

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}
