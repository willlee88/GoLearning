package double

import "testing"

func double(n int) int { return n * 2 }

func TestDouble(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"zero", 0, 0},
		{"pos", 3, 6},
		{"neg", -2, -4},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := double(tc.in); got != tc.want {
				t.Fatalf("got %d want %d", got, tc.want)
			}
		})
	}
}
