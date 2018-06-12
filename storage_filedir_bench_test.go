package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func createBenchData(count int) (items, values [][]byte) {
	keys := make([][]byte, count)
	vals := make([][]byte, count)

	r := rand.New(rand.NewSource(123))

	for i := 0; i < count; i++ {
		keys[i] = []byte(strconv.Itoa(r.Int()))
		vals[i] = []byte("test-" + strconv.Itoa(r.Int()) + "-value")
	}
	return keys, vals
}

func BenchmarkStorageFiles_Store(b *testing.B) {
	keys, vals := createBenchData(b.N)
	tmpDir, err := ioutil.TempDir("", "benchmark-files")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	s := NewStorageFiles(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Store(keys[i], vals[i])
	}

	b.StopTimer()
}
