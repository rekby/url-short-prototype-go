package main

import "sync"

type StorageMap struct {
	m map[string][]byte
	mutex sync.RWMutex
}

func NewStorageMap()*StorageMap{
	return &StorageMap{
		m:make(map[string][]byte),
	}
}

func (s *StorageMap) Store (key, value []byte)error {
	keyString := string(key)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exist := s.m[keyString]; exist {
		return errDuplicate
	}
	valCopy := make([]byte, len(value))
	copy(valCopy, value)
	s.m[keyString] = valCopy
	return nil
}

func (s *StorageMap)Get(key []byte)(value []byte, err error){
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if val, exist := s.m[string(key)]; exist {
		return val, nil
	}
	return nil, errNoKey
}