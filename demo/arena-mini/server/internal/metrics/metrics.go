package metrics

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Snapshot is process-level gauges/counters for teaching (/metrics JSON).
type Metrics struct {
	StartedAt time.Time

	Connections   atomic.Int64
	Rooms         atomic.Int64
	MessagesIn    atomic.Int64
	MessagesOut   atomic.Int64
	InputsApplied atomic.Int64
	Errors        atomic.Int64
}

func New() *Metrics {
	return &Metrics{StartedAt: time.Now().UTC()}
}

func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service":         "arena-mini",
			"version":         "m5.1",
			"uptime_sec":      int(time.Since(m.StartedAt).Seconds()),
			"connections":     m.Connections.Load(),
			"rooms":           m.Rooms.Load(),
			"messages_in":     m.MessagesIn.Load(),
			"messages_out":    m.MessagesOut.Load(),
			"inputs_applied":  m.InputsApplied.Load(),
			"errors":          m.Errors.Load(),
		})
	}
}
