package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

const (
	TEST_REDIS_SERVER_NETWORK = "tcp"
	TEST_REDIS_ADDRESS        = "127.0.0.1:6379"
)

type TestSkip interface {
	Skip(args ...interface{})
}

func redisInit(t TestSkip) *StorageRedis {
	testDb, err := strconv.Atoi(os.Getenv("REKBY_REDIS_TEST_DB"))
	if err != nil {
		fmt.Print(`
For test redis env set value REKBY_REDIS_TEST_DB to number or redis DB
WARNING: The database will be flushed (REMOVE ALL DATA).
`)
		t.Skip(err)
	}
	s := NewStorageRedis(TEST_REDIS_SERVER_NETWORK, TEST_REDIS_ADDRESS,
		testDb)
	err = s.redisPool.Cmd("FLUSHDB").Err
	if err != nil {
		panic(err)
	}
	return s
}

func TestStorageRedis_Store(t *testing.T) {
	s := redisInit(t)
	err := s.Store([]byte("123"), []byte("234"))
	if err != nil {
		t.Error(err)
	}

	resp := s.redisPool.Cmd("GET", "123")
	str, err := resp.Str()
	if err != nil || str != "234" {
		t.Error(resp.Err, resp.String())
	}
}

func TestStorageRedis_StoreDuplicate(t *testing.T) {
	s := redisInit(t)
	s.Store([]byte("123"), []byte("234"))
	err := s.Store([]byte("123"), []byte("234"))
	if err != errDuplicate {
		t.Error(err)
	}

	resp := s.redisPool.Cmd("GET", "123")
	str, err := resp.Str()
	if err != nil || str != "234" {
		t.Error(resp.Err, resp.String())
	}
}

func TestStorageRedis_Get(t *testing.T) {
	s := redisInit(t)
	s.redisPool.Cmd("SET", "234", "567")
	val, err := s.Get([]byte("234"))
	if err != nil || string(val) != "567" {
		t.Error(err, string(val))
	}
}

func TestStorageRedis_GetNoKey(t *testing.T) {
	s := redisInit(t)
	val, err := s.Get([]byte("234"))
	if err != errNoKey {
		t.Error(err, string(val))
	}
}
