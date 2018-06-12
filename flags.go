package main

import "flag"

var (
	bindAddress           = flag.String("bind", ":8080", "Bind address for http handler")
	storeFolder           = flag.String("store-folder", "_storage", "path to storage folder")
	urlPrefix             = flag.String("url-prefix", "http://localhost:8080/", "Url prefix before id")
	urlPrefixBytes        []byte
	maxRetryCount         = flag.Int("max-retry-save", 100, "Max count for save hash on any error")
	addRandomBytesOnRetry = flag.Int("add-random-bytes", 8, "Add random bytes to url for generate new id")
)
