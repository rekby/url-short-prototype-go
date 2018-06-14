package main

import (
	"testing"
)

var (
	_ Storage = NewStorageMap()
)

//nolint:deadcode,megacheck
func TestStorageMap_Store(t *testing.T) {
	s := NewStorageMap()
	err := s.Store([]byte("123"), []byte("222"))
	if err != nil {
		t.Error(err)
	}

	val, exist := s.m["123"]
	if !exist || string(val) != "222" {
		t.Error(exist, string(val))
	}
}

//nolint:deadcode,megacheck,errcheck
func TestStorageMap_StoreDuplicate(t *testing.T) {
	s := NewStorageMap()
	s.Store([]byte("123"), []byte("222"))
	err := s.Store([]byte("123"), []byte("222"))
	if err != errDuplicate {
		t.Error(err)
	}
}

//nolint:deadcode,megacheck
func TestStorageMap_Get(t *testing.T) {
	s := NewStorageMap()
	s.m["123"] = []byte("222")

	val, err := s.Get([]byte("123"))
	if err != nil || string(val) != "222" {
		t.Error(err, val)
	}
}

//nolint:deadcode,megacheck
func TestStorageMap_GetNoKey(t *testing.T) {
	s := NewStorageMap()
	val, err := s.Get([]byte("123"))
	if err != errNoKey || val != nil {
		t.Error(err, val)
	}
}
