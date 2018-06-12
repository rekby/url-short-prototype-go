package main

import "encoding/base32"

type MakeUrlFunc func (prefix, id []byte)(url []byte)

var base32Encoding = base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ-2345679").WithPadding(base32.NoPadding)
func encodeUrlBase32(prefix, val []byte)[]byte{
	resLen := base32Encoding.EncodedLen(len(val))
	res := make([]byte, resLen + len(prefix))
	copy(res, prefix)
	base32Encoding.Encode(res[len(prefix):], val)
	return res
}

