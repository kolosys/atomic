package collection

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"sort"
	"sync"
)

// Keep is used for merge operations to indicate whether to keep a value and what value to keep.
type Keep[V any] struct {
	Keep  bool
	Value V
}

// Comparator is a function that compares two values and their keys, returning -1, 0, or 1.
type Comparator[K comparable, V any] func(firstValue, secondValue V, firstKey, secondKey K) int

// Collection is a generic map-like structure with additional utility methods.
// It is safe for concurrent use.
type Collection[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

// New creates a new Collection.
func New[K comparable, V any]() *Collection[K, V] {
	return &Collection[K, V]{items: make(map[K]V)}
}

// Set adds or updates an item in the collection.
func (c *Collection[K, V]) Set(key K, value V) *Collection[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
	return c
}

// Get retrieves an item from the collection.
func (c *Collection[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.items[key]
	return val, ok
}

// Has checks if a key exists in the collection.
func (c *Collection[K, V]) Has(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.items[key]
	return ok
}

// Delete removes an item from the collection.
func (c *Collection[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, existed := c.items[key]
	delete(c.items, key)
	return existed
}

// Clear removes all items from the collection.
func (c *Collection[K, V]) Clear() *Collection[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[K]V)
	return c
}

// Size returns the number of items in the collection.
func (c *Collection[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Keys returns all keys in the collection.
func (c *Collection[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.keysUnlocked()
}

// Values returns all values in the collection.
func (c *Collection[K, V]) Values() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]V, 0, len(c.items))
	for _, v := range c.items {
		values = append(values, v)
	}
	return values
}

// Entries returns all key-value pairs in the collection.
func (c *Collection[K, V]) Entries() [][2]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entries := make([][2]any, 0, len(c.items))
	for k, v := range c.items {
		entries = append(entries, [2]any{k, v})
	}
	return entries
}

// Clone creates a shallow copy of the collection.
func (c *Collection[K, V]) Clone() *Collection[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	clone := New[K, V]()
	for k, v := range c.items {
		clone.items[k] = v
	}
	return clone
}

// Ensure obtains the value for the given key if it exists, otherwise sets and returns the value provided by the default value generator.
func (c *Collection[K, V]) Ensure(key K, defaultValueGenerator func(key K, collection *Collection[K, V]) V) V {
	c.mu.RLock()
	if val, ok := c.items[key]; ok {
		c.mu.RUnlock()
		return val
	}
	c.mu.RUnlock()

	// Generate the default value without holding any locks
	def := defaultValueGenerator(key, c)

	// Now acquire write lock to set the value, but check again if it was set by another goroutine
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		return val // Another goroutine set it while we were generating
	}
	c.items[key] = def
	return def
}

// HasAll checks if all of the provided keys exist in the collection.
func (c *Collection[K, V]) HasAll(keys ...K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, k := range keys {
		if _, ok := c.items[k]; !ok {
			return false
		}
	}
	return true
}

// HasAny checks if any of the provided keys exist in the collection.
func (c *Collection[K, V]) HasAny(keys ...K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, k := range keys {
		if _, ok := c.items[k]; ok {
			return true
		}
	}
	return false
}

// First returns the first value(s) in the collection.
// If amount is 0, returns nil. If amount < 0, returns Last(-amount).
func (c *Collection[K, V]) First(amount ...int) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	if len(keys) == 0 {
		return nil
	}
	if len(amount) == 0 {
		return c.items[keys[0]]
	}
	n := amount[0]
	if n == 0 {
		return nil
	}
	if n < 0 {
		return c.Last(-n)
	}
	if n >= len(keys) {
		res := make([]V, 0, len(keys))
		for _, k := range keys {
			res = append(res, c.items[k])
		}
		return res
	}
	res := make([]V, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, c.items[keys[i]])
	}
	return res
}

// FirstKey returns the first key(s) in the collection.
func (c *Collection[K, V]) FirstKey(amount ...int) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	if len(keys) == 0 {
		return nil
	}
	if len(amount) == 0 {
		return keys[0]
	}
	n := amount[0]
	if n < 0 {
		return c.LastKey(-n)
	}
	if n >= len(keys) {
		return keys
	}
	return keys[:n]
}

// Last returns the last value(s) in the collection.
func (c *Collection[K, V]) Last(amount ...int) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	if len(keys) == 0 {
		return nil
	}
	if len(amount) == 0 {
		return c.items[keys[len(keys)-1]]
	}
	n := amount[0]
	if n < 0 {
		return c.First(-n)
	}
	if n == 0 {
		return []V{}
	}
	if n >= len(keys) {
		res := make([]V, 0, len(keys))
		for _, k := range keys {
			res = append(res, c.items[k])
		}
		return res
	}
	res := make([]V, 0, n)
	for i := len(keys) - n; i < len(keys); i++ {
		res = append(res, c.items[keys[i]])
	}
	return res
}

// LastKey returns the last key(s) in the collection.
func (c *Collection[K, V]) LastKey(amount ...int) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	if len(keys) == 0 {
		return nil
	}
	if len(amount) == 0 {
		return keys[len(keys)-1]
	}
	n := amount[0]
	if n < 0 {
		return c.FirstKey(-n)
	}
	if n == 0 {
		return []K{}
	}
	if n >= len(keys) {
		return keys
	}
	return keys[len(keys)-n:]
}

// At returns the value at a given index, allowing for positive and negative integers.
func (c *Collection[K, V]) At(index int) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	if index < 0 {
		index += len(keys)
	}
	if index < 0 || index >= len(keys) {
		var zero V
		return zero, false
	}
	return c.items[keys[index]], true
}

// KeyAt returns the key at a given index, allowing for positive and negative integers.
func (c *Collection[K, V]) KeyAt(index int) (K, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	if index < 0 {
		index += len(keys)
	}
	if index < 0 || index >= len(keys) {
		var zero K
		return zero, false
	}
	return keys[index], true
}

// Random returns a random value or n unique random values from the collection.
func (c *Collection[K, V]) Random(amount ...int) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	if len(keys) == 0 {
		return nil
	}
	if len(amount) == 0 {
		k := keys[rand.Intn(len(keys))]
		return c.items[k]
	}
	n := amount[0]
	if n <= 0 {
		return []V{}
	}
	if n > len(keys) {
		n = len(keys)
	}
	perm := rand.Perm(len(keys))
	res := make([]V, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, c.items[keys[perm[i]]])
	}
	return res
}

// RandomKey returns a random key or n unique random keys from the collection.
func (c *Collection[K, V]) RandomKey(amount ...int) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	if len(keys) == 0 {
		return nil
	}
	if len(amount) == 0 {
		return keys[rand.Intn(len(keys))]
	}
	n := amount[0]
	if n <= 0 {
		return []K{}
	}
	if n > len(keys) {
		n = len(keys)
	}
	perm := rand.Perm(len(keys))
	res := make([]K, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, keys[perm[i]])
	}
	return res
}

// Reverse reverses the order of the collection in place.
func (c *Collection[K, V]) Reverse() *Collection[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := c.keysUnlocked()
	for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
		keys[i], keys[j] = keys[j], keys[i]
	}
	newItems := make(map[K]V, len(c.items))
	for _, k := range keys {
		newItems[k] = c.items[k]
	}
	c.items = newItems
	return c
}

// Find returns the first value for which fn returns true.
func (c *Collection[K, V]) Find(fn func(value V, key K, collection *Collection[K, V]) bool) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.items {
		if fn(v, k, c) {
			return v, true
		}
	}
	var zero V
	return zero, false
}

// FindKey returns the first key for which fn returns true.
func (c *Collection[K, V]) FindKey(fn func(value V, key K, collection *Collection[K, V]) bool) (K, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.items {
		if fn(v, k, c) {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// FindLast returns the last value for which fn returns true.
func (c *Collection[K, V]) FindLast(fn func(value V, key K, collection *Collection[K, V]) bool) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	for i := len(keys) - 1; i >= 0; i-- {
		k := keys[i]
		v := c.items[k]
		if fn(v, k, c) {
			return v, true
		}
	}
	var zero V
	return zero, false
}

// FindLastKey returns the last key for which fn returns true.
func (c *Collection[K, V]) FindLastKey(fn func(value V, key K, collection *Collection[K, V]) bool) (K, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.keysUnlocked()
	for i := len(keys) - 1; i >= 0; i-- {
		k := keys[i]
		v := c.items[k]
		if fn(v, k, c) {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// Sweep removes items that satisfy the provided filter function. Returns the number of removed entries.
func (c *Collection[K, V]) Sweep(fn func(value V, key K, collection *Collection[K, V]) bool) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	count := 0
	for k, v := range c.items {
		if fn(v, k, c) {
			delete(c.items, k)
			count++
		}
	}
	return count
}

// Filter returns a new collection containing only the items for which fn returns true.
func (c *Collection[K, V]) Filter(fn func(value V, key K, collection *Collection[K, V]) bool) *Collection[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := New[K, V]()
	for k, v := range c.items {
		if fn(v, k, c) {
			res.items[k] = v
		}
	}
	return res
}

// Partition splits the collection into two collections: the first contains items that passed, the second those that failed.
func (c *Collection[K, V]) Partition(fn func(value V, key K, collection *Collection[K, V]) bool) (*Collection[K, V], *Collection[K, V]) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	pass := New[K, V]()
	fail := New[K, V]()
	for k, v := range c.items {
		if fn(v, k, c) {
			pass.items[k] = v
		} else {
			fail.items[k] = v
		}
	}
	return pass, fail
}

// FlatMap maps each item into a collection, then joins the results into a single collection.
func (c *Collection[K, V]) FlatMap(fn func(value V, key K, collection *Collection[K, V]) *Collection[K, V]) *Collection[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := New[K, V]()
	for k, v := range c.items {
		sub := fn(v, k, c)
		for subk, subv := range sub.items {
			result.items[subk] = subv
		}
	}
	return result
}

// Some returns true if any item passes the test.
func (c *Collection[K, V]) Some(fn func(value V, key K, collection *Collection[K, V]) bool) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.items {
		if fn(v, k, c) {
			return true
		}
	}
	return false
}

// Every returns true if all items pass the test.
func (c *Collection[K, V]) Every(fn func(value V, key K, collection *Collection[K, V]) bool) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.items {
		if !fn(v, k, c) {
			return false
		}
	}
	return true
}

// Each executes fn for each element and returns the collection.
func (c *Collection[K, V]) Each(fn func(value V, key K, collection *Collection[K, V])) *Collection[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.items {
		fn(v, k, c)
	}
	return c
}

// Tap runs a function on the collection and returns the collection.
func (c *Collection[K, V]) Tap(fn func(collection *Collection[K, V])) *Collection[K, V] {
	fn(c)
	return c
}

// Concat combines this collection with others into a new collection.
func (c *Collection[K, V]) Concat(collections ...*Collection[K, V]) *Collection[K, V] {
	result := c.Clone()
	for _, coll := range collections {
		coll.mu.RLock()
		for k, v := range coll.items {
			result.items[k] = v
		}
		coll.mu.RUnlock()
	}
	return result
}

// Equals checks if this collection shares identical items with another.
func (c *Collection[K, V]) Equals(other *Collection[K, V]) bool {
	if c == other {
		return true
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	if len(c.items) != len(other.items) {
		return false
	}
	for k, v := range c.items {
		ov, ok := other.items[k]
		if !ok || !reflect.DeepEqual(v, ov) {
			return false
		}
	}
	return true
}

// Sort sorts the items of a collection in place and returns it.
func (c *Collection[K, V]) Sort(compare Comparator[K, V]) *Collection[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := c.keysUnlocked()
	sort.SliceStable(keys, func(i, j int) bool {
		return compare(c.items[keys[i]], c.items[keys[j]], keys[i], keys[j]) < 0
	})
	newItems := make(map[K]V, len(c.items))
	for _, k := range keys {
		newItems[k] = c.items[k]
	}
	c.items = newItems
	return c
}

// Intersection returns a new collection containing the items where the key is present in both collections.
func (c *Collection[K, V]) Intersection(other *Collection[K, any]) *Collection[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	res := New[K, V]()
	for k, v := range c.items {
		if _, ok := other.items[k]; ok {
			res.items[k] = v
		}
	}
	return res
}

// Union returns a new collection containing the items where the key is present in either of the collections.
func (c *Collection[K, V]) Union(other *Collection[K, V]) *Collection[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	res := c.Clone()
	for k, v := range other.items {
		if _, ok := res.items[k]; !ok {
			res.items[k] = v
		}
	}
	return res
}

// Difference returns a new collection containing the items where the key is present in this collection but not the other.
func (c *Collection[K, V]) Difference(other *Collection[K, any]) *Collection[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	res := New[K, V]()
	for k, v := range c.items {
		if _, ok := other.items[k]; !ok {
			res.items[k] = v
		}
	}
	return res
}

// SymmetricDifference returns a new collection containing only the items where the keys are present in either collection, but not both.
func (c *Collection[K, V]) SymmetricDifference(other *Collection[K, V]) *Collection[K, V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	res := New[K, V]()
	for k, v := range c.items {
		if _, ok := other.items[k]; !ok {
			res.items[k] = v
		}
	}
	for k, v := range other.items {
		if _, ok := c.items[k]; !ok {
			res.items[k] = v
		}
	}
	return res
}

// ToReversed returns a new collection with the items in reverse order.
func (c *Collection[K, V]) ToReversed() *Collection[K, V] {
	return c.Clone().Reverse()
}

// ToSorted returns a shallow copy of the collection with the items sorted.
func (c *Collection[K, V]) ToSorted(compare Comparator[K, V]) *Collection[K, V] {
	return c.Clone().Sort(compare)
}

// ToJSON returns the collection as a JSON array of [key, value] pairs.
func (c *Collection[K, V]) ToJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	pairs := make([][2]any, 0, len(c.items))
	for k, v := range c.items {
		pairs = append(pairs, [2]any{k, v})
	}
	return json.Marshal(pairs)
}

// keysUnlocked returns the keys in insertion order. (Go maps are unordered, so this is not guaranteed.)
func (c *Collection[K, V]) keysUnlocked() []K {
	keys := make([]K, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}
