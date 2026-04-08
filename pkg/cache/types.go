package cache

import (
	"sync"
	"time"
)

type KvCache[T any] struct {
	cache     map[any]T
	ttl       time.Duration
	createdAt time.Time
	mutex     sync.RWMutex
}

type GenericCache[T any] struct {
	cache     *T
	ttl       time.Duration
	createdAt time.Time
	mutex     sync.RWMutex
}
