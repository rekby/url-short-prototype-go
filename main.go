package main

import (
	"bytes"
	"flag"

	"net/http"

	"net/url"

	"math/rand"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	DEFAULT_DIR_MODE  = 0700
	DEFAULT_FILE_MODE = 0600
)

var (
	storage         Storage     = nil
	hashFunc        HashFunc    = dchestSipHash_48bitFast
	makeUrl         MakeUrlFunc = encodeUrlBase64
	hashDecoderFunc IdDecoder   = decodeUrlBase64
)

func main() {
	flag.Parse()
	urlPrefixBytes = []byte(*urlPrefix)
	rand.Seed(time.Now().UnixNano())

	storage = NewStorageFiles(*storeFolder)

	fasthttp.ListenAndServe(*bindAddress, handleRequest)
}

func handleRequest(ctx *fasthttp.RequestCtx) {
	addrBytes := ctx.FormValue("url")
	if len(addrBytes) > 0 {
		handlreStoreRequest(ctx, addrBytes)
	} else {
		handleReadRequest(ctx)
	}
}

func handleReadRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain")
	encodedId := ctx.Request.RequestURI()
	if len(encodedId) < 2 {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	encodedId = encodedId[1:]

	binaryId, err := hashDecoderFunc(encodedId)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	destUrl, err := storage.Get(binaryId)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(destUrl)
}

func handlreStoreRequest(ctx *fasthttp.RequestCtx, urlBytes []byte) {
	ctx.SetContentType("text/plain")
	if !checkUrl(urlBytes) {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	bytesForHash := urlBytes
	var resultUrl []byte
	var saveErr error
	for tryIndex := 0; tryIndex < *maxRetryCount; tryIndex++ {
		urlHash := hashFunc(bytesForHash)
		saveErr = storage.Store(urlHash, urlBytes)
		if saveErr == nil {
			resultUrl = makeUrl(urlPrefixBytes, urlHash)
			break
		}

		if tryIndex == 0 {
			bytesForHash = make([]byte, len(urlBytes), len(urlBytes)+*maxRetryCount**addRandomBytesOnRetry)
			copy(bytesForHash, urlBytes)
		}
		buf := make([]byte, *addRandomBytesOnRetry)
		rand.Read(buf)
		bytesForHash = append(bytesForHash, buf...)
	}
	if saveErr != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.WriteString(saveErr.Error())
		return
	}
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.SetContentType("text/plain")
	ctx.Write(resultUrl)
}

var checkUrlAllowedPrefixes = [][]byte{
	[]byte("http://"),
	[]byte("https://"),
	[]byte("ftp://"),
}

func checkUrl(urlBytes []byte) bool {
	if len(urlBytes) == 0 || len(urlBytes) > 3000 {
		return false
	}

	hasAllowedPrefix := false
	for _, prefix := range checkUrlAllowedPrefixes {
		if bytes.HasPrefix(urlBytes, prefix) {
			hasAllowedPrefix = true
			break
		}
	}
	if !hasAllowedPrefix {
		return false
	}

	parsedUrl, err := url.Parse(string(urlBytes))
	if err != nil {
		return false
	}

	if parsedUrl.Host == "" || parsedUrl.Scheme == "" {
		return false
	}

	return true
}
