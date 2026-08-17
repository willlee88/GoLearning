package envelope_test

import (
	"testing"

	"github.com/willyliao/golearning/examples/g01-envelope"
)

func TestParseMove(t *testing.T) {
	raw := []byte(`{"v":1,"type":"move","payload":{"dx":1,"dy":0}}`)
	env, err := envelope.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	m, err := envelope.ParseMove(env)
	if err != nil || m.DX != 1 {
		t.Fatalf("%+v %v", m, err)
	}
}

func TestBadVersion(t *testing.T) {
	_, err := envelope.Parse([]byte(`{"v":99,"type":"chat","payload":{}}`))
	if err == nil {
		t.Fatal("expected error")
	}
}
