package collection

import "reflect"

// Map returns a slice of values produced by applying fn to each item.
func MapCollection[K comparable, V, R any](c *Collection[K, V], fn func(value V, key K, collection *Collection[K, V]) R) []R {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]R, 0, len(c.items))
	for k, v := range c.items {
		res = append(res, fn(v, k, c))
	}
	return res
}

// MapValues returns a new collection with the same keys but values mapped by fn.
func MapCollectionValues[K comparable, V, R any](c *Collection[K, V], fn func(value V, key K, collection *Collection[K, V]) R) *Collection[K, R] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := New[K, R]()
	for k, v := range c.items {
		res.items[k] = fn(v, k, c)
	}
	return res
}

// Reduce applies a function to produce a single value.
func ReduceCollection[K comparable, V, R any](c *Collection[K, V], fn func(accumulator R, value V, key K, collection *Collection[K, V]) R, initialValue R) R {
	c.mu.RLock()
	defer c.mu.RUnlock()
	acc := initialValue
	for k, v := range c.items {
		acc = fn(acc, v, k, c)
	}
	return acc
}

// ReduceRight applies a function to produce a single value, iterating from the end.
func ReduceRightCollection[K comparable, V, R any](c *Collection[K, V], fn func(accumulator R, value V, key K, collection *Collection[K, V]) R, initialValue R) R {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	acc := initialValue
	for i := len(keys) - 1; i >= 0; i-- {
		k := keys[i]
		v := c.items[k]
		acc = fn(acc, v, k, c)
	}
	return acc
}

// Merge merges two collections together into a new collection.
func MergeCollection[K comparable, V, O, R any](
	c *Collection[K, V],
	other *Collection[K, O],
	whenInSelf func(value V, key K) Keep[R],
	whenInOther func(valueOther O, key K) Keep[R],
	whenInBoth func(value V, valueOther O, key K) Keep[R],
) *Collection[K, R] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	res := New[K, R]()
	keys := make(map[K]struct{})
	for k := range c.items {
		keys[k] = struct{}{}
	}
	for k := range other.items {
		keys[k] = struct{}{}
	}
	for k := range keys {
		_, inSelf := c.items[k]
		_, inOther := other.items[k]
		switch {
		case inSelf && inOther:
			keep := whenInBoth(c.items[k], other.items[k], k)
			if keep.Keep {
				res.items[k] = keep.Value
			}
		case inSelf:
			keep := whenInSelf(c.items[k], k)
			if keep.Keep {
				res.items[k] = keep.Value
			}
		case inOther:
			keep := whenInOther(other.items[k], k)
			if keep.Keep {
				res.items[k] = keep.Value
			}
		}
	}
	return res
}

// DefaultSort is the default sort comparison algorithm used in ECMAScript.
func DefaultSort[K comparable, V any](firstValue, secondValue V, firstKey, secondKey K) int {
	x := toString(firstValue)
	y := toString(secondValue)
	if x < y {
		return -1
	}
	if y < x {
		return 1
	}
	return 0
}

// CombineEntries creates a Collection from a list of entries.
func CombineEntries[K comparable, V any](
	entries [][2]any,
	combine func(firstValue, secondValue V, key K) V,
) *Collection[K, V] {
	coll := New[K, V]()
	for _, entry := range entries {
		k := entry[0].(K)
		v := entry[1].(V)
		if old, ok := coll.items[k]; ok {
			coll.items[k] = combine(old, v, k)
		} else {
			coll.items[k] = v
		}
	}
	return coll
}

// GroupBy groups items by a key selector function.
func GroupBy[K comparable, Item any](items []Item, keySelector func(item Item, index int) K) *Collection[K, []Item] {
	res := New[K, []Item]()
	for i, item := range items {
		k := keySelector(item, i)
		res.items[k] = append(res.items[k], item)
	}
	return res
}

// toString attempts to convert a value to string for sorting.
func toString(v any) string {
	return reflect.ValueOf(v).String()
}
