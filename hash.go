package main

import "crypto/sha256"
import dchest "github.com/dchest/siphash"
import (
	//nolint:gas
	"crypto/md5"

	cryptorand "crypto/rand"

	"math/rand"

	aead "github.com/aead/siphash"
)

type HashFunc func(value []byte) []byte

func hashMD5_48Bit(value []byte) []byte {
	//nolint:gas
	hash := md5.Sum(value)
	return hash[:6]
}

func hashSha256(value []byte) []byte {
	hash := sha256.Sum256(value)
	return hash[:]
}

//nolint:deadcode,megacheck
func hashSha256_48Bit(value []byte) []byte {
	hash := sha256.Sum256(value)
	return hash[:6]
}

//nolint:deadcode,megacheck
func hashSha256_64Bit(value []byte) []byte {
	hash := sha256.Sum256(value)
	return hash[:8]
}

var sipHashKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

func hashSipDchest_48bit(value []byte) []byte {
	hash := dchest.New(sipHashKey)
	//nolint:errcheck
	hash.Write(value)
	res := hash.Sum(nil)
	return res[:6]
}

var hashSipKeyUint0 = uint64(sipHashKey[0]) | uint64(sipHashKey[1])<<8 |
	uint64(sipHashKey[2])<<16 | uint64(sipHashKey[3])<<24 |
	uint64(sipHashKey[4])<<32 | uint64(sipHashKey[5])<<40 | uint64(sipHashKey[6])<<48 |
	uint64(sipHashKey[7])<<56
var hashSipKeyUint1 = uint64(sipHashKey[8]) | uint64(sipHashKey[9])<<8 |
	uint64(sipHashKey[10])<<16 | uint64(sipHashKey[11])<<24 |
	uint64(sipHashKey[12])<<32 | uint64(sipHashKey[13])<<40 | uint64(sipHashKey[14])<<48 |
	uint64(sipHashKey[15])<<56

func hashSipDchestFast_48bit(value []byte) []byte {
	resUint := dchest.Hash(hashSipKeyUint0, hashSipKeyUint1, value)
	res := make([]byte, 6)
	res[0] = byte(resUint)
	res[1] = byte(resUint >> 8)
	res[2] = byte(resUint >> 16)
	res[3] = byte(resUint >> 24)
	res[4] = byte(resUint >> 32)
	res[5] = byte(resUint >> 40)
	return res
}

func hashSipAead_48bit(value []byte) []byte {
	hash, _ := aead.New64(sipHashKey)
	//nolint:errcheck
	hash.Write(value)
	res := hash.Sum(nil)
	return res[:6]
}

func hashRandomCrypto_48Bit([]byte) []byte {
	res := make([]byte, 6)
	if _, err := cryptorand.Read(res); err != nil {
		panic(err)
	}
	return res
}

func hashRandom_48Bit([]byte) []byte {
	res := make([]byte, 6)
	//nolint:errcheck,gas
	rand.Read(res)
	return res
}
