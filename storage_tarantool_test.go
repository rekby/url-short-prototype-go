package main

import (
	"testing"

	"github.com/tarantool/go-tarantool"
)

const (
	TEST_TARANTOOL_SERVER   = "localhost:3301"
	TEST_TARANTOOL_USER     = "admin"
	TEST_TARANTOOL_PASSWORD = ""
	TEST_TARANTOOL_SPACE    = "test"
)

func tarantoolTestInit() *StorageTarantool {
	s := NewStorageTarantool(TEST_TARANTOOL_SERVER, TEST_TARANTOOL_USER, TEST_TARANTOOL_PASSWORD, TEST_TARANTOOL_SPACE)
	if space, exist := s.conn.Schema.Spaces[TEST_TARANTOOL_SPACE]; exist {
		_, err := s.conn.Call("box.schema.space.drop", []interface{}{space.Id})
		if err != nil {
			panic(err)
		}
	}

	s.conn.Call("box.schema.space.create", []interface{}{TEST_TARANTOOL_SPACE})

	createIndexTuple := struct {
		Type  string        `msgpack:"type"`
		Parts []interface{} `msgpack:"parts"`
	}{
		Type:  "hash",
		Parts: []interface{}{1, "string"},
	}

	_, err := s.conn.Call("box.space."+TEST_TARANTOOL_SPACE+":create_index", []interface{}{"primary", createIndexTuple})
	if err != nil {
		panic(err)
	}

	s.Close()

	s = NewStorageTarantool(TEST_TARANTOOL_SERVER, TEST_TARANTOOL_USER, TEST_TARANTOOL_PASSWORD, TEST_TARANTOOL_SPACE)
	return s
}

func TestStorageTarantool_Store(t *testing.T) {
	s := tarantoolTestInit()
	err := s.Store([]byte("123"), []byte("234"))
	if err != nil {
		t.Error(err)
	}

	var v []tarantoolTuple
	err = s.conn.SelectTyped(TEST_TARANTOOL_SPACE, "primary", 0, 1, tarantool.IterEq, tarantool.StringKey{"123"}, &v)
	if err != nil || v[0].ID != "123" || string(v[0].Value) != "234" {
		t.Error(err, v[0].ID, string(v[0].Value))
	}
}

func TestStorageTarantool_StoreDuplicate(t *testing.T) {
	s := tarantoolTestInit()
	s.Store([]byte("123"), []byte("234"))
	err := s.Store([]byte("123"), []byte("234"))
	if err != errDuplicate {
		t.Error(err)
	}

	var v []tarantoolTuple
	err = s.conn.SelectTyped(TEST_TARANTOOL_SPACE, "primary", 0, 1, tarantool.IterEq, tarantool.StringKey{"123"}, &v)
	if err != nil || v[0].ID != "123" || string(v[0].Value) != "234" {
		t.Error(err, v[0].ID, string(v[0].Value))
	}
}

func TestStorageTarantool_Get(t *testing.T) {
	s := tarantoolTestInit()
	s.conn.Insert(TEST_TARANTOOL_SPACE, []interface{}{"222", "123"})
	val, err := s.Get([]byte("222"))
	if err != nil || string(val) != "123" {
		t.Error(err, string(val))
	}
}

func TestStorageTarantool_GetNoKey(t *testing.T) {
	s := tarantoolTestInit()
	val, err := s.Get([]byte("222"))
	if err != errNoKey {
		t.Error(err, string(val))
	}
}
