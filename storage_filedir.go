package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type StorageFiles struct {
	Dir string
}

func NewStorageFiles(dir string) StorageFiles {
	os.MkdirAll(dir, DEFAULT_DIR_MODE)
	return StorageFiles{Dir: dir}
}

func (s StorageFiles) Store(key, value []byte) error {
	fileName := filepath.Join(s.Dir, string(makeUrl(nil, key))+".txt")
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_EXCL, DEFAULT_FILE_MODE)
	if err != nil {
		if os.IsExist(err) {
			err = errDuplicate
		}
		return err
	}
	_, err = f.Write(value)
	if err != nil {
		f.Close()
		return err
	}
	err = f.Close()
	return err
}

func (s StorageFiles) Get(key []byte) (res []byte, err error) {
	fileName := filepath.Join(s.Dir, string(makeUrl(nil, key))+".txt")
	res, err = ioutil.ReadFile(fileName)
	if os.IsNotExist(err) {
		err = errNoKey
	}
	return res, err
}
