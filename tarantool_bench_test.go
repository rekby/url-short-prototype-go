package main

import (
	"runtime"
	"sync"
	"testing"
)

func BenchmarkStorageTarantool_Store(b *testing.B) {
	defer func() {
		err := recover()
		if err != nil {
			b.Skip(err)
		}
	}()
	tarantoolTestInit().Close()

	var goroutinesCount = benchmarkParalellism * runtime.GOMAXPROCS(-1)
	connections := make([]*StorageTarantool, goroutinesCount)
	keys := make([][][]byte, goroutinesCount)
	vals := make([][][]byte, goroutinesCount)
	var localMutex sync.Mutex

	for i := 0; i < goroutinesCount; i++ {
		connections[i] = NewStorageTarantool(TEST_TARANTOOL_SERVER, TEST_TARANTOOL_USER, TEST_TARANTOOL_PASSWORD, TEST_TARANTOOL_SPACE)
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
