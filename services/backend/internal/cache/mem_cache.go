package cache

import (
	"encoding/json"
	"fmt"
)

type MemCache struct {
	store map[string][]byte
}

func NewMemCache() Cache {
	return &MemCache{
		store: make(map[string][]byte),
	}
}

func (m *MemCache) Set(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m.store[key] = data
	return nil
}

func (m *MemCache) Get(key string, dest any) error {
	data, exists := m.store[key]
	if !exists {
		return fmt.Errorf("key %s not found in cache", key)
	}
	return json.Unmarshal(data, dest)
}

func (m *MemCache) Delete(key string) error {
	delete(m.store, key)
	return nil
}

func (m *MemCache) SetString(key, value string) error {
	m.store[key] = []byte(value)
	return nil
}

func (m *MemCache) GetString(key string) (string, error) {
	data, exists := m.store[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in cache", key)
	}
	return string(data), nil
}
