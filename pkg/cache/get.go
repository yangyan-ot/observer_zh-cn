package cache

func (a *GenericCache[T]) Get() T {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	if a.cache == nil {
		var zeroValue T
		return zeroValue
	}
	return *a.cache
}

func (c *KvCache[T]) Get(key any) (T, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	data, ok := c.cache[key]

	return data, ok
}
