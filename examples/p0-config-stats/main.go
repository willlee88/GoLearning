// p0-config-stats 對應 P0.2 / P0.4：讀 JSON 設定、統計假玩家。
//
// 對照 Python 直覺：
//   with open(...) as f: data = json.load(f)
// 在 Go 是明確的 error 處理與靜態 struct。
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Config 是設定檔的靜態契約（對照 Python 的 dict + 執行期 KeyError）。
type Config struct {
	RoomCapacity int      `json:"room_capacity"`
	Players      []Player `json:"players"`
}

type Player struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}
	if cfg.RoomCapacity <= 0 {
		return Config{}, errors.New("room_capacity must be positive")
	}
	return cfg, nil
}

func summarize(cfg Config) (totalScore int, top Player) {
	if len(cfg.Players) == 0 {
		return 0, Player{}
	}
	top = cfg.Players[0]
	for _, p := range cfg.Players {
		totalScore += p.Score
		if p.Score > top.Score {
			top = p
		}
	}
	return totalScore, top
}

func main() {
	path := "players.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	cfg, err := loadConfig(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	total, top := summarize(cfg)
	fmt.Printf("players=%d capacity=%d total_score=%d\n", len(cfg.Players), cfg.RoomCapacity, total)
	fmt.Printf("top_player=%s score=%d\n", top.Name, top.Score)

	if len(cfg.Players) > cfg.RoomCapacity {
		fmt.Println("warning: player count exceeds room_capacity")
	}
}
