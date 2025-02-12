package utils

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	m map[K]V
	sync.RWMutex
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: make(map[K]V),
	}
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.RLock()
	defer s.RUnlock()
	value, exists := s.m[key]
	return value, exists
}

func (s *SafeMap[K, V]) Delete(key K) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, key)
}

func (s *SafeMap[K, V]) Keys() []K {
	s.RLock()
	defer s.RUnlock()
	keys := make([]K, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}
