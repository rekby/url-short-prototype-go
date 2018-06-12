package main

import "testing"

var benchmarkBytesForHash = []byte("12345678")

func BenchmarkMD5_48Bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashMD5_48Bit(benchmarkBytesForHash)
	}
}

func BenchmarkSha256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashSha256(benchmarkBytesForHash)
	}
}

func BenchmarkDchestSipHash_48bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dchestSipHash_48bit(benchmarkBytesForHash)
	}
}

func BenchmarkDchestSipHash_48bitFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dchestSipHash_48bitFast(benchmarkBytesForHash)
	}
}

func BenchmarkAeadSipHash_48bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		aeadSipHash_48bit(benchmarkBytesForHash)
	}
}
