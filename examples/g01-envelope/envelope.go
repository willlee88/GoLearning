package envelope

import (
	"encoding/json"
	"fmt"
)

type Envelope struct {
	V       int             `json:"v"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type MovePayload struct {
	DX int `json:"dx"`
	DY int `json:"dy"`
}

func Parse(data []byte) (Envelope, error) {
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return Envelope{}, err
	}
	if env.V != 1 {
		return Envelope{}, fmt.Errorf("unsupported version %d", env.V)
	}
	if env.Type == "" {
		return Envelope{}, fmt.Errorf("missing type")
	}
	return env, nil
}

func ParseMove(env Envelope) (MovePayload, error) {
	if env.Type != "move" {
		return MovePayload{}, fmt.Errorf("not move")
	}
	var m MovePayload
	if err := json.Unmarshal(env.Payload, &m); err != nil {
		return MovePayload{}, err
	}
	if m.DX < -1 || m.DX > 1 || m.DY < -1 || m.DY > 1 {
		return MovePayload{}, fmt.Errorf("move out of range")
	}
	return m, nil
}
