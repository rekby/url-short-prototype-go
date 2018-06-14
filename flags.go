package main

import "flag"

var (
	bindAddress    = flag.String("bind", ":8080", "Bind address for http handler")
	storeFolder    = flag.String("store-folder", "_storage", "path to storage folder")
	urlPrefix      = flag.String("url-prefix", "http://localhost:8080/", "Url prefix before id")
	urlPrefixBytes []byte
	maxRetryCount  = flag.Int("max-retry-save", 100, "Max count for save hash on any error")

	storageType = flag.String("storage-type", "files", "files|memory-map|redis|tarantool")

	redisAddress  = flag.String("redis-addr", "127.0.0.1:6379", "redis addr")
	redisDatabase = flag.Int("redis-database", 0, "")

	tarantoolServer   = flag.String("tarantool-server", "127.0.0.1:3301", "")
	tarantoolUser     = flag.String("tarantool-user", "admin", "")
	tarantoolPassword = flag.String("tarantool-password", "", "")
	tarantoolSpace    = flag.String("tarantool-space", "url-short",
		"Space have to be existed. In space have to be existed primary index for first field, type scalar.")
)
