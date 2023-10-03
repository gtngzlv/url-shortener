package util

import (
	"testing"
)

func BenchmarkRandStringRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandStringRunes()
	}
}
