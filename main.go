package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"net/http"
)

const (
	DEFAULT_DIR_MODE = 0700
	DEFAULT_FILE_MODE = 0600
)

var (
	storage Storage
	hashFunc HashFunc
	makeUrl MakeUrlFunc
)

func main(){
	flag.Parse()
	urlPrefixBytes = []byte(*urlPrefix)

	storage = NewStorageFiles(*storeFolder)
	hashFunc = hashSha256
	makeUrl = encodeUrlBase32

	fasthttp.ListenAndServe(*bindAddress, handleRequest)
}

func handleRequest(ctx *fasthttp.RequestCtx){
	addrBytes := ctx.FormValue("url")
	if !checkUrl(addrBytes){
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	urlHash := hashFunc(addrBytes)
	id := makeUrl(urlPrefixBytes, urlHash)
	storage.Store(id, addrBytes)
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
