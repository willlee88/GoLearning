package parse

import "testing"

func TestParseOK(t *testing.T) {
	dx, dy, err := ParseDXDY("1,0")
	if err != nil || dx != 1 || dy != 0 {
		t.Fatalf("%d %d %v", dx, dy, err)
	}
}

func FuzzParseDXDY(f *testing.F) {
	f.Add("1,0")
	f.Add("0,0")
	f.Add("-1,1")
	f.Fuzz(func(t *testing.T, s string) {
		_, _, _ = ParseDXDY(s) // must not panic
	})
}
