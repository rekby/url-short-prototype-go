package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	_ Storage = StorageFiles{}
)

//nolint:deadcode,megacheck,errcheck
func TestStorageFiles_Store(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "url-short")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	s := NewStorageFiles(tmpDir)
	err = s.Store([]byte("123"), []byte("222"))
	if err != nil {
		t.Error(err)
	}

	val, err := ioutil.ReadFile(filepath.Join(tmpDir, string(makeUrl(nil, []byte("123")))+".txt"))
	if err != nil || string(val) != "222" {
		t.Error(err, string(val))
	}
}

//nolint:deadcode,megacheck,errcheck
func TestStorageFiles_StoreDuplicate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "url-short")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	s := NewStorageFiles(tmpDir)
	s.Store([]byte("123"), []byte("222"))
	err = s.Store([]byte("123"), []byte("asdasd"))
	if err != errDuplicate {
		t.Error(err)
	}
}

//nolint:deadcode,megacheck,errcheck
func TestStorageFiles_Get(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "url-short")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	s := NewStorageFiles(tmpDir)
	ioutil.WriteFile(filepath.Join(tmpDir, string(makeUrl(nil, []byte("222")))+".txt"), []byte("234"), DEFAULT_FILE_MODE)
	value, err := s.Get([]byte("222"))
	if string(value) != "234" || err != nil {
		t.Error(err, value)
	}
}

//nolint:deadcode,megacheck,errcheck
func TestStorageFiles_GetNoKey(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "url-short")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	s := NewStorageFiles(tmpDir)
	value, err := s.Get([]byte("222"))
	if err != errNoKey || value != nil {
		t.Error(err, value)
	}
}
