package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

const benchmarkParalellism = 4

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

//nolint:deadcode,megacheck,errcheck
func BenchmarkStorageFiles_Store(b *testing.B) {
	tmpDir, err := ioutil.TempDir("", "benchmark-files")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	var goroutinesCount = benchmarkParalellism * runtime.GOMAXPROCS(-1)
	connections := make([]StorageFiles, goroutinesCount)
	keys := make([][][]byte, goroutinesCount)
	vals := make([][][]byte, goroutinesCount)
	var localMutex sync.Mutex

	for i := 0; i < goroutinesCount; i++ {
		connections[i] = NewStorageFiles(tmpDir)
		keys[i], vals[i] = createBenchData(b.N)
	}

	b.SetParallelism(benchmarkParalellism)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		localMutex.Lock()
		s := connections[0]
		localKeys, localVals := keys[len(connections)-1], vals[len(connections)-1]
		connections = connections[1:]
		localMutex.Unlock()

		for pb.Next() {
			s.Store(localKeys[0], localVals[0])
			localKeys, localVals = localKeys[1:], localVals[1:]
		}
	})
}
