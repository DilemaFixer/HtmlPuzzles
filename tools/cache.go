package tools

type Cache[TKey comparable, TValue any] struct {
	data      map[TKey]TValue
	queue     []TKey
	maxCached uint
}

func NewCache[TKey comparable, TValue any](maxCached uint) *Cache[TKey, TValue] {
	return &Cache[TKey, TValue]{
		data:      make(map[TKey]TValue),
		queue:     make([]TKey, 0),
		maxCached: maxCached,
	}
}

func (c *Cache[TKey, TValue]) Get(key TKey) (TValue, bool) {
	value, exists := c.data[key]
	if exists {
		c.moveToEnd(key)
	}
	return value, exists
}

func (c *Cache[TKey, TValue]) Set(key TKey, value TValue) {
	if _, exists := c.data[key]; exists {
		c.data[key] = value
		c.moveToEnd(key)
		return
	}

	if uint(len(c.data)) >= c.maxCached {
		c.evictOldest()
	}

	c.data[key] = value
	c.queue = append(c.queue, key)
}

func (c *Cache[TKey, TValue]) Delete(key TKey) bool {
	if _, exists := c.data[key]; !exists {
		return false
	}

	delete(c.data, key)
	c.removeFromQueue(key)
	return true
}

func (c *Cache[TKey, TValue]) Clear() {
	c.data = make(map[TKey]TValue)
	c.queue = c.queue[:0]
}

func (c *Cache[TKey, TValue]) Size() int {
	return len(c.data)
}

func (c *Cache[TKey, TValue]) Contains(key TKey) bool {
	_, exists := c.data[key]
	return exists
}

func (c *Cache[TKey, TValue]) Keys() []TKey {
	keys := make([]TKey, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

func (c *Cache[TKey, TValue]) GetOrSet(key TKey, factory func() TValue) TValue {
	if value, exists := c.Get(key); exists {
		return value
	}

	value := factory()
	c.Set(key, value)
	return value
}

func (c *Cache[TKey, TValue]) Find(predicate func(key TKey, value TValue) bool) (TKey, TValue, bool) {
	for key, value := range c.data {
		if predicate(key, value) {
			c.moveToEnd(key)
			return key, value, true
		}
	}

	var zeroKey TKey
	var zeroValue TValue
	return zeroKey, zeroValue, false
}

func (c *Cache[TKey, TValue]) FindAll(predicate func(key TKey, value TValue) bool) map[TKey]TValue {
	result := make(map[TKey]TValue)

	for key, value := range c.data {
		if predicate(key, value) {
			result[key] = value
			c.moveToEnd(key)
		}
	}

	return result
}

func (c *Cache[TKey, TValue]) FindKey(predicate func(key TKey, value TValue) bool) (TKey, bool) {
	key, _, found := c.Find(predicate)
	return key, found
}

func (c *Cache[TKey, TValue]) FindValue(predicate func(key TKey, value TValue) bool) (TValue, bool) {
	_, value, found := c.Find(predicate)
	return value, found
}

func (c *Cache[TKey, TValue]) Filter(predicate func(key TKey, value TValue) bool) *Cache[TKey, TValue] {
	filtered := NewCache[TKey, TValue](c.maxCached)

	for key, value := range c.data {
		if predicate(key, value) {
			filtered.Set(key, value)
		}
	}

	return filtered
}

func (c *Cache[TKey, TValue]) Any(predicate func(key TKey, value TValue) bool) bool {
	for key, value := range c.data {
		if predicate(key, value) {
			return true
		}
	}
	return false
}

func (c *Cache[TKey, TValue]) All(predicate func(key TKey, value TValue) bool) bool {
	for key, value := range c.data {
		if !predicate(key, value) {
			return false
		}
	}
	return true
}

func (c *Cache[TKey, TValue]) ForEach(action func(key TKey, value TValue)) {
	for key, value := range c.data {
		action(key, value)
	}
}

func (c *Cache[TKey, TValue]) Count(predicate func(key TKey, value TValue) bool) int {
	count := 0
	for key, value := range c.data {
		if predicate(key, value) {
			count++
		}
	}
	return count
}

func (c *Cache[TKey, TValue]) moveToEnd(key TKey) {
	for i, k := range c.queue {
		if k == key {
			copy(c.queue[i:], c.queue[i+1:])
			c.queue = c.queue[:len(c.queue)-1]
			break
		}
	}
	c.queue = append(c.queue, key)
}

func (c *Cache[TKey, TValue]) removeFromQueue(key TKey) {
	for i, k := range c.queue {
		if k == key {
			copy(c.queue[i:], c.queue[i+1:])
			c.queue = c.queue[:len(c.queue)-1]
			break
		}
	}
}

func (c *Cache[TKey, TValue]) evictOldest() {
	if len(c.queue) == 0 {
		return
	}

	oldestKey := c.queue[0]
	delete(c.data, oldestKey)
	c.queue = c.queue[1:]
}
