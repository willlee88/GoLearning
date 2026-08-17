package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigAndSummarize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "players.json")
	content := `{
		"room_capacity": 2,
		"players": [
			{"id": 1, "name": "A", "score": 3},
			{"id": 2, "name": "B", "score": 9}
		]
	}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.RoomCapacity != 2 || len(cfg.Players) != 2 {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}

	total, top := summarize(cfg)
	if total != 12 {
		t.Fatalf("total=%d want 12", total)
	}
	if top.Name != "B" || top.Score != 9 {
		t.Fatalf("top=%+v", top)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := loadConfig(filepath.Join(t.TempDir(), "nope.json"))
	if err == nil {
		t.Fatal("expected error")
	}
}
