package main

import (
	"testing"
	"strconv"
	"math/rand"
	"io/ioutil"
	"os"
)

func BenchmarkStorageFiles_Store(b *testing.B){
	keys := make([][]byte, b.N)
	val := []byte("testVal")

	r := rand.New(rand.NewSource(123))

	for i := 0; i < b.N; i++{
		keys[i] = []byte(strconv.Itoa(r.Int()))
	}

	tmpDir, err := ioutil.TempDir("", "benchmark-files")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	s := NewStorageFiles(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Store(keys[i], val)
	}

	b.StopTimer()
}