package util

type OrderedMap[K comparable, V any] struct {
	data map[K]V
	keys []K
}

type Item[K comparable, V any] struct {
	Key   K
	Value V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		data: make(map[K]V),
		keys: make([]K, 0),
	}
}

func (om *OrderedMap[K, V]) Length() int {
	return len(om.keys)
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := om.data[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.data[key] = value
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, exists := om.data[key]
	return value, exists
}

func (om *OrderedMap[K, V]) Delete(key K) bool {
	if _, exists := om.data[key]; !exists {
		return false
	}
	delete(om.data, key)

	for i, k := range om.keys {
		if k == key {
			om.keys = append(om.keys[:i], om.keys[i+1:]...)
			return true
		}
	}

	return false
}

func (om *OrderedMap[K, V]) Items() []Item[K, V] {
	items := make([]Item[K, V], 0, len(om.keys))

	for _, key := range om.keys {
		value, ok := om.data[key]
		if ok {
			items = append(items, Item[K, V]{
				Key:   key,
				Value: value,
			})
		}
	}
	return items
}
