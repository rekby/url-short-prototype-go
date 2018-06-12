package main

import (
	"encoding/base32"
	"math/big"
	"encoding/base64"
)

type MakeUrlFunc func (prefix, id []byte)(url []byte)

var base32Encoding = base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ-2345679").WithPadding(base32.NoPadding)
func encodeUrlBase32(prefix, val []byte)[]byte{
	resLen := base32Encoding.EncodedLen(len(val))
	res := make([]byte, resLen + len(prefix))
	copy(res, prefix)
	base32Encoding.Encode(res[len(prefix):], val)
	return res
}

func encodeUrlBase64(prefix, val []byte)[]byte{
	resLen := base64.RawURLEncoding.EncodedLen(len(val))
	res := make([]byte, resLen + len(prefix))
	copy(res, prefix)
	base64.RawURLEncoding.Encode(res[len(prefix):], val)
	return res
}

func encodeUrlBase62(prefix, val []byte)[]byte{
	bigInt := big.Int{}
	bigInt.SetBytes(val)
	return bigInt.Append(prefix, 62)
}