package random

import "testing"

func Test_random(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size = 1", 1},
		{"size = 5", 5},
		{"size = 10", 10},
		{"size = 20", 20},
		{"size = 30", 30},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result1 := NewRandomString(tc.size)
			result2 := NewRandomString(tc.size)
			if len(result1) != tc.size {
				t.Errorf("expected %d, got %d", tc.size, len(result1))
			}
			if len(result2) != tc.size {
				t.Errorf("expected %d, got %d", tc.size, len(result2))
			}
			if result1 != result2 {
				t.Errorf("two consecutive calls returned identical results %q", result1)
			}
		})
	}
}
