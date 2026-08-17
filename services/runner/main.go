// GoLearning code runner: accepts small Go programs, runs them with a timeout, returns output.
// Intended for local / Docker learning use with resource limits — not a multi-tenant sandbox.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	defaultAddr     = ":8091"
	maxCodeBytes    = 64 * 1024
	maxOutputBytes  = 256 * 1024
	defaultTimeout  = 5 * time.Second
	maxTimeout      = 10 * time.Second
	maxConcurrent   = 4
	workdirParent   = "/tmp/golearning-runs"
)

var (
	sem = make(chan struct{}, maxConcurrent)
	// crude rate limit: per IP simple counter window
	rateMu   sync.Mutex
	rateMap  = map[string]*rateBucket{}
	maxPerMin = 30
)

type rateBucket struct {
	n     int
	reset time.Time
}

type runRequest struct {
	Code    string `json:"code"`
	Timeout int    `json:"timeout_sec,omitempty"` // 1..10
}

type runResponse struct {
	OK         bool   `json:"ok"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
	Truncated  bool   `json:"truncated,omitempty"`
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)

	addr := env("ADDR", defaultAddr)
	if err := os.MkdirAll(workdirParent, 0o755); err != nil {
		log.Error("mkdir workdir", "err", err)
		os.Exit(1)
	}

	// Warm module cache with a hello (optional, best-effort)
	_ = warmup()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "service": "golearning-runner"})
	})
	mux.HandleFunc("POST /run", handleRun)
	mux.HandleFunc("OPTIONS /run", func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		w.WriteHeader(http.StatusNoContent)
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Info("runner listening", "addr", addr, "timeout", defaultTimeout.String())
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("listen", "err", err)
		os.Exit(1)
	}
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	ip := clientIP(r)
	if !allow(ip) {
		writeJSON(w, http.StatusTooManyRequests, runResponse{OK: false, Error: "rate limit: too many runs, try again later"})
		return
	}

	var req runRequest
	r.Body = http.MaxBytesReader(w, r.Body, maxCodeBytes+1024)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, runResponse{OK: false, Error: "invalid JSON body"})
		return
	}
	code := strings.ReplaceAll(req.Code, "\r\n", "\n")
	if strings.TrimSpace(code) == "" {
		writeJSON(w, http.StatusBadRequest, runResponse{OK: false, Error: "code is empty"})
		return
	}
	if len(code) > maxCodeBytes {
		writeJSON(w, http.StatusBadRequest, runResponse{OK: false, Error: fmt.Sprintf("code too large (max %d bytes)", maxCodeBytes)})
		return
	}
	if err := validateCode(code); err != nil {
		writeJSON(w, http.StatusBadRequest, runResponse{OK: false, Error: err.Error()})
		return
	}

	timeout := defaultTimeout
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
		if timeout > maxTimeout {
			timeout = maxTimeout
		}
		if timeout < time.Second {
			timeout = time.Second
		}
	}

	select {
	case sem <- struct{}{}:
		defer func() { <-sem }()
	default:
		writeJSON(w, http.StatusServiceUnavailable, runResponse{OK: false, Error: "server busy, retry shortly"})
		return
	}

	res := runGo(r.Context(), code, timeout)
	status := http.StatusOK
	if res.Error != "" && res.Stdout == "" && res.Stderr == "" && res.ExitCode == 0 {
		status = http.StatusBadRequest
	}
	writeJSON(w, status, res)
}

func validateCode(code string) error {
	// Single-file teaching programs only.
	if !regexp.MustCompile(`(?m)^\s*package\s+main\b`).MatchString(code) {
		return errors.New("code must use package main")
	}
	// Block obvious attempts to touch the host / heavy imports for classroom safety.
	blocked := []string{
		"os/exec",
		"syscall",
		"unsafe",
		"net/http",
		"plugin",
		"golang.org/x/",
		"github.com/",
		"gopkg.in/",
		"C.", // cgo
	}
	lower := code
	for _, b := range blocked {
		if strings.Contains(lower, b) {
			return fmt.Errorf("import or feature not allowed in playground: %s", b)
		}
	}
	// Disallow go: directives that pull external tools
	if strings.Contains(code, "//go:generate") || strings.Contains(code, "import \"C\"") {
		return errors.New("go:generate and cgo are not allowed")
	}
	return nil
}

func runGo(parent context.Context, code string, timeout time.Duration) runResponse {
	start := time.Now()
	dir, err := os.MkdirTemp(workdirParent, "run-*")
	if err != nil {
		return runResponse{OK: false, Error: "cannot create temp dir", DurationMs: time.Since(start).Milliseconds()}
	}
	defer os.RemoveAll(dir)

	mainPath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainPath, []byte(code), 0o600); err != nil {
		return runResponse{OK: false, Error: "cannot write code", DurationMs: time.Since(start).Milliseconds()}
	}
	mod := "module playground\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o600); err != nil {
		return runResponse{OK: false, Error: "cannot write go.mod", DurationMs: time.Since(start).Milliseconds()}
	}

	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GO111MODULE=on",
		"CGO_ENABLED=0",
		// Runner may have no outbound network (Docker internal network).
		"GOPROXY=off",
		"GOSUMDB=off",
		"GOTOOLCHAIN=local",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &limitedWriter{buf: &stdout, limit: maxOutputBytes}
	cmd.Stderr = &limitedWriter{buf: &stderr, limit: maxOutputBytes}

	err = cmd.Run()
	dur := time.Since(start).Milliseconds()
	res := runResponse{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		DurationMs: dur,
		ExitCode:   0,
	}
	if lw, ok := cmd.Stdout.(*limitedWriter); ok && lw.truncated {
		res.Truncated = true
	}
	if lw, ok := cmd.Stderr.(*limitedWriter); ok && lw.truncated {
		res.Truncated = true
	}

	if ctx.Err() == context.DeadlineExceeded {
		res.OK = false
		res.Error = fmt.Sprintf("timeout after %s", timeout)
		res.ExitCode = -1
		if res.Stderr != "" {
			res.Stderr += "\n"
		}
		res.Stderr += res.Error
		return res
	}
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			res.ExitCode = ee.ExitCode()
		} else {
			res.ExitCode = 1
			if res.Stderr == "" {
				res.Stderr = err.Error()
			}
		}
		res.OK = false
		return res
	}
	res.OK = true
	return res
}

type limitedWriter struct {
	buf       *bytes.Buffer
	limit     int
	truncated bool
}

func (w *limitedWriter) Write(p []byte) (int, error) {
	remain := w.limit - w.buf.Len()
	if remain <= 0 {
		w.truncated = true
		return len(p), nil
	}
	if len(p) > remain {
		w.buf.Write(p[:remain])
		w.truncated = true
		return len(p), nil
	}
	return w.buf.Write(p)
}

func warmup() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "version")
	return cmd.Run()
}

func allow(ip string) bool {
	rateMu.Lock()
	defer rateMu.Unlock()
	now := time.Now()
	b, ok := rateMap[ip]
	if !ok || now.After(b.reset) {
		rateMap[ip] = &rateBucket{n: 1, reset: now.Add(time.Minute)}
		return true
	}
	if b.n >= maxPerMin {
		return false
	}
	b.n++
	return true
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i >= 0 {
		return host[:i]
	}
	return host
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
