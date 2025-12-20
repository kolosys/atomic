package collection

import "reflect"

// Map returns a slice of values produced by applying fn to each item.
func MapCollection[K comparable, V, R any](c *Collection[K, V], fn func(value V, key K, collection *Collection[K, V]) R) []R {
	keys := c.Keys()
	res := make([]R, 0, len(keys))
	for _, k := range keys {
		v, _ := c.Get(k)
		res = append(res, fn(v, k, c))
	}
	return res
}

// MapValues returns a new collection with the same keys but values mapped by fn.
func MapCollectionValues[K comparable, V, R any](c *Collection[K, V], fn func(value V, key K, collection *Collection[K, V]) R) *Collection[K, R] {
	keys := c.Keys()
	res := New[K, R]()
	for _, k := range keys {
		v, _ := c.Get(k)
		res.Set(k, fn(v, k, c))
	}
	return res
}

// Reduce applies a function to produce a single value.
func ReduceCollection[K comparable, V, R any](c *Collection[K, V], fn func(accumulator R, value V, key K, collection *Collection[K, V]) R, initialValue R) R {
	keys := c.Keys()
	acc := initialValue
	for _, k := range keys {
		v, _ := c.Get(k)
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
	res := New[K, R]()
	keysInSelf := c.Keys()
	keysInOther := other.Keys()

	seen := make(map[K]struct{})
	allKeys := make([]K, 0, len(keysInSelf)+len(keysInOther))

	for _, k := range keysInSelf {
		allKeys = append(allKeys, k)
		seen[k] = struct{}{}
	}
	for _, k := range keysInOther {
		if _, ok := seen[k]; !ok {
			allKeys = append(allKeys, k)
			seen[k] = struct{}{}
		}
	}

	for _, k := range allKeys {
		v, inSelf := c.Get(k)
		vo, inOther := other.Get(k)
		switch {
		case inSelf && inOther:
			keep := whenInBoth(v, vo, k)
			if keep.Keep {
				res.Set(k, keep.Value)
			}
		case inSelf:
			keep := whenInSelf(v, k)
			if keep.Keep {
				res.Set(k, keep.Value)
			}
		case inOther:
			keep := whenInOther(vo, k)
			if keep.Keep {
				res.Set(k, keep.Value)
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
		if old, ok := coll.Get(k); ok {
			coll.Set(k, combine(old, v, k))
		} else {
			coll.Set(k, v)
		}
	}
	return coll
}

// GroupBy groups items by a key selector function.
func GroupBy[K comparable, Item any](items []Item, keySelector func(item Item, index int) K) *Collection[K, []Item] {
	res := New[K, []Item]()
	for i, item := range items {
		k := keySelector(item, i)
		current, _ := res.Get(k)
		res.Set(k, append(current, item))
	}
	return res
}

// toString attempts to convert a value to string for sorting.
func toString(v any) string {
	return reflect.ValueOf(v).String()
}
