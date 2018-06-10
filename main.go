package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"crypto/sha256"
	"os"
	"path/filepath"
	"encoding/base32"
	"io/ioutil"
	"net/http"
)

const (
	DEFAULT_DIR_MODE = 0700
	DEFAULT_FILE_MODE = 0600
)

func main(){
	flag.Parse()

	os.MkdirAll(*storeFolder, DEFAULT_DIR_MODE)

	fasthttp.ListenAndServe(*bindAddress, handleRequest)
}

func handleRequest(ctx *fasthttp.RequestCtx){
	addrBytes := ctx.FormValue("url")
	if !checkUrl(addrBytes){
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	urlHash := hashBytes(addrBytes)
	id := encodeBytesToUrl(urlHash)
	storeUrl(id, addrBytes)
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.Header.Set("Content-type", "text/plain")
	ctx.Write(id)
}

func checkUrl(url []byte)bool {
	if len(url) == 0 || len(url) > 3000{
		return false
	}
	return true
}

func hashBytes(value []byte)[]byte{
	hash :=  sha256.Sum256(value)
	return hash[:]
}

var base32Encoding = base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ-2345679").WithPadding(base32.NoPadding)
func encodeBytesToUrl(val []byte)[]byte{
	resLen := base32Encoding.EncodedLen(len(val))
	res := make([]byte, resLen)
	base32Encoding.Encode(res, val)
	return res
}

func storeUrl(id, value []byte){
	fileName := filepath.Join(*storeFolder, string(id) + ".txt")
	ioutil.WriteFile(fileName, value, DEFAULT_FILE_MODE)
}