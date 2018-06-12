package main

import (
	"encoding/base32"
	"encoding/base64"
	"math/big"
)

type MakeUrlFunc func(prefix, id []byte) (url []byte)
type IdDecoder func(urlHash []byte) ([]byte, error)

var base32Encoding = base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ-2345679").WithPadding(base32.NoPadding)

func encodeUrlBase32(prefix, val []byte) []byte {
	resLen := base32Encoding.EncodedLen(len(val))
	res := make([]byte, resLen+len(prefix))
	copy(res, prefix)
	base32Encoding.Encode(res[len(prefix):], val)
	return res
}

func encodeUrlBase64(prefix, val []byte) []byte {
	resLen := base64.RawURLEncoding.EncodedLen(len(val))
	res := make([]byte, resLen+len(prefix))
	copy(res, prefix)
	base64.RawURLEncoding.Encode(res[len(prefix):], val)
	return res
}

func decodeUrlBase64(val []byte) ([]byte, error) {
	maxLen := base64.RawURLEncoding.DecodedLen(len(val))
	res := make([]byte, maxLen)
	realLen, err := base64.RawURLEncoding.Decode(res, val)
	if err != nil {
		return nil, err
	}
	return res[:realLen], nil
}

func encodeUrlBase62(prefix, val []byte) []byte {
	bigInt := big.Int{}
	bigInt.SetBytes(val)
	return bigInt.Append(prefix, 62)
}
