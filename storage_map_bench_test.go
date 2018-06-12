package main

import (
	"testing"
)

func BenchmarkStorageMap_Store(b *testing.B) {
	keys, vals := createBenchData(b.N)
	s := NewStorageMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Store(keys[i], vals[i])
	}
}
