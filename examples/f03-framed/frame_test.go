package framed_test

import (
	"bytes"
	"errors"
	"testing"

	framed "github.com/willyliao/golearning/examples/f03-framed"
)

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	payload := []byte(`{"type":"chat","v":1}`)
	if err := framed.WriteFrame(&buf, payload); err != nil {
		t.Fatal(err)
	}
	got, err := framed.ReadFrame(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("got %s", got)
	}
}

func TestTooLarge(t *testing.T) {
	big := make([]byte, framed.MaxFrame+1)
	err := framed.WriteFrame(ioDiscard{}, big)
	if !errors.Is(err, framed.ErrFrameTooLarge) {
		t.Fatalf("err=%v", err)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
