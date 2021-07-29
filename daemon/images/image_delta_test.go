package images

import (
	"fmt"
	"math"
	"testing"
)

func Test_roundToPowerOf2(t *testing.T) {
	tests := []struct {
		x    uint32
		want uint32
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{11, 16},
		{127, 128},
		{128, 128},
		{129, 256},
		{1234, 2048},
		{15000, 16384},
		{30001, 32768},
		{44444, 65536},
		{4190000, 4194304},
		{1073711111, 1073741824},
		{1073741824, 1073741824},
		{uint32(math.Pow(2, 31)), uint32(math.Pow(2, 31))},
		{2147483649, 2147483648},
		{4294967295, 2147483648},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("roundToPowerOf2(%v)", tt.x), func(t *testing.T) {
			if got := roundToPowerOf2(tt.x); got != tt.want {
				t.Errorf("got roundToPowerOf2(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}
