package main

import (
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkStorageMap_Store(b *testing.B) {
	keys := make([][]byte, b.N)
	val := []byte("testVal")

	r := rand.New(rand.NewSource(123))

	for i := 0; i < b.N; i++ {
		keys[i] = []byte(strconv.Itoa(r.Int()))
	}

	s := NewStorageMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Store(keys[i], val)
	}
}
