package main

import "testing"

func BenchmarkSha256_8chars(b *testing.B){
	val := []byte("12345678")
	for i := 0; i < b.N; i++{
		hashSha256(val)
	}
}

func BenchmarkSha256_32chars(b *testing.B){
	val := []byte("12345678901234567890123456789012")
	for i := 0; i < b.N; i++{
		hashSha256(val)
	}
}
