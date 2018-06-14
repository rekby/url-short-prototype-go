package main

import "testing"

var benchmarkBytesForHash = []byte("https://yandex.ru/search/?text=%D0%BA%D0%B0%D0%BA%D0%B8%D0%B5%20%D0%BD%D0%BE%D0%B2%D0%BE%D1%81%D1%82%D0%B8%20%D0%BD%D0%B0%20%D1%81%D0%B5%D0%B3%D0%BE%D0%B4%D0%BD%D1%8F%3F&lr=213")

func BenchmarkHashMD5_48Bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashMD5_48Bit(benchmarkBytesForHash)
	}
}

func BenchmarkHashSha256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashSha256(benchmarkBytesForHash)
	}
}

func BenchmarkHashSipDchest_48bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashSipDchest_48bit(benchmarkBytesForHash)
	}
}

func BenchmarkHashSipDchestFast_48bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashSipDchestFast_48bit(benchmarkBytesForHash)
	}
}

func BenchmarkHashSipAeadSip_48bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashSipAead_48bit(benchmarkBytesForHash)
	}
}

func BenchmarkHashRandom_48bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashRandom_48Bit(benchmarkBytesForHash)
	}
}

func BenchmarkHashRandomCrypto_48Bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hashRandomCrypto_48Bit(benchmarkBytesForHash)
	}
}
