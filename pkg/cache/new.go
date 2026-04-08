package cache

import "time"

func NewGeneric[T any](ttl time.Duration) GenericCache[T] {
	return GenericCache[T]{ttl: ttl}
}

func NewKv[T any](ttl time.Duration) KvCache[T] {
	return KvCache[T]{
		cache: map[any]T{},
		ttl:   ttl,
	}
}
