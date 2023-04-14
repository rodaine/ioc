package ioc

import (
	"sync"
)

type syncMap[K, V any] struct {
	inner sync.Map
}

func (sm *syncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := sm.inner.Load(key)
	if ok {
		value = v.(V)
	}
	return value, ok
}

func (sm *syncMap[K, V]) Store(key K, value V) {
	sm.inner.Store(key, value)
}
