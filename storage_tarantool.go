package main

import "github.com/tarantool/go-tarantool"

const (
	ER_TUPLE_FOUND = 3
)

type StorageTarantool struct {
	conn  *tarantool.Connection
	space string
}

type tarantoolTuple struct {
	//nolint:structcheck,megacheck
	_msgpack struct{} `msgpack:",asArray"`
	ID       string
	Value    []byte
}

func NewStorageTarantool(host, user, password, space string) *StorageTarantool {
	opts := tarantool.Opts{
		User: user,
		Pass: password,
	}
	conn, err := tarantool.Connect(host, opts)
	if err != nil {
		panic(err)
	}
	_, err = conn.Ping()
	if err != nil {
		panic(err)
	}

	return &StorageTarantool{
		conn:  conn,
		space: space,
	}
}

func (s *StorageTarantool) Store(key, value []byte) error {
	tuple := tarantoolTuple{
		ID:    string(key),
		Value: value,
	}
	_, err := s.conn.Insert(s.space, tuple)
	if err != nil {
		if tarantoolErr, ok := err.(tarantool.Error); ok {
			if tarantoolErr.Code == ER_TUPLE_FOUND {
				err = errDuplicate
			}
		}
	}
	return err
}

func (s *StorageTarantool) Get(key []byte) (value []byte, err error) {
	var items []tarantoolTuple
	err = s.conn.SelectTyped(s.space, "primary", 0, 1,
		tarantool.IterEq, tarantool.StringKey{string(key)}, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errNoKey
	}
	return items[0].Value, nil
}

func (s *StorageTarantool) Close() error {
	return s.conn.Close()
}
