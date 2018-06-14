package main

import (
	"bytes"
	cryptorand "crypto/rand"
	"flag"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/valyala/fasthttp"

	"log"
	"math"
	"math/big"
)

const (
	DEFAULT_DIR_MODE  = 0700
	DEFAULT_FILE_MODE = 0600
)

var (
	storage         Storage     = nil
	hashFunc        HashFunc    = hashRandom_48Bit
	makeUrl         MakeUrlFunc = encodeUrlBase64
	hashDecoderFunc IdDecoder   = decodeUrlBase64
)

func main() {
	flag.Parse()
	urlPrefixBytes = []byte(*urlPrefix)
	randIntSeed, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	rand.Seed(randIntSeed.Int64())

	switch *storageType {
	case "files":
		storage = NewStorageFiles(*storeFolder)
	case "memory-map":
		storage = NewStorageMap()
	case "tarantool":
		t := NewStorageTarantool(*tarantoolServer, *tarantoolUser, *tarantoolPassword, *tarantoolSpace)
		defer t.Close()
		storage = t
	case "redis":
		storage = NewStorageRedis("tcp", *redisAddress, *redisDatabase)
	default:
		log.Fatalf("Unknown type of storage: '%v'", *storageType)
	}

	if err := fasthttp.ListenAndServe(*bindAddress, handleRequest); err != nil {
		log.Println(err)
	}
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
		if _, err := ctx.WriteString(err.Error()); err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
		}
		return
	}

	destUrl, err := storage.Get(binaryId)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		if _, err := ctx.WriteString(err.Error()); err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}
	ctx.SetStatusCode(http.StatusOK)
	if _, err := ctx.Write(destUrl); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
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

		bytesForHash = urlHash
	}
	if saveErr != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		if _, err := ctx.WriteString(saveErr.Error()); err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		return
	}
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.SetContentType("text/plain")
	if _, err := ctx.Write(resultUrl); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
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
