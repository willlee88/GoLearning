package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type Metrics struct {
	Connections atomic.Int64 `json:"-"`
	MessagesIn  atomic.Int64 `json:"-"`
}

func (m *Metrics) snapshot() map[string]int64 {
	return map[string]int64{
		"connections": m.Connections.Load(),
		"messages_in": m.MessagesIn.Load(),
	}
}

func main() {
	var m Metrics
	mux := http.NewServeMux()
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(m.snapshot())
	})
	mux.HandleFunc("POST /hit", func(w http.ResponseWriter, r *http.Request) {
		m.MessagesIn.Add(1)
		w.WriteHeader(http.StatusNoContent)
	})
	// simulate connections
	m.Connections.Store(3)

	srv := &http.Server{Addr: ":8099", Handler: mux, ReadHeaderTimeout: 3 * time.Second}
	log.Println("metrics demo :8099")
	log.Fatal(srv.ListenAndServe())
}
