package main

import "crypto/sha256"
import dchest "github.com/dchest/siphash"
import (
	"crypto/md5"

	aead "github.com/aead/siphash"
)

type HashFunc func(value []byte) []byte

func hashMD5_48Bit(value []byte) []byte {
	hash := md5.Sum(value)
	return hash[:6]
}

func hashSha256(value []byte) []byte {
	hash := sha256.Sum256(value)
	return hash[:]
}

func hashSha256_48Bit(value []byte) []byte {
	hash := sha256.Sum256(value)
	return hash[:6]
}

func hashSha256_64Bit(value []byte) []byte {
	hash := sha256.Sum256(value)
	return hash[:8]
}

var sipHashKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

func dchestSipHash_48bit(value []byte) []byte {
	hash := dchest.New(sipHashKey)
	hash.Write(value)
	res := hash.Sum(nil)
	return res[:6]
}

var sipHashKeyUint0 = uint64(sipHashKey[0]) | uint64(sipHashKey[1])<<8 |
	uint64(sipHashKey[2])<<16 | uint64(sipHashKey[3])<<24 |
	uint64(sipHashKey[4])<<32 | uint64(sipHashKey[5])<<40 | uint64(sipHashKey[6])<<48 |
	uint64(sipHashKey[7])<<56
var sipHashKeyUint1 = uint64(sipHashKey[8]) | uint64(sipHashKey[9])<<8 |
	uint64(sipHashKey[10])<<16 | uint64(sipHashKey[11])<<24 |
	uint64(sipHashKey[12])<<32 | uint64(sipHashKey[13])<<40 | uint64(sipHashKey[14])<<48 |
	uint64(sipHashKey[15])<<56

func dchestSipHash_48bitFast(value []byte) []byte {
	resUint := dchest.Hash(sipHashKeyUint0, sipHashKeyUint1, value)
	res := make([]byte, 6)
	res[0] = byte(resUint)
	res[1] = byte(resUint >> 8)
	res[2] = byte(resUint >> 16)
	res[3] = byte(resUint >> 24)
	res[4] = byte(resUint >> 32)
	res[5] = byte(resUint >> 40)
	return res
}

func aeadSipHash_48bit(value []byte) []byte {
	hash, _ := aead.New64(sipHashKey)
	hash.Write(value)
	res := hash.Sum(nil)
	return res[:6]
}
