package main

import "errors"

var (
	errNoKey = errors.New("Key doesn't exist")
	errDuplicate = errors.New("Key duplication")
)

type Storage interface {
	Store (key, value []byte)error
	Get(key []byte)(value []byte, err error)
}