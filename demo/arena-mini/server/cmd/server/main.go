// Arena Mini — authoritative match, metrics, slog, graceful shutdown.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/willyliao/golearning/demo/arena-mini/internal/hub"
	"github.com/willyliao/golearning/demo/arena-mini/internal/metrics"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)

	addr := env("ADDR", ":8080")
	m := metrics.New()
	h := hub.New(m, log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{
			"ok":      true,
			"service": "arena-mini",
			"version": "m5.1",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"rooms":   h.RoomCount(),
		})
	})
	mux.HandleFunc("GET /rooms", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"rooms": h.RoomStats()})
	})
	mux.HandleFunc("GET /metrics", m.Handler())
	mux.HandleFunc("GET /ws", h.HandleWS)

	webDir := env("WEB_DIR", "../web")
	mux.Handle("GET /", http.FileServer(http.Dir(webDir)))

	srv := &http.Server{
		Addr:              addr,
		Handler:           withCORS(withAccessLog(log, mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("arena-mini listening", "addr", addr, "web", webDir, "version", "m5.1")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen failed", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Info("shutdown signal received")

	// stop game loops first, then HTTP
	h.Close()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown", "err", err)
	}
	log.Info("bye")
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withAccessLog(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skip noisy static / high-frequency if needed
		if r.URL.Path == "/metrics" || r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		next.ServeHTTP(w, r)
		if r.URL.Path != "/" && r.Header.Get("Upgrade") == "" {
			log.Info("http", "method", r.Method, "path", r.URL.Path, "dur", time.Since(start).String())
		}
	})
}
