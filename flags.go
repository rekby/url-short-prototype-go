package main

import "flag"

var (
	bindAddress = flag.String("bind", ":8080", "Bind address for http handler")
	storeFolder = flag.String("store-folder", "_storage", "path to storage folder")
)