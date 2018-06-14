package main

import (
	"runtime"
	"sync"
	"testing"
)

//nolint:deadcode,megacheck
func BenchmarkStorageRedis_Store(b *testing.B) {
	var goroutinesCount = benchmarkParalellism * runtime.GOMAXPROCS(-1)
	connections := make([]*StorageRedis, goroutinesCount)
	keys := make([][][]byte, goroutinesCount)
	vals := make([][][]byte, goroutinesCount)
	var localMutex sync.Mutex

	for i := 0; i < goroutinesCount; i++ {
		connections[i] = redisInit(b)
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
			if err := s.Store(localKeys[0], localVals[0]); err != nil {
				b.Error(err)
			}
			localKeys, localVals = localKeys[1:], localVals[1:]
		}
	})
}
