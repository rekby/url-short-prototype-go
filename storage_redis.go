package main

import (
	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type StorageRedis struct {
	redisPool *pool.Pool
}

func NewStorageRedis(network, address string, database int) *StorageRedis {
	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		databaseString := strconv.Itoa(database)
		if err = client.Cmd("SELECT", databaseString).Err; err != nil {
			client.Close()
			return nil, err
		}
		return client, nil

	}
	redisPool, err := pool.NewCustom(network, address, 10, df)
	if err != nil {
		panic(err)
	}
	resp := redisPool.Cmd("PING")
	if resp.Err != nil {
		panic(err)
	}
	return &StorageRedis{
		redisPool: redisPool,
	}
}

func (s *StorageRedis) Store(key, value []byte) error {
	resp := s.redisPool.Cmd("SET", key, value, "NX")
	err := resp.Err
	if resp.IsType(redis.Nil) {
		return errDuplicate
	}
	return err
}

func (s *StorageRedis) Get(key []byte) (value []byte, err error) {
	resp := s.redisPool.Cmd("GET", key)
	if resp.IsType(redis.Nil) {
		return nil, errNoKey
	}
	return resp.Bytes()
}
