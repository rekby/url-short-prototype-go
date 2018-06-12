package main

import "crypto/sha256"

type HashFunc func(value []byte)[]byte

func hashSha256(value []byte)[]byte{
	hash :=  sha256.Sum256(value)
	return hash[:]
}

