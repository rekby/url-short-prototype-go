package main

import (
	"runtime"
	"sync"
	"testing"
)

func BenchmarkStorageMap_Store(b *testing.B) {
	var goroutinesCount = benchmarkParalellism * runtime.GOMAXPROCS(-1)
	connections := make([]*StorageMap, goroutinesCount)
	keys := make([][][]byte, goroutinesCount)
	vals := make([][][]byte, goroutinesCount)
	var localMutex sync.Mutex

	for i := 0; i < goroutinesCount; i++ {
		connections[i] = NewStorageMap()
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
