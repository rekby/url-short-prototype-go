package main

import (
	"bytes"
	"flag"

	"net/http"

	"github.com/valyala/fasthttp"
)

const (
	DEFAULT_DIR_MODE  = 0700
	DEFAULT_FILE_MODE = 0600
)

var (
	storage  Storage
	hashFunc HashFunc
	makeUrl  MakeUrlFunc
)

func main() {
	flag.Parse()
	urlPrefixBytes = []byte(*urlPrefix)

	storage = NewStorageFiles(*storeFolder)
	hashFunc = dchestSipHash_48bitFast
	makeUrl = encodeUrlBase64

	fasthttp.ListenAndServe(*bindAddress, handleRequest)
}

func handleRequest(ctx *fasthttp.RequestCtx) {
	addrBytes := ctx.FormValue("url")
	if len(addrBytes) > 0 {
		handlreStoreRequest(ctx, addrBytes)
	} else {
		handleRedirectRequest(ctx)
	}
}

func handleRedirectRequest(ctx *fasthttp.RequestCtx) {
	ctx.Request.RequestURI()
}

func handlreStoreRequest(ctx *fasthttp.RequestCtx, urlBytes []byte) {
	if !checkUrl(urlBytes) {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	urlHash := hashFunc(urlBytes)
	id := makeUrl(urlPrefixBytes, urlHash)
	storage.Store(id, urlBytes)
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.Header.Set("Content-type", "text/plain")
	ctx.Write(id)
}

var checkUrlAllowedPrefixes = [][]byte{
	[]byte("http://"),
	[]byte("https://"),
	[]byte("ftp://"),
}

func checkUrl(url []byte) bool {
	if len(url) == 0 || len(url) > 3000 {
		return false
	}

	hasAllowedPrefix := false
	for _, prefix := range checkUrlAllowedPrefixes {
		if bytes.HasPrefix(url, prefix) {
			hasAllowedPrefix = true
			break
		}
	}
	if !hasAllowedPrefix {
		return false
	}

	return true
}
