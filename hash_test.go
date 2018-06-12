package main

import (
	"bytes"
	"testing"
)

func TestSipHashes(t *testing.T) {
	const cnt = 10000
	_, values := createBenchData(cnt)
	for i := 0; i < cnt; i++ {
		dchestSipHash_48bitResult := dchestSipHash_48bit(values[i])
		dchestSipHash_48bitFastResult := dchestSipHash_48bitFast(values[i])
		aeadSipHash_48bitResult := aeadSipHash_48bit(values[i])
		if !bytes.Equal(dchestSipHash_48bitResult, dchestSipHash_48bitFastResult) ||
			!bytes.Equal(dchestSipHash_48bitFastResult, aeadSipHash_48bitResult) {
			t.Error("Hashes doesn't equal")
			t.Log(dchestSipHash_48bitResult)
			t.Log(dchestSipHash_48bitFastResult)
			t.Log(aeadSipHash_48bitResult)
		}
	}
}
