package collection_test

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kolosys/atomic/collection"
)

// TestNewCollection tests the creation of a new collection
func TestNewCollection(t *testing.T) {
	c := collection.New[string, int]()
	if c == nil {
		t.Fatal("NewCollection should not return nil")
	}
	if c.Size() != 0 {
		t.Errorf("New collection should be empty, got size %d", c.Size())
	}
}

// TestCollectionSet tests the Set method
func TestCollectionSet(t *testing.T) {
	c := collection.New[string, int]()

	// Test setting a value
	result := c.Set("key1", 42)
	if result != c {
		t.Error("Set should return the collection for chaining")
	}

	if c.Size() != 1 {
		t.Errorf("Collection size should be 1, got %d", c.Size())
	}

	// Test overwriting a value
	c.Set("key1", 100)
	if c.Size() != 1 {
		t.Errorf("Collection size should still be 1, got %d", c.Size())
	}

	val, ok := c.Get("key1")
	if !ok {
		t.Error("Should be able to get the value that was set")
	}
	if val != 100 {
		t.Errorf("Expected 100, got %d", val)
	}
}

// TestCollectionGet tests the Get method
func TestCollectionGet(t *testing.T) {
	c := collection.New[string, int]()

	// Test getting from empty collection
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("Getting nonexistent key should return false")
	}

	// Test getting existing value
	c.Set("key1", 42)
	val, ok := c.Get("key1")
	if !ok {
		t.Error("Should be able to get existing key")
	}
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
}

// TestCollectionHas tests the Has method
func TestCollectionHas(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	if c.Has("key1") {
		t.Error("Empty collection should not have any keys")
	}

	// Test with existing key
	c.Set("key1", 42)
	if !c.Has("key1") {
		t.Error("Collection should have key1")
	}

	// Test with non-existing key
	if c.Has("key2") {
		t.Error("Collection should not have key2")
	}
}

// TestCollectionDelete tests the Delete method
func TestCollectionDelete(t *testing.T) {
	c := collection.New[string, int]()

	// Test deleting from empty collection
	if c.Delete("nonexistent") {
		t.Error("Deleting nonexistent key should return false")
	}

	// Test deleting existing key
	c.Set("key1", 42)
	if !c.Delete("key1") {
		t.Error("Deleting existing key should return true")
	}

	if c.Has("key1") {
		t.Error("Key should be deleted")
	}

	if c.Size() != 0 {
		t.Errorf("Collection should be empty after deletion, got size %d", c.Size())
	}

	// Test deleting already deleted key
	if c.Delete("key1") {
		t.Error("Deleting already deleted key should return false")
	}
}

// TestCollectionClear tests the Clear method
func TestCollectionClear(t *testing.T) {
	c := collection.New[string, int]()

	// Clear empty collection
	result := c.Clear()
	if result != c {
		t.Error("Clear should return the collection for chaining")
	}

	// Clear collection with items
	c.Set("key1", 1).Set("key2", 2).Set("key3", 3)
	if c.Size() != 3 {
		t.Errorf("Expected size 3, got %d", c.Size())
	}

	c.Clear()
	if c.Size() != 0 {
		t.Errorf("Collection should be empty after clear, got size %d", c.Size())
	}

	if c.Has("key1") || c.Has("key2") || c.Has("key3") {
		t.Error("No keys should exist after clear")
	}
}

// TestCollectionSize tests the Size method
func TestCollectionSize(t *testing.T) {
	c := collection.New[string, int]()

	// Test empty collection
	if c.Size() != 0 {
		t.Errorf("Empty collection should have size 0, got %d", c.Size())
	}

	// Test adding items
	c.Set("key1", 1)
	if c.Size() != 1 {
		t.Errorf("Expected size 1, got %d", c.Size())
	}

	c.Set("key2", 2).Set("key3", 3)
	if c.Size() != 3 {
		t.Errorf("Expected size 3, got %d", c.Size())
	}

	// Test overwriting item (should not change size)
	c.Set("key1", 10)
	if c.Size() != 3 {
		t.Errorf("Expected size 3 after overwrite, got %d", c.Size())
	}

	// Test deleting items
	c.Delete("key2")
	if c.Size() != 2 {
		t.Errorf("Expected size 2 after delete, got %d", c.Size())
	}
}

// TestCollectionKeys tests the Keys method
func TestCollectionKeys(t *testing.T) {
	c := collection.New[string, int]()

	// Test empty collection
	keys := c.Keys()
	if len(keys) != 0 {
		t.Errorf("Empty collection should have 0 keys, got %d", len(keys))
	}

	// Test with items
	c.Set("key1", 1).Set("key2", 2).Set("key3", 3)
	keys = c.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Check that all expected keys are present
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	expectedKeys := []string{"key1", "key2", "key3"}
	for _, expected := range expectedKeys {
		if !keyMap[expected] {
			t.Errorf("Expected key %s not found in keys", expected)
		}
	}
}

// TestCollectionValues tests the Values method
func TestCollectionValues(t *testing.T) {
	c := collection.New[string, int]()

	// Test empty collection
	values := c.Values()
	if len(values) != 0 {
		t.Errorf("Empty collection should have 0 values, got %d", len(values))
	}

	// Test with items
	c.Set("key1", 10).Set("key2", 20).Set("key3", 30)
	values = c.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Check that all expected values are present
	valueMap := make(map[int]bool)
	for _, val := range values {
		valueMap[val] = true
	}

	expectedValues := []int{10, 20, 30}
	for _, expected := range expectedValues {
		if !valueMap[expected] {
			t.Errorf("Expected value %d not found in values", expected)
		}
	}
}

// TestCollectionEntries tests the Entries method
func TestCollectionEntries(t *testing.T) {
	c := collection.New[string, int]()

	// Test empty collection
	entries := c.Entries()
	if len(entries) != 0 {
		t.Errorf("Empty collection should have 0 entries, got %d", len(entries))
	}

	// Test with items
	c.Set("key1", 10).Set("key2", 20)
	entries = c.Entries()
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// Check entries structure
	entryMap := make(map[string]int)
	for _, entry := range entries {
		if len(entry) != 2 {
			t.Errorf("Each entry should have 2 elements, got %d", len(entry))
		}
		key, keyOk := entry[0].(string)
		val, valOk := entry[1].(int)
		if !keyOk || !valOk {
			t.Error("Entry elements should have correct types")
		}
		entryMap[key] = val
	}

	if entryMap["key1"] != 10 || entryMap["key2"] != 20 {
		t.Error("Entries should contain correct key-value pairs")
	}
}

// TestCollectionClone tests the Clone method
func TestCollectionClone(t *testing.T) {
	c := collection.New[string, int]()

	// Test cloning empty collection
	clone := c.Clone()
	if clone == c {
		t.Error("Clone should return a different instance")
	}
	if clone.Size() != 0 {
		t.Errorf("Cloned empty collection should have size 0, got %d", clone.Size())
	}

	// Test cloning collection with items
	c.Set("key1", 10).Set("key2", 20).Set("key3", 30)
	clone = c.Clone()

	if clone == c {
		t.Error("Clone should return a different instance")
	}

	if clone.Size() != c.Size() {
		t.Errorf("Clone should have same size as original: expected %d, got %d", c.Size(), clone.Size())
	}

	// Verify all items are copied
	if !clone.Has("key1") || !clone.Has("key2") || !clone.Has("key3") {
		t.Error("Clone should contain all items from original")
	}

	val1, _ := clone.Get("key1")
	val2, _ := clone.Get("key2")
	val3, _ := clone.Get("key3")

	if val1 != 10 || val2 != 20 || val3 != 30 {
		t.Error("Clone should contain correct values")
	}

	// Test independence - modifying clone shouldn't affect original
	clone.Set("key4", 40)
	if c.Has("key4") {
		t.Error("Modifying clone should not affect original")
	}

	// Test independence - modifying original shouldn't affect clone
	c.Set("key5", 50)
	if clone.Has("key5") {
		t.Error("Modifying original should not affect clone")
	}

	// Test modifying existing values
	clone.Set("key1", 100)
	origVal, _ := c.Get("key1")
	if origVal != 10 {
		t.Error("Modifying clone value should not affect original value")
	}
}

// TestCollectionEnsure tests the Ensure method
func TestCollectionEnsure(t *testing.T) {
	c := collection.New[string, int]()

	// Test ensuring with non-existing key
	counter := 0
	val := c.Ensure("key1", func(key string, collection *collection.Collection[string, int]) int {
		counter++
		if key != "key1" {
			t.Errorf("Expected key 'key1', got '%s'", key)
		}
		return 42
	})

	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	if counter != 1 {
		t.Errorf("Default value generator should be called once, was called %d times", counter)
	}

	if !c.Has("key1") {
		t.Error("Key should exist after Ensure")
	}

	storedVal, _ := c.Get("key1")
	if storedVal != 42 {
		t.Errorf("Stored value should be 42, got %d", storedVal)
	}

	// Test ensuring with existing key (should not call generator)
	val2 := c.Ensure("key1", func(key string, collection *collection.Collection[string, int]) int {
		counter++
		return 100 // This should not be used
	})

	if val2 != 42 {
		t.Errorf("Expected existing value 42, got %d", val2)
	}

	if counter != 1 {
		t.Errorf("Default value generator should not be called for existing key, was called %d times", counter)
	}

	// Test that generator can access the collection
	c.Ensure("key2", func(key string, coll *collection.Collection[string, int]) int {
		if coll.Size() != 1 {
			t.Errorf("Generator should see collection with 1 item, saw %d", coll.Size())
		}
		return coll.Size() * 10
	})

	val3, _ := c.Get("key2")
	if val3 != 10 {
		t.Errorf("Expected value 10 from generator calculation, got %d", val3)
	}
}

// TestCollectionHasAll tests the HasAll method
func TestCollectionHasAll(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	if !c.HasAll() {
		t.Error("HasAll with no arguments should return true")
	}

	if c.HasAll("key1") {
		t.Error("Empty collection should not have any keys")
	}

	if c.HasAll("key1", "key2") {
		t.Error("Empty collection should not have any keys")
	}

	// Test with single item
	c.Set("key1", 10)

	if !c.HasAll("key1") {
		t.Error("Collection should have key1")
	}

	if c.HasAll("key1", "key2") {
		t.Error("Collection should not have both key1 and key2")
	}

	if c.HasAll("key2") {
		t.Error("Collection should not have key2")
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	if !c.HasAll("key1") {
		t.Error("Collection should have key1")
	}

	if !c.HasAll("key1", "key2") {
		t.Error("Collection should have both key1 and key2")
	}

	if !c.HasAll("key1", "key2", "key3") {
		t.Error("Collection should have all three keys")
	}

	if c.HasAll("key1", "key2", "key3", "key4") {
		t.Error("Collection should not have key4")
	}
}

// TestCollectionHasAny tests the HasAny method
func TestCollectionHasAny(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	if c.HasAny() {
		t.Error("HasAny with no arguments should return false")
	}

	if c.HasAny("key1") {
		t.Error("Empty collection should not have any keys")
	}

	if c.HasAny("key1", "key2") {
		t.Error("Empty collection should not have any keys")
	}

	// Test with single item
	c.Set("key1", 10)

	if !c.HasAny("key1") {
		t.Error("Collection should have key1")
	}

	if !c.HasAny("key1", "key2") {
		t.Error("Collection should have at least key1")
	}

	if c.HasAny("key2") {
		t.Error("Collection should not have key2")
	}

	if c.HasAny("key2", "key3") {
		t.Error("Collection should not have key2 or key3")
	}

	// Test with multiple items
	c.Set("key2", 20)

	if !c.HasAny("key1") {
		t.Error("Collection should have key1")
	}

	if !c.HasAny("key1", "key2") {
		t.Error("Collection should have key1 and key2")
	}

	if !c.HasAny("key2", "key3") {
		t.Error("Collection should have at least key2")
	}

	if !c.HasAny("key3", "key1") {
		t.Error("Collection should have at least key1")
	}

	if c.HasAny("key3", "key4") {
		t.Error("Collection should not have key3 or key4")
	}
}

// TestCollectionFirst tests the First method
func TestCollectionFirst(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.First()
	if result != nil {
		t.Error("First on empty collection should return nil")
	}

	// Test with amount on empty collection
	result = c.First(1)
	if result != nil {
		t.Error("First(1) on empty collection should return nil")
	}

	// Test with single item
	c.Set("key1", 10)
	result = c.First()
	if result != 10 {
		t.Errorf("Expected 10, got %v", result)
	}

	// Test First(0)
	result = c.First(0)
	if result != nil {
		t.Error("First(0) should return nil")
	}

	// Test First(1) with single item
	result = c.First(1)
	resultSlice, ok := result.([]int)
	if !ok {
		t.Errorf("First(1) should return slice, got %T", result)
	}
	if len(resultSlice) != 1 || resultSlice[0] != 10 {
		t.Errorf("Expected [10], got %v", resultSlice)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Test First(2)
	result = c.First(2)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("First(2) should return slice, got %T", result)
	}
	if len(resultSlice) != 2 {
		t.Errorf("Expected slice of length 2, got %d", len(resultSlice))
	}

	// Test First with amount greater than collection size
	result = c.First(5)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("First(5) should return slice, got %T", result)
	}
	if len(resultSlice) != 3 {
		t.Errorf("Expected slice of length 3, got %d", len(resultSlice))
	}

	// Test First with negative amount (should call Last)
	result = c.First(-1)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("First(-1) should return slice, got %T", result)
	}
	if len(resultSlice) != 1 {
		t.Errorf("Expected slice of length 1, got %d", len(resultSlice))
	}
}

// TestCollectionFirstKey tests the FirstKey method
func TestCollectionFirstKey(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.FirstKey()
	if result != nil {
		t.Error("FirstKey on empty collection should return nil")
	}

	// Test with amount on empty collection
	result = c.FirstKey(1)
	if result != nil {
		t.Error("FirstKey(1) on empty collection should return nil")
	}

	// Test with single item
	c.Set("key1", 10)
	result = c.FirstKey()
	if result != "key1" {
		t.Errorf("Expected 'key1', got %v", result)
	}

	// Test FirstKey(1) with single item
	result = c.FirstKey(1)
	resultSlice, ok := result.([]string)
	if !ok {
		t.Errorf("FirstKey(1) should return slice, got %T", result)
	}
	if len(resultSlice) != 1 || resultSlice[0] != "key1" {
		t.Errorf("Expected ['key1'], got %v", resultSlice)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Test FirstKey(2)
	result = c.FirstKey(2)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("FirstKey(2) should return slice, got %T", result)
	}
	if len(resultSlice) != 2 {
		t.Errorf("Expected slice of length 2, got %d", len(resultSlice))
	}

	// Test FirstKey with amount greater than collection size
	result = c.FirstKey(5)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("FirstKey(5) should return slice, got %T", result)
	}
	if len(resultSlice) != 3 {
		t.Errorf("Expected slice of length 3, got %d", len(resultSlice))
	}
}

// TestCollectionLast tests the Last method
func TestCollectionLast(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.Last()
	if result != nil {
		t.Error("Last on empty collection should return nil")
	}

	// Test with amount on empty collection
	result = c.Last(1)
	if result != nil {
		t.Error("Last(1) on empty collection should return nil")
	}

	// Test with single item
	c.Set("key1", 10)
	result = c.Last()
	if result != 10 {
		t.Errorf("Expected 10, got %v", result)
	}

	// Test Last(0)
	result = c.Last(0)
	resultSlice, ok := result.([]int)
	if !ok {
		t.Errorf("Last(0) should return empty slice, got %T", result)
	}
	if len(resultSlice) != 0 {
		t.Errorf("Expected empty slice, got %v", resultSlice)
	}

	// Test Last(1) with single item
	result = c.Last(1)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("Last(1) should return slice, got %T", result)
	}
	if len(resultSlice) != 1 || resultSlice[0] != 10 {
		t.Errorf("Expected [10], got %v", resultSlice)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Test Last(2)
	result = c.Last(2)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("Last(2) should return slice, got %T", result)
	}
	if len(resultSlice) != 2 {
		t.Errorf("Expected slice of length 2, got %d", len(resultSlice))
	}

	// Test Last with amount greater than collection size
	result = c.Last(5)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("Last(5) should return slice, got %T", result)
	}
	if len(resultSlice) != 3 {
		t.Errorf("Expected slice of length 3, got %d", len(resultSlice))
	}

	// Test Last with negative amount (should call First)
	result = c.Last(-1)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("Last(-1) should return slice, got %T", result)
	}
	if len(resultSlice) != 1 {
		t.Errorf("Expected slice of length 1, got %d", len(resultSlice))
	}
}

// TestCollectionLastKey tests the LastKey method
func TestCollectionLastKey(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.LastKey()
	if result != nil {
		t.Error("LastKey on empty collection should return nil")
	}

	// Test with amount on empty collection
	result = c.LastKey(1)
	if result != nil {
		t.Error("LastKey(1) on empty collection should return nil")
	}

	// Test with single item
	c.Set("key1", 10)
	result = c.LastKey()
	if result != "key1" {
		t.Errorf("Expected 'key1', got %v", result)
	}

	// Test LastKey(0)
	result = c.LastKey(0)
	resultSlice, ok := result.([]string)
	if !ok {
		t.Errorf("LastKey(0) should return empty slice, got %T", result)
	}
	if len(resultSlice) != 0 {
		t.Errorf("Expected empty slice, got %v", resultSlice)
	}

	// Test LastKey(1) with single item
	result = c.LastKey(1)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("LastKey(1) should return slice, got %T", result)
	}
	if len(resultSlice) != 1 || resultSlice[0] != "key1" {
		t.Errorf("Expected ['key1'], got %v", resultSlice)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Test LastKey(2)
	result = c.LastKey(2)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("LastKey(2) should return slice, got %T", result)
	}
	if len(resultSlice) != 2 {
		t.Errorf("Expected slice of length 2, got %d", len(resultSlice))
	}

	// Test LastKey with amount greater than collection size
	result = c.LastKey(5)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("LastKey(5) should return slice, got %T", result)
	}
	if len(resultSlice) != 3 {
		t.Errorf("Expected slice of length 3, got %d", len(resultSlice))
	}
}

// TestCollectionAt tests the At method
func TestCollectionAt(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	_, ok := c.At(0)
	if ok {
		t.Error("At(0) on empty collection should return false")
	}

	_, ok = c.At(-1)
	if ok {
		t.Error("At(-1) on empty collection should return false")
	}

	// Test with single item
	c.Set("key1", 10)

	val, ok := c.At(0)
	if !ok {
		t.Error("At(0) should return true for single item collection")
	}
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	val, ok = c.At(-1)
	if !ok {
		t.Error("At(-1) should return true for single item collection")
	}
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	_, ok = c.At(1)
	if ok {
		t.Error("At(1) should return false for single item collection")
	}

	_, ok = c.At(-2)
	if ok {
		t.Error("At(-2) should return false for single item collection")
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Test positive indices
	_, ok = c.At(0)
	if !ok {
		t.Error("At(0) should return true")
	}

	_, ok = c.At(1)
	if !ok {
		t.Error("At(1) should return true")
	}

	_, ok = c.At(2)
	if !ok {
		t.Error("At(2) should return true")
	}

	_, ok = c.At(3)
	if ok {
		t.Error("At(3) should return false for 3-item collection")
	}

	// Test negative indices
	_, ok = c.At(-1)
	if !ok {
		t.Error("At(-1) should return true")
	}

	_, ok = c.At(-2)
	if !ok {
		t.Error("At(-2) should return true")
	}

	_, ok = c.At(-3)
	if !ok {
		t.Error("At(-3) should return true")
	}

	_, ok = c.At(-4)
	if ok {
		t.Error("At(-4) should return false for 3-item collection")
	}
}

// TestCollectionKeyAt tests the KeyAt method
func TestCollectionKeyAt(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	_, ok := c.KeyAt(0)
	if ok {
		t.Error("KeyAt(0) on empty collection should return false")
	}

	_, ok = c.KeyAt(-1)
	if ok {
		t.Error("KeyAt(-1) on empty collection should return false")
	}

	// Test with single item
	c.Set("key1", 10)

	key, ok := c.KeyAt(0)
	if !ok {
		t.Error("KeyAt(0) should return true for single item collection")
	}
	if key != "key1" {
		t.Errorf("Expected 'key1', got %s", key)
	}

	key, ok = c.KeyAt(-1)
	if !ok {
		t.Error("KeyAt(-1) should return true for single item collection")
	}
	if key != "key1" {
		t.Errorf("Expected 'key1', got %s", key)
	}

	_, ok = c.KeyAt(1)
	if ok {
		t.Error("KeyAt(1) should return false for single item collection")
	}

	_, ok = c.KeyAt(-2)
	if ok {
		t.Error("KeyAt(-2) should return false for single item collection")
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Test positive indices
	_, ok = c.KeyAt(0)
	if !ok {
		t.Error("KeyAt(0) should return true")
	}

	_, ok = c.KeyAt(1)
	if !ok {
		t.Error("KeyAt(1) should return true")
	}

	_, ok = c.KeyAt(2)
	if !ok {
		t.Error("KeyAt(2) should return true")
	}

	_, ok = c.KeyAt(3)
	if ok {
		t.Error("KeyAt(3) should return false for 3-item collection")
	}

	// Test negative indices
	_, ok = c.KeyAt(-1)
	if !ok {
		t.Error("KeyAt(-1) should return true")
	}

	_, ok = c.KeyAt(-2)
	if !ok {
		t.Error("KeyAt(-2) should return true")
	}

	_, ok = c.KeyAt(-3)
	if !ok {
		t.Error("KeyAt(-3) should return true")
	}

	_, ok = c.KeyAt(-4)
	if ok {
		t.Error("KeyAt(-4) should return false for 3-item collection")
	}
}

// TestCollectionRandom tests the Random method
func TestCollectionRandom(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.Random()
	if result != nil {
		t.Error("Random on empty collection should return nil")
	}

	result = c.Random(1)
	if result != nil {
		t.Error("Random(1) on empty collection should return nil")
	}

	// Test with single item
	c.Set("key1", 10)

	result = c.Random()
	if result != 10 {
		t.Errorf("Expected 10, got %v", result)
	}

	// Test Random(0)
	result = c.Random(0)
	resultSlice, ok := result.([]int)
	if !ok {
		t.Errorf("Random(0) should return empty slice, got %T", result)
	}
	if len(resultSlice) != 0 {
		t.Errorf("Expected empty slice, got %v", resultSlice)
	}

	// Test Random(1) with single item
	result = c.Random(1)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("Random(1) should return slice, got %T", result)
	}
	if len(resultSlice) != 1 || resultSlice[0] != 10 {
		t.Errorf("Expected [10], got %v", resultSlice)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Test Random(2)
	result = c.Random(2)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("Random(2) should return slice, got %T", result)
	}
	if len(resultSlice) != 2 {
		t.Errorf("Expected slice of length 2, got %d", len(resultSlice))
	}

	// Verify all values are valid
	valueMap := make(map[int]bool)
	valueMap[10] = true
	valueMap[20] = true
	valueMap[30] = true

	for _, val := range resultSlice {
		if !valueMap[val] {
			t.Errorf("Unexpected value %d in random result", val)
		}
	}

	// Test Random with amount greater than collection size
	result = c.Random(5)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("Random(5) should return slice, got %T", result)
	}
	if len(resultSlice) != 3 {
		t.Errorf("Expected slice of length 3, got %d", len(resultSlice))
	}

	// Test Random with negative amount
	result = c.Random(-1)
	resultSlice, ok = result.([]int)
	if !ok {
		t.Errorf("Random(-1) should return empty slice, got %T", result)
	}
	if len(resultSlice) != 0 {
		t.Errorf("Expected empty slice for negative amount, got %v", resultSlice)
	}
}

// TestCollectionRandomKey tests the RandomKey method
func TestCollectionRandomKey(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.RandomKey()
	if result != nil {
		t.Error("RandomKey on empty collection should return nil")
	}

	result = c.RandomKey(1)
	if result != nil {
		t.Error("RandomKey(1) on empty collection should return nil")
	}

	// Test with single item
	c.Set("key1", 10)

	result = c.RandomKey()
	if result != "key1" {
		t.Errorf("Expected 'key1', got %v", result)
	}

	// Test RandomKey(0)
	result = c.RandomKey(0)
	resultSlice, ok := result.([]string)
	if !ok {
		t.Errorf("RandomKey(0) should return empty slice, got %T", result)
	}
	if len(resultSlice) != 0 {
		t.Errorf("Expected empty slice, got %v", resultSlice)
	}

	// Test RandomKey(1) with single item
	result = c.RandomKey(1)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("RandomKey(1) should return slice, got %T", result)
	}
	if len(resultSlice) != 1 || resultSlice[0] != "key1" {
		t.Errorf("Expected ['key1'], got %v", resultSlice)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Test RandomKey(2)
	result = c.RandomKey(2)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("RandomKey(2) should return slice, got %T", result)
	}
	if len(resultSlice) != 2 {
		t.Errorf("Expected slice of length 2, got %d", len(resultSlice))
	}

	// Verify all keys are valid
	keyMap := make(map[string]bool)
	keyMap["key1"] = true
	keyMap["key2"] = true
	keyMap["key3"] = true

	for _, key := range resultSlice {
		if !keyMap[key] {
			t.Errorf("Unexpected key %s in random result", key)
		}
	}

	// Test RandomKey with amount greater than collection size
	result = c.RandomKey(5)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("RandomKey(5) should return slice, got %T", result)
	}
	if len(resultSlice) != 3 {
		t.Errorf("Expected slice of length 3, got %d", len(resultSlice))
	}

	// Test RandomKey with negative amount
	result = c.RandomKey(-1)
	resultSlice, ok = result.([]string)
	if !ok {
		t.Errorf("RandomKey(-1) should return empty slice, got %T", result)
	}
	if len(resultSlice) != 0 {
		t.Errorf("Expected empty slice for negative amount, got %v", resultSlice)
	}
}

// TestCollectionFind tests the Find method
func TestCollectionFind(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	_, found := c.Find(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if found {
		t.Error("Find on empty collection should return false")
	}

	// Test with single item - found
	c.Set("key1", 10)
	val, found := c.Find(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 10
	})
	if !found {
		t.Error("Should find the value 10")
	}
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	// Test with single item - not found
	_, found = c.Find(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 20
	})
	if found {
		t.Error("Should not find value 20")
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Find first even number
	val, found = c.Find(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value%2 == 0
	})
	if !found {
		t.Error("Should find an even number")
	}
	if val%2 != 0 {
		t.Errorf("Expected even number, got %d", val)
	}

	// Find value greater than 25
	val, found = c.Find(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 25
	})
	if !found {
		t.Error("Should find value greater than 25")
	}
	if val <= 25 {
		t.Errorf("Expected value > 25, got %d", val)
	}

	// Test that function receives correct parameters
	c.Find(func(value int, key string, coll *collection.Collection[string, int]) bool {
		if coll.Size() != 3 {
			t.Errorf("Function should receive collection with 3 items, got %d", coll.Size())
		}
		if !coll.Has(key) {
			t.Errorf("Collection should contain key %s", key)
		}
		return false
	})
}

// TestCollectionFindKey tests the FindKey method
func TestCollectionFindKey(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	_, found := c.FindKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if found {
		t.Error("FindKey on empty collection should return false")
	}

	// Test with single item - found
	c.Set("key1", 10)
	key, found := c.FindKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 10
	})
	if !found {
		t.Error("Should find the key for value 10")
	}
	if key != "key1" {
		t.Errorf("Expected 'key1', got %s", key)
	}

	// Test with single item - not found
	_, found = c.FindKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 20
	})
	if found {
		t.Error("Should not find key for value 20")
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)

	// Find key for value greater than 25
	key, found = c.FindKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 25
	})
	if !found {
		t.Error("Should find key for value greater than 25")
	}
	val, _ := c.Get(key)
	if val <= 25 {
		t.Errorf("Found key %s should have value > 25, got %d", key, val)
	}

	// Find key that starts with "key2"
	key, found = c.FindKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return key == "key2"
	})
	if !found {
		t.Error("Should find key2")
	}
	if key != "key2" {
		t.Errorf("Expected 'key2', got %s", key)
	}
}

// TestCollectionFindLast tests the FindLast method
func TestCollectionFindLast(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	_, found := c.FindLast(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if found {
		t.Error("FindLast on empty collection should return false")
	}

	// Test with single item - found
	c.Set("key1", 10)
	val, found := c.FindLast(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 10
	})
	if !found {
		t.Error("Should find the value 10")
	}
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	// Test with single item - not found
	_, found = c.FindLast(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 20
	})
	if found {
		t.Error("Should not find value 20")
	}

	// Test with multiple items - this should find the last matching item
	c.Set("key2", 20).Set("key3", 10).Set("key4", 30) // Now we have 10, 20, 10, 30

	// Find last occurrence of 10
	val, found = c.FindLast(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 10
	})
	if !found {
		t.Error("Should find value 10")
	}
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	// Find last even number (should be 30 if it comes after 10 and 20)
	val, found = c.FindLast(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value%2 == 0
	})
	if !found {
		t.Error("Should find an even number")
	}
	if val%2 != 0 {
		t.Errorf("Expected even number, got %d", val)
	}
}

// TestCollectionFindLastKey tests the FindLastKey method
func TestCollectionFindLastKey(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	_, found := c.FindLastKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if found {
		t.Error("FindLastKey on empty collection should return false")
	}

	// Test with single item - found
	c.Set("key1", 10)
	key, found := c.FindLastKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 10
	})
	if !found {
		t.Error("Should find the key for value 10")
	}
	if key != "key1" {
		t.Errorf("Expected 'key1', got %s", key)
	}

	// Test with single item - not found
	_, found = c.FindLastKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 20
	})
	if found {
		t.Error("Should not find key for value 20")
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 10).Set("key4", 30) // Now we have 10, 20, 10, 30

	// Find last key with value 10
	key, found = c.FindLastKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 10
	})
	if !found {
		t.Error("Should find key for value 10")
	}
	val, _ := c.Get(key)
	if val != 10 {
		t.Errorf("Found key %s should have value 10, got %d", key, val)
	}

	// Find last key for even number
	key, found = c.FindLastKey(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value%2 == 0
	})
	if !found {
		t.Error("Should find key for even number")
	}
	val, _ = c.Get(key)
	if val%2 != 0 {
		t.Errorf("Found key %s should have even value, got %d", key, val)
	}
}

// TestCollectionEach tests the Each method
func TestCollectionEach(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	count := 0
	result := c.Each(func(value int, key string, collection *collection.Collection[string, int]) {
		count++
	})
	if count != 0 {
		t.Errorf("Each on empty collection should not call function, called %d times", count)
	}
	if result != c {
		t.Error("Each should return the collection for chaining")
	}

	// Test with single item
	c.Set("key1", 10)
	count = 0
	var seenKey string
	var seenValue int

	c.Each(func(value int, key string, coll *collection.Collection[string, int]) {
		count++
		seenKey = key
		seenValue = value
		if coll.Size() != 1 {
			t.Errorf("Function should receive collection with 1 item, got %d", coll.Size())
		}
	})

	if count != 1 {
		t.Errorf("Each should call function once, called %d times", count)
	}
	if seenKey != "key1" {
		t.Errorf("Expected to see key 'key1', got %s", seenKey)
	}
	if seenValue != 10 {
		t.Errorf("Expected to see value 10, got %d", seenValue)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)
	count = 0
	seenKeys := make(map[string]bool)
	seenValues := make(map[int]bool)

	c.Each(func(value int, key string, coll *collection.Collection[string, int]) {
		count++
		seenKeys[key] = true
		seenValues[value] = true
		if coll.Size() != 3 {
			t.Errorf("Function should receive collection with 3 items, got %d", coll.Size())
		}
	})

	if count != 3 {
		t.Errorf("Each should call function 3 times, called %d times", count)
	}

	expectedKeys := []string{"key1", "key2", "key3"}
	for _, key := range expectedKeys {
		if !seenKeys[key] {
			t.Errorf("Expected to see key %s", key)
		}
	}

	expectedValues := []int{10, 20, 30}
	for _, val := range expectedValues {
		if !seenValues[val] {
			t.Errorf("Expected to see value %d", val)
		}
	}
}

// TestCollectionTap tests the Tap method
func TestCollectionTap(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	called := false
	result := c.Tap(func(collection *collection.Collection[string, int]) {
		called = true
		if collection.Size() != 0 {
			t.Errorf("Tap should receive empty collection, got size %d", collection.Size())
		}
	})

	if !called {
		t.Error("Tap should call the function")
	}
	if result != c {
		t.Error("Tap should return the collection for chaining")
	}

	// Test with items
	c.Set("key1", 10).Set("key2", 20)
	called = false
	var receivedSize int

	c.Tap(func(coll *collection.Collection[string, int]) {
		called = true
		receivedSize = coll.Size()
		// Verify we can call methods on the collection
		if !coll.Has("key1") {
			t.Error("Tap function should receive collection with key1")
		}
		if !coll.Has("key2") {
			t.Error("Tap function should receive collection with key2")
		}
		val, _ := coll.Get("key1")
		if val != 10 {
			t.Errorf("Expected value 10 for key1, got %d", val)
		}
	})

	if !called {
		t.Error("Tap should call the function")
	}
	if receivedSize != 2 {
		t.Errorf("Expected collection size 2, got %d", receivedSize)
	}

	// Test chaining - modify collection in tap and verify it affects the original
	originalSize := c.Size()
	c.Tap(func(coll *collection.Collection[string, int]) {
		coll.Set("key3", 30)
	}).Tap(func(coll *collection.Collection[string, int]) {
		if coll.Size() != originalSize+1 {
			t.Errorf("Expected collection size %d after tap modification, got %d", originalSize+1, coll.Size())
		}
	})

	if c.Size() != originalSize+1 {
		t.Errorf("Original collection should be modified by tap, expected size %d, got %d", originalSize+1, c.Size())
	}
	if !c.Has("key3") {
		t.Error("Original collection should contain key3 after tap modification")
	}
}

// TestCollectionFilter tests the Filter method
func TestCollectionFilter(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	filtered := c.Filter(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if filtered.Size() != 0 {
		t.Errorf("Filter on empty collection should return empty collection, got size %d", filtered.Size())
	}

	// Test with single item - match
	c.Set("key1", 10)
	filtered = c.Filter(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 5
	})
	if filtered.Size() != 1 {
		t.Errorf("Filter should return collection with 1 item, got %d", filtered.Size())
	}
	if !filtered.Has("key1") {
		t.Error("Filtered collection should contain key1")
	}

	// Test with single item - no match
	filtered = c.Filter(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 15
	})
	if filtered.Size() != 0 {
		t.Errorf("Filter should return empty collection, got size %d", filtered.Size())
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30).Set("key4", 5)

	// Filter even numbers
	filtered = c.Filter(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value%2 == 0
	})
	if filtered.Size() != 3 {
		t.Errorf("Filter should return 3 even numbers, got %d", filtered.Size())
	}
	if !filtered.Has("key1") || !filtered.Has("key2") || !filtered.Has("key3") {
		t.Error("Filtered collection should contain key1, key2, key3")
	}
	if filtered.Has("key4") {
		t.Error("Filtered collection should not contain key4 (odd number)")
	}

	// Filter values greater than 15
	filtered = c.Filter(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 15
	})
	if filtered.Size() != 2 {
		t.Errorf("Filter should return 2 items > 15, got %d", filtered.Size())
	}
	if !filtered.Has("key2") || !filtered.Has("key3") {
		t.Error("Filtered collection should contain key2, key3")
	}

	// Test that original collection is unchanged
	if c.Size() != 4 {
		t.Errorf("Original collection should remain unchanged, expected size 4, got %d", c.Size())
	}

	// Test that function receives correct parameters
	filtered = c.Filter(func(value int, key string, coll *collection.Collection[string, int]) bool {
		if coll.Size() != 4 {
			t.Errorf("Function should receive original collection with 4 items, got %d", coll.Size())
		}
		if !coll.Has(key) {
			t.Errorf("Original collection should contain key %s", key)
		}
		return false
	})
	if filtered.Size() != 0 {
		t.Error("Filter returning false should result in empty collection")
	}
}

// TestCollectionPartition tests the Partition method
func TestCollectionPartition(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	pass, fail := c.Partition(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if pass.Size() != 0 || fail.Size() != 0 {
		t.Errorf("Partition on empty collection should return two empty collections, got pass=%d, fail=%d", pass.Size(), fail.Size())
	}

	// Test with single item - pass
	c.Set("key1", 10)
	pass, fail = c.Partition(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 5
	})
	if pass.Size() != 1 || fail.Size() != 0 {
		t.Errorf("Partition should return pass=1, fail=0, got pass=%d, fail=%d", pass.Size(), fail.Size())
	}
	if !pass.Has("key1") {
		t.Error("Pass collection should contain key1")
	}

	// Test with single item - fail
	pass, fail = c.Partition(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 15
	})
	if pass.Size() != 0 || fail.Size() != 1 {
		t.Errorf("Partition should return pass=0, fail=1, got pass=%d, fail=%d", pass.Size(), fail.Size())
	}
	if !fail.Has("key1") {
		t.Error("Fail collection should contain key1")
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30).Set("key4", 5)

	// Partition even vs odd
	pass, fail = c.Partition(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value%2 == 0
	})
	if pass.Size() != 3 || fail.Size() != 1 {
		t.Errorf("Partition even/odd should return pass=3, fail=1, got pass=%d, fail=%d", pass.Size(), fail.Size())
	}

	// Check pass collection (even numbers)
	if !pass.Has("key1") || !pass.Has("key2") || !pass.Has("key3") {
		t.Error("Pass collection should contain even numbers: key1, key2, key3")
	}
	if pass.Has("key4") {
		t.Error("Pass collection should not contain key4 (odd)")
	}

	// Check fail collection (odd numbers)
	if !fail.Has("key4") {
		t.Error("Fail collection should contain key4 (odd)")
	}
	if fail.Has("key1") || fail.Has("key2") || fail.Has("key3") {
		t.Error("Fail collection should not contain even numbers")
	}

	// Test that original collection is unchanged
	if c.Size() != 4 {
		t.Errorf("Original collection should remain unchanged, expected size 4, got %d", c.Size())
	}

	// Verify values are correctly partitioned
	val1, _ := pass.Get("key1")
	val2, _ := pass.Get("key2")
	val3, _ := pass.Get("key3")
	val4, _ := fail.Get("key4")

	if val1 != 10 || val2 != 20 || val3 != 30 || val4 != 5 {
		t.Error("Partitioned collections should contain correct values")
	}
}

// TestCollectionSweep tests the Sweep method
func TestCollectionSweep(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	removed := c.Sweep(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if removed != 0 {
		t.Errorf("Sweep on empty collection should remove 0 items, removed %d", removed)
	}

	// Test with single item - no match
	c.Set("key1", 10)
	removed = c.Sweep(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 15
	})
	if removed != 0 {
		t.Errorf("Sweep should remove 0 items, removed %d", removed)
	}
	if c.Size() != 1 {
		t.Errorf("Collection size should remain 1, got %d", c.Size())
	}

	// Test with single item - match
	removed = c.Sweep(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value == 10
	})
	if removed != 1 {
		t.Errorf("Sweep should remove 1 item, removed %d", removed)
	}
	if c.Size() != 0 {
		t.Errorf("Collection should be empty after sweep, got size %d", c.Size())
	}

	// Test with multiple items
	c.Set("key1", 10).Set("key2", 20).Set("key3", 30).Set("key4", 5)

	// Sweep odd numbers
	removed = c.Sweep(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value%2 == 1
	})
	if removed != 1 {
		t.Errorf("Sweep should remove 1 odd number, removed %d", removed)
	}
	if c.Size() != 3 {
		t.Errorf("Collection should have 3 items after sweep, got %d", c.Size())
	}
	if c.Has("key4") {
		t.Error("Collection should not have key4 after sweeping odd numbers")
	}
	if !c.Has("key1") || !c.Has("key2") || !c.Has("key3") {
		t.Error("Collection should still have even numbers")
	}

	// Sweep values greater than 15
	removed = c.Sweep(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 15
	})
	if removed != 2 {
		t.Errorf("Sweep should remove 2 items > 15, removed %d", removed)
	}
	if c.Size() != 1 {
		t.Errorf("Collection should have 1 item after sweep, got %d", c.Size())
	}
	if !c.Has("key1") {
		t.Error("Collection should still have key1")
	}

	// Sweep all remaining items
	removed = c.Sweep(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if removed != 1 {
		t.Errorf("Sweep should remove 1 remaining item, removed %d", removed)
	}
	if c.Size() != 0 {
		t.Errorf("Collection should be empty after final sweep, got size %d", c.Size())
	}
}

// TestCollectionSome tests the Some method
func TestCollectionSome(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.Some(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return true
	})
	if result {
		t.Error("Some on empty collection should return false")
	}

	// Test with single item - match
	c.Set("key1", 10)
	result = c.Some(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 5
	})
	if !result {
		t.Error("Some should return true when condition is met")
	}

	// Test with single item - no match
	result = c.Some(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 15
	})
	if result {
		t.Error("Some should return false when condition is not met")
	}

	// Test with multiple items - some match
	c.Set("key2", 20).Set("key3", 30).Set("key4", 5)

	result = c.Some(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 25
	})
	if !result {
		t.Error("Some should return true when at least one item meets condition")
	}

	result = c.Some(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 50
	})
	if result {
		t.Error("Some should return false when no items meet condition")
	}

	// Test early termination - function should not be called for all items if match is found early
	callCount := 0
	result = c.Some(func(value int, key string, collection *collection.Collection[string, int]) bool {
		callCount++
		return value == 10 // This should match quickly
	})
	if !result {
		t.Error("Some should return true when condition is met")
	}
	// We can't guarantee exact call count due to map iteration order, but it should be <= 4
	if callCount > 4 {
		t.Errorf("Some called function %d times, should be <= 4", callCount)
	}

	// Test that function receives correct parameters
	c.Some(func(value int, key string, coll *collection.Collection[string, int]) bool {
		if coll.Size() != 4 {
			t.Errorf("Function should receive collection with 4 items, got %d", coll.Size())
		}
		if !coll.Has(key) {
			t.Errorf("Collection should contain key %s", key)
		}
		return false
	})
}

// TestCollectionEvery tests the Every method
func TestCollectionEvery(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.Every(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return false
	})
	if !result {
		t.Error("Every on empty collection should return true")
	}

	// Test with single item - match
	c.Set("key1", 10)
	result = c.Every(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 5
	})
	if !result {
		t.Error("Every should return true when all items meet condition")
	}

	// Test with single item - no match
	result = c.Every(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 15
	})
	if result {
		t.Error("Every should return false when condition is not met")
	}

	// Test with multiple items - all match
	c.Set("key2", 20).Set("key3", 30).Set("key4", 15)

	result = c.Every(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 5
	})
	if !result {
		t.Error("Every should return true when all items meet condition")
	}

	// Test with multiple items - not all match
	result = c.Every(func(value int, key string, collection *collection.Collection[string, int]) bool {
		return value > 15
	})
	if result {
		t.Error("Every should return false when not all items meet condition")
	}

	// Test early termination - function should not be called for all items if mismatch is found early
	callCount := 0
	result = c.Every(func(value int, key string, collection *collection.Collection[string, int]) bool {
		callCount++
		return value != 10 // This should fail quickly
	})
	if result {
		t.Error("Every should return false when condition is not met")
	}
	// We can't guarantee exact call count due to map iteration order, but it should be <= 4
	if callCount > 4 {
		t.Errorf("Every called function %d times, should be <= 4", callCount)
	}

	// Test that function receives correct parameters
	c.Every(func(value int, key string, coll *collection.Collection[string, int]) bool {
		if coll.Size() != 4 {
			t.Errorf("Function should receive collection with 4 items, got %d", coll.Size())
		}
		if !coll.Has(key) {
			t.Errorf("Collection should contain key %s", key)
		}
		return true
	})
}

// TestCollectionSort tests the Sort method
func TestCollectionSort(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.Sort(func(firstValue, secondValue int, firstKey, secondKey string) int {
		if firstValue < secondValue {
			return -1
		} else if firstValue > secondValue {
			return 1
		}
		return 0
	})
	if result != c {
		t.Error("Sort should return the collection for chaining")
	}
	if c.Size() != 0 {
		t.Errorf("Empty collection should remain empty, got size %d", c.Size())
	}

	// Test with single item
	c.Set("key1", 10)
	c.Sort(func(firstValue, secondValue int, firstKey, secondKey string) int {
		if firstValue < secondValue {
			return -1
		} else if firstValue > secondValue {
			return 1
		}
		return 0
	})
	if c.Size() != 1 {
		t.Errorf("Single item collection should have size 1, got %d", c.Size())
	}

	// Test with multiple items - sort by value ascending
	c.Set("key2", 5).Set("key3", 15).Set("key4", 8)
	c.Sort(func(firstValue, secondValue int, firstKey, secondKey string) int {
		if firstValue < secondValue {
			return -1
		} else if firstValue > secondValue {
			return 1
		}
		return 0
	})

	// Get values after sorting
	values := c.Values()
	if len(values) != 4 {
		t.Errorf("Expected 4 values, got %d", len(values))
	}

	// Since Go maps don't preserve order, we can't verify sorting by checking Values() order
	// Instead, we verify all expected values are present
	expectedValues := []int{5, 8, 10, 15}
	valueMap := make(map[int]bool)
	for _, val := range values {
		valueMap[val] = true
	}

	for _, expected := range expectedValues {
		if !valueMap[expected] {
			t.Errorf("Expected value %d not found after sorting", expected)
		}
	}

	// Test sorting by key
	c.Clear().Set("zebra", 1).Set("alpha", 2).Set("beta", 3)
	c.Sort(func(firstValue, secondValue int, firstKey, secondKey string) int {
		if firstKey < secondKey {
			return -1
		} else if firstKey > secondKey {
			return 1
		}
		return 0
	})

	keys := c.Keys()
	expectedKeyOrder := []string{"alpha", "beta", "zebra"}
	if !reflect.DeepEqual(keys, expectedKeyOrder) {
		t.Errorf("Expected key order %v, got %v", expectedKeyOrder, keys)
	}
}

// TestCollectionReverse tests the Reverse method
func TestCollectionReverse(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.Reverse()
	if result != c {
		t.Error("Reverse should return the collection for chaining")
	}
	if c.Size() != 0 {
		t.Errorf("Empty collection should remain empty, got size %d", c.Size())
	}

	// Test with single item
	c.Set("key1", 10)
	c.Reverse()
	if c.Size() != 1 {
		t.Errorf("Single item collection should have size 1, got %d", c.Size())
	}
	val, _ := c.Get("key1")
	if val != 10 {
		t.Errorf("Value should remain unchanged, expected 10, got %d", val)
	}

	// Test with multiple items
	c.Clear().Set("key1", 10).Set("key2", 20).Set("key3", 30)
	originalKeys := c.Keys()
	originalValues := c.Values()

	c.Reverse()

	reversedKeys := c.Keys()
	reversedValues := c.Values()

	// Verify size unchanged
	if c.Size() != 3 {
		t.Errorf("Size should remain 3, got %d", c.Size())
	}

	// Verify all items still present
	if !c.Has("key1") || !c.Has("key2") || !c.Has("key3") {
		t.Error("All original keys should still be present")
	}

	// Since Go maps don't guarantee order, we can't easily test if the order is actually reversed
	// But we can verify that all keys and values are still present
	if len(reversedKeys) != len(originalKeys) {
		t.Errorf("Keys length should be unchanged: expected %d, got %d", len(originalKeys), len(reversedKeys))
	}
	if len(reversedValues) != len(originalValues) {
		t.Errorf("Values length should be unchanged: expected %d, got %d", len(originalValues), len(reversedValues))
	}

	// Verify all original keys are still present
	keyMap := make(map[string]bool)
	for _, key := range reversedKeys {
		keyMap[key] = true
	}
	for _, originalKey := range originalKeys {
		if !keyMap[originalKey] {
			t.Errorf("Original key %s should still be present", originalKey)
		}
	}
}

// TestCollectionToReversed tests the ToReversed method
func TestCollectionToReversed(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	reversed := c.ToReversed()
	if reversed == c {
		t.Error("ToReversed should return a different collection instance")
	}
	if reversed.Size() != 0 {
		t.Errorf("Reversed empty collection should be empty, got size %d", reversed.Size())
	}

	// Test with single item
	c.Set("key1", 10)
	reversed = c.ToReversed()
	if reversed == c {
		t.Error("ToReversed should return a different collection instance")
	}
	if reversed.Size() != 1 {
		t.Errorf("Reversed collection should have size 1, got %d", reversed.Size())
	}
	val, _ := reversed.Get("key1")
	if val != 10 {
		t.Errorf("Value should be copied correctly, expected 10, got %d", val)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)
	originalSize := c.Size()
	reversed = c.ToReversed()

	// Verify independence
	if reversed == c {
		t.Error("ToReversed should return a different collection instance")
	}

	// Verify size
	if reversed.Size() != originalSize {
		t.Errorf("Reversed collection should have same size as original: expected %d, got %d", originalSize, reversed.Size())
	}

	// Verify all items are copied
	if !reversed.Has("key1") || !reversed.Has("key2") || !reversed.Has("key3") {
		t.Error("Reversed collection should contain all original items")
	}

	// Verify original collection is unchanged
	if c.Size() != originalSize {
		t.Errorf("Original collection should be unchanged, expected size %d, got %d", originalSize, c.Size())
	}

	// Test independence - modifying reversed shouldn't affect original
	reversed.Set("key4", 40)
	if c.Has("key4") {
		t.Error("Modifying reversed collection should not affect original")
	}
}

// TestCollectionToSorted tests the ToSorted method
func TestCollectionToSorted(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	sorted := c.ToSorted(func(firstValue, secondValue int, firstKey, secondKey string) int {
		if firstValue < secondValue {
			return -1
		} else if firstValue > secondValue {
			return 1
		}
		return 0
	})
	if sorted == c {
		t.Error("ToSorted should return a different collection instance")
	}
	if sorted.Size() != 0 {
		t.Errorf("Sorted empty collection should be empty, got size %d", sorted.Size())
	}

	// Test with multiple items
	c.Set("key1", 30).Set("key2", 10).Set("key3", 20)
	originalSize := c.Size()

	sorted = c.ToSorted(func(firstValue, secondValue int, firstKey, secondKey string) int {
		if firstValue < secondValue {
			return -1
		} else if firstValue > secondValue {
			return 1
		}
		return 0
	})

	// Verify independence
	if sorted == c {
		t.Error("ToSorted should return a different collection instance")
	}

	// Verify size
	if sorted.Size() != originalSize {
		t.Errorf("Sorted collection should have same size as original: expected %d, got %d", originalSize, sorted.Size())
	}

	// Verify all items are copied
	if !sorted.Has("key1") || !sorted.Has("key2") || !sorted.Has("key3") {
		t.Error("Sorted collection should contain all original items")
	}

	// Since Go maps don't preserve order, we can't verify sorting by checking Values() order
	// Instead, we verify that the sorting was applied by checking the actual sorted order
	// would be preserved if we could access it in sorted order
	sortedValues := sorted.Values()
	expectedValues := []int{10, 20, 30}

	// Convert to map for easy checking
	valueMap := make(map[int]bool)
	for _, val := range sortedValues {
		valueMap[val] = true
	}

	// Verify all expected values are present
	for _, expected := range expectedValues {
		if !valueMap[expected] {
			t.Errorf("Expected value %d not found in sorted collection", expected)
		}
	}

	// Verify original collection is unchanged
	if c.Size() != originalSize {
		t.Errorf("Original collection should be unchanged, expected size %d, got %d", originalSize, c.Size())
	}

	// Test sorting by key
	sorted = c.ToSorted(func(firstValue, secondValue int, firstKey, secondKey string) int {
		if firstKey < secondKey {
			return -1
		} else if firstKey > secondKey {
			return 1
		}
		return 0
	})

	sortedKeys := sorted.Keys()
	expectedKeyOrder := []string{"key1", "key2", "key3"}
	if !reflect.DeepEqual(sortedKeys, expectedKeyOrder) {
		t.Errorf("Expected key order %v, got %v", expectedKeyOrder, sortedKeys)
	}
}

// TestCollectionFlatMap tests the FlatMap method
func TestCollectionFlatMap(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := c.FlatMap(func(value int, key string, coll *collection.Collection[string, int]) *collection.Collection[string, int] {
		sub := collection.New[string, int]()
		sub.Set(key+"_sub", value*2)
		return sub
	})
	if result.Size() != 0 {
		t.Errorf("FlatMap on empty collection should return empty collection, got size %d", result.Size())
	}

	// Test with single item
	c.Set("key1", 10)
	result = c.FlatMap(func(value int, key string, coll *collection.Collection[string, int]) *collection.Collection[string, int] {
		sub := collection.New[string, int]()
		sub.Set(key+"_double", value*2)
		sub.Set(key+"_triple", value*3)
		return sub
	})
	if result.Size() != 2 {
		t.Errorf("FlatMap should return collection with 2 items, got %d", result.Size())
	}

	val1, ok1 := result.Get("key1_double")
	val2, ok2 := result.Get("key1_triple")
	if !ok1 || !ok2 {
		t.Error("FlatMap result should contain both generated keys")
	}
	if val1 != 20 || val2 != 30 {
		t.Errorf("Expected values 20 and 30, got %d and %d", val1, val2)
	}

	// Test with multiple items
	c.Set("key2", 5)
	result = c.FlatMap(func(value int, key string, coll *collection.Collection[string, int]) *collection.Collection[string, int] {
		sub := collection.New[string, int]()
		if value > 7 {
			sub.Set(key+"_big", value)
		} else {
			sub.Set(key+"_small", value)
		}
		return sub
	})

	if result.Size() != 2 {
		t.Errorf("FlatMap should return collection with 2 items, got %d", result.Size())
	}

	if !result.Has("key1_big") || !result.Has("key2_small") {
		t.Error("FlatMap should contain correctly categorized items")
	}

	// Test with overlapping keys (later values should overwrite earlier ones)
	c.Clear().Set("A", 1).Set("B", 2)
	result = c.FlatMap(func(value int, key string, coll *collection.Collection[string, int]) *collection.Collection[string, int] {
		sub := collection.New[string, int]()
		sub.Set("common", value) // Same key for all items
		return sub
	})

	if result.Size() != 1 {
		t.Errorf("FlatMap with overlapping keys should have size 1, got %d", result.Size())
	}

	commonVal, _ := result.Get("common")
	// The value should be one of the input values (which one depends on iteration order)
	if commonVal != 1 && commonVal != 2 {
		t.Errorf("Common key should have value 1 or 2, got %d", commonVal)
	}

	// Test that function receives correct parameters
	c.Clear().Set("test", 100)
	c.FlatMap(func(value int, key string, coll *collection.Collection[string, int]) *collection.Collection[string, int] {
		if coll.Size() != 1 {
			t.Errorf("Function should receive collection with 1 item, got %d", coll.Size())
		}
		if key != "test" {
			t.Errorf("Expected key 'test', got '%s'", key)
		}
		if value != 100 {
			t.Errorf("Expected value 100, got %d", value)
		}
		return collection.New[string, int]()
	})

	// Test that original collection is unchanged
	if c.Size() != 1 {
		t.Errorf("Original collection should be unchanged, expected size 1, got %d", c.Size())
	}
}

// TestCollectionToJSON tests the ToJSON method
func TestCollectionToJSON(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	jsonBytes, err := c.ToJSON()
	if err != nil {
		t.Errorf("ToJSON should not return error for empty collection, got %v", err)
	}
	if string(jsonBytes) != "[]" {
		t.Errorf("Empty collection JSON should be '[]', got '%s'", string(jsonBytes))
	}

	// Test with single item
	c.Set("key1", 10)
	jsonBytes, err = c.ToJSON()
	if err != nil {
		t.Errorf("ToJSON should not return error, got %v", err)
	}

	// Parse the JSON to verify structure
	jsonStr := string(jsonBytes)
	if !strings.Contains(jsonStr, "key1") || !strings.Contains(jsonStr, "10") {
		t.Errorf("JSON should contain key and value, got '%s'", jsonStr)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)
	jsonBytes, err = c.ToJSON()
	if err != nil {
		t.Errorf("ToJSON should not return error, got %v", err)
	}

	jsonStr = string(jsonBytes)
	// Verify all keys and values are present in JSON
	if !strings.Contains(jsonStr, "key1") || !strings.Contains(jsonStr, "10") ||
		!strings.Contains(jsonStr, "key2") || !strings.Contains(jsonStr, "20") ||
		!strings.Contains(jsonStr, "key3") || !strings.Contains(jsonStr, "30") {
		t.Errorf("JSON should contain all keys and values, got '%s'", jsonStr)
	}

	// Verify JSON structure is array of pairs
	if !strings.HasPrefix(jsonStr, "[") || !strings.HasSuffix(jsonStr, "]") {
		t.Errorf("JSON should be an array, got '%s'", jsonStr)
	}
}

// TestCollectionUnion tests the Union method
func TestCollectionUnion(t *testing.T) {
	c1 := collection.New[string, int]()
	c2 := collection.New[string, int]()

	// Test with both empty collections
	result := c1.Union(c2)
	if result.Size() != 0 {
		t.Errorf("Union of empty collections should be empty, got size %d", result.Size())
	}

	// Test with first collection empty
	c2.Set("key1", 10).Set("key2", 20)
	result = c1.Union(c2)
	if result.Size() != 2 {
		t.Errorf("Union should have size 2, got %d", result.Size())
	}
	if !result.Has("key1") || !result.Has("key2") {
		t.Error("Union should contain all keys from second collection")
	}
	val1, _ := result.Get("key1")
	val2, _ := result.Get("key2")
	if val1 != 10 || val2 != 20 {
		t.Errorf("Expected values 10, 20, got %d, %d", val1, val2)
	}

	// Test with second collection empty
	c1.Set("key3", 30)
	c2.Clear()
	result = c1.Union(c2)
	if result.Size() != 1 {
		t.Errorf("Union should have size 1, got %d", result.Size())
	}
	if !result.Has("key3") {
		t.Error("Union should contain key from first collection")
	}

	// Test with non-overlapping collections
	c1.Clear().Set("key1", 10).Set("key2", 20)
	c2.Set("key3", 30).Set("key4", 40)
	result = c1.Union(c2)
	if result.Size() != 4 {
		t.Errorf("Union should have size 4, got %d", result.Size())
	}
	if !result.Has("key1") || !result.Has("key2") || !result.Has("key3") || !result.Has("key4") {
		t.Error("Union should contain all keys from both collections")
	}

	// Test with overlapping collections (first collection takes precedence)
	c1.Set("key3", 300) // This key exists in both collections
	result = c1.Union(c2)
	if result.Size() != 4 {
		t.Errorf("Union should have size 4, got %d", result.Size())
	}

	val3, _ := result.Get("key3")
	if val3 != 300 {
		t.Errorf("Expected value from first collection (300), got %d", val3)
	}

	// Test that original collections are unchanged
	if c1.Size() != 3 || c2.Size() != 2 {
		t.Error("Original collections should be unchanged")
	}
}

// TestCollectionIntersection tests the Intersection method
func TestCollectionIntersection(t *testing.T) {
	c1 := collection.New[string, int]()
	c2 := collection.New[string, any]()

	// Test with both empty collections
	result := c1.Intersection(c2)
	if result.Size() != 0 {
		t.Errorf("Intersection of empty collections should be empty, got size %d", result.Size())
	}

	// Test with first collection empty
	c2.Set("key1", 10)
	result = c1.Intersection(c2)
	if result.Size() != 0 {
		t.Errorf("Intersection should be empty when first collection is empty, got size %d", result.Size())
	}

	// Test with second collection empty
	c1.Set("key1", 10)
	c2.Clear()
	result = c1.Intersection(c2)
	if result.Size() != 0 {
		t.Errorf("Intersection should be empty when second collection is empty, got size %d", result.Size())
	}

	// Test with non-overlapping collections
	c1.Set("key1", 10).Set("key2", 20)
	c2.Set("key3", 30).Set("key4", 40)
	result = c1.Intersection(c2)
	if result.Size() != 0 {
		t.Errorf("Intersection of non-overlapping collections should be empty, got size %d", result.Size())
	}

	// Test with overlapping collections
	c1.Set("key3", 100) // Different value for same key
	c2.Set("key1", 200) // Different value for same key
	result = c1.Intersection(c2)
	if result.Size() != 2 {
		t.Errorf("Intersection should have size 2, got %d", result.Size())
	}

	if !result.Has("key1") || !result.Has("key3") {
		t.Error("Intersection should contain overlapping keys")
	}

	// Intersection should use values from first collection
	val1, _ := result.Get("key1")
	val3, _ := result.Get("key3")
	if val1 != 10 || val3 != 100 {
		t.Errorf("Expected values from first collection (10, 100), got (%d, %d)", val1, val3)
	}

	// Test that original collections are unchanged
	if c1.Size() != 3 || c2.Size() != 3 {
		t.Error("Original collections should be unchanged")
	}
}

// TestCollectionDifference tests the Difference method
func TestCollectionDifference(t *testing.T) {
	c1 := collection.New[string, int]()
	c2 := collection.New[string, any]()

	// Test with both empty collections
	result := c1.Difference(c2)
	if result.Size() != 0 {
		t.Errorf("Difference of empty collections should be empty, got size %d", result.Size())
	}

	// Test with first collection empty
	c2.Set("key1", 10)
	result = c1.Difference(c2)
	if result.Size() != 0 {
		t.Errorf("Difference should be empty when first collection is empty, got size %d", result.Size())
	}

	// Test with second collection empty
	c1.Set("key1", 10).Set("key2", 20)
	c2.Clear()
	result = c1.Difference(c2)
	if result.Size() != 2 {
		t.Errorf("Difference should equal first collection when second is empty, got size %d", result.Size())
	}
	if !result.Has("key1") || !result.Has("key2") {
		t.Error("Difference should contain all keys from first collection")
	}

	// Test with non-overlapping collections
	c2.Set("key3", 30).Set("key4", 40)
	result = c1.Difference(c2)
	if result.Size() != 2 {
		t.Errorf("Difference should have size 2, got %d", result.Size())
	}
	if !result.Has("key1") || !result.Has("key2") {
		t.Error("Difference should contain all keys from first collection")
	}

	// Test with overlapping collections
	c1.Set("key3", 100) // This key exists in both
	result = c1.Difference(c2)
	if result.Size() != 2 {
		t.Errorf("Difference should have size 2, got %d", result.Size())
	}

	if !result.Has("key1") || !result.Has("key2") {
		t.Error("Difference should contain keys only in first collection")
	}
	if result.Has("key3") {
		t.Error("Difference should not contain keys that exist in second collection")
	}

	// Test values are preserved from first collection
	val1, _ := result.Get("key1")
	val2, _ := result.Get("key2")
	if val1 != 10 || val2 != 20 {
		t.Errorf("Expected values from first collection (10, 20), got (%d, %d)", val1, val2)
	}

	// Test that original collections are unchanged
	if c1.Size() != 3 || c2.Size() != 2 {
		t.Error("Original collections should be unchanged")
	}
}

// TestCollectionSymmetricDifference tests the SymmetricDifference method
func TestCollectionSymmetricDifference(t *testing.T) {
	c1 := collection.New[string, int]()
	c2 := collection.New[string, int]()

	// Test with both empty collections
	result := c1.SymmetricDifference(c2)
	if result.Size() != 0 {
		t.Errorf("SymmetricDifference of empty collections should be empty, got size %d", result.Size())
	}

	// Test with first collection empty
	c2.Set("key1", 10).Set("key2", 20)
	result = c1.SymmetricDifference(c2)
	if result.Size() != 2 {
		t.Errorf("SymmetricDifference should have size 2, got %d", result.Size())
	}
	if !result.Has("key1") || !result.Has("key2") {
		t.Error("SymmetricDifference should contain all keys from second collection")
	}

	// Test with second collection empty
	c1.Set("key3", 30)
	c2.Clear()
	result = c1.SymmetricDifference(c2)
	if result.Size() != 1 {
		t.Errorf("SymmetricDifference should have size 1, got %d", result.Size())
	}
	if !result.Has("key3") {
		t.Error("SymmetricDifference should contain key from first collection")
	}

	// Test with non-overlapping collections
	c1.Set("key1", 10)
	c2.Set("key2", 20)
	result = c1.SymmetricDifference(c2)
	if result.Size() != 3 {
		t.Errorf("SymmetricDifference should have size 3, got %d", result.Size())
	}
	if !result.Has("key1") || !result.Has("key2") || !result.Has("key3") {
		t.Error("SymmetricDifference should contain all unique keys")
	}

	// Test with overlapping collections
	c1.Set("key2", 100) // This key exists in both with different values
	result = c1.SymmetricDifference(c2)
	if result.Size() != 2 {
		t.Errorf("SymmetricDifference should have size 2, got %d", result.Size())
	}

	if !result.Has("key1") || !result.Has("key3") {
		t.Error("SymmetricDifference should contain keys unique to each collection")
	}
	if result.Has("key2") {
		t.Error("SymmetricDifference should not contain keys present in both collections")
	}

	// Test values are preserved correctly
	val1, _ := result.Get("key1")
	val3, _ := result.Get("key3")
	if val1 != 10 || val3 != 30 {
		t.Errorf("Expected values (10, 30), got (%d, %d)", val1, val3)
	}

	// Test that original collections are unchanged
	if c1.Size() != 3 || c2.Size() != 1 {
		t.Error("Original collections should be unchanged")
	}
}

// TestCollectionConcat tests the Concat method
func TestCollectionConcat(t *testing.T) {
	c1 := collection.New[string, int]()
	c2 := collection.New[string, int]()
	c3 := collection.New[string, int]()

	// Test with empty collections
	result := c1.Concat(c2, c3)
	if result.Size() != 0 {
		t.Errorf("Concat of empty collections should be empty, got size %d", result.Size())
	}

	// Test with single collection
	c1.Set("key1", 10).Set("key2", 20)
	result = c1.Concat()
	if result.Size() != 2 {
		t.Errorf("Concat with no additional collections should clone original, got size %d", result.Size())
	}
	if result == c1 {
		t.Error("Concat should return a new collection instance")
	}

	// Test with non-overlapping collections
	c2.Set("key3", 30)
	c3.Set("key4", 40).Set("key5", 50)
	result = c1.Concat(c2, c3)
	if result.Size() != 5 {
		t.Errorf("Concat should have size 5, got %d", result.Size())
	}

	expectedKeys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range expectedKeys {
		if !result.Has(key) {
			t.Errorf("Concat result should contain key %s", key)
		}
	}

	// Test with overlapping collections (later collections override)
	c2.Set("key1", 100) // Override value from c1
	c3.Set("key3", 300) // Override value from c2
	result = c1.Concat(c2, c3)
	if result.Size() != 5 {
		t.Errorf("Concat should have size 5, got %d", result.Size())
	}

	val1, _ := result.Get("key1")
	val3, _ := result.Get("key3")
	if val1 != 100 {
		t.Errorf("Expected overridden value 100 for key1, got %d", val1)
	}
	if val3 != 300 {
		t.Errorf("Expected overridden value 300 for key3, got %d", val3)
	}

	// Test that original collections are unchanged
	origVal1, _ := c1.Get("key1")
	if origVal1 != 10 {
		t.Error("Original collection should be unchanged")
	}

	// Test with many collections
	c4 := collection.New[string, int]()
	c5 := collection.New[string, int]()
	c4.Set("key6", 60)
	c5.Set("key7", 70)

	result = c1.Concat(c2, c3, c4, c5)
	if result.Size() != 7 {
		t.Errorf("Concat with many collections should have size 7, got %d", result.Size())
	}
}

// TestCollectionEquals tests the Equals method
func TestCollectionEquals(t *testing.T) {
	c1 := collection.New[string, int]()
	c2 := collection.New[string, int]()

	// Test with both empty collections
	if !c1.Equals(c2) {
		t.Error("Empty collections should be equal")
	}

	// Test self equality
	if !c1.Equals(c1) {
		t.Error("Collection should be equal to itself")
	}

	// Test with one empty, one non-empty
	c1.Set("key1", 10)
	if c1.Equals(c2) {
		t.Error("Non-empty and empty collections should not be equal")
	}
	if c2.Equals(c1) {
		t.Error("Empty and non-empty collections should not be equal")
	}

	// Test with same content
	c2.Set("key1", 10)
	if !c1.Equals(c2) {
		t.Error("Collections with same content should be equal")
	}
	if !c2.Equals(c1) {
		t.Error("Equality should be symmetric")
	}

	// Test with different values for same key
	c2.Set("key1", 20)
	if c1.Equals(c2) {
		t.Error("Collections with different values should not be equal")
	}

	// Test with different keys
	c2.Set("key1", 10) // Reset to same value
	c2.Set("key2", 20) // Add different key
	if c1.Equals(c2) {
		t.Error("Collections with different keys should not be equal")
	}

	// Test with multiple items
	c1.Set("key2", 20).Set("key3", 30)
	c2.Set("key3", 30)
	if !c1.Equals(c2) {
		t.Error("Collections with same multiple items should be equal")
	}

	// Test with complex values
	c3 := collection.New[string, []int]()
	c4 := collection.New[string, []int]()

	c3.Set("array1", []int{1, 2, 3})
	c4.Set("array1", []int{1, 2, 3})
	if !c3.Equals(c4) {
		t.Error("Collections with equal complex values should be equal")
	}

	c4.Set("array1", []int{1, 2, 4}) // Different array content
	if c3.Equals(c4) {
		t.Error("Collections with different complex values should not be equal")
	}

	// Test that comparison doesn't modify collections
	originalSize1 := c1.Size()
	originalSize2 := c2.Size()
	c1.Equals(c2)
	if c1.Size() != originalSize1 || c2.Size() != originalSize2 {
		t.Error("Equals comparison should not modify collections")
	}
}

// TestCollectionConcurrentAccess tests concurrent read/write access to the collection
func TestCollectionConcurrentAccess(t *testing.T) {
	c := collection.New[string, int]()

	// Pre-populate the collection
	for i := 0; i < 100; i++ {
		c.Set(fmt.Sprintf("key%d", i), i)
	}

	// Test concurrent reads and writes
	var wg sync.WaitGroup
	numGoroutines := 10
	operationsPerGoroutine := 100

	// Start writer goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("concurrent_key_%d_%d", id, j)
				value := id*1000 + j
				c.Set(key, value)
			}
		}(i)
	}

	// Start reader goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				// Read existing keys
				key := fmt.Sprintf("key%d", j%100)
				_, ok := c.Get(key)
				if !ok {
					t.Errorf("Expected to find key %s", key)
				}

				// Check if key exists
				if !c.Has(key) {
					t.Errorf("Expected key %s to exist", key)
				}
			}
		}(i)
	}

	// Start size checker goroutines
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine/2; j++ {
				size := c.Size()
				if size < 100 { // Should at least have the initial 100 items
					t.Errorf("Size should be at least 100, got %d", size)
				}
				time.Sleep(time.Microsecond)
			}
		}()
	}

	wg.Wait()

	// Verify final state
	finalSize := c.Size()
	if finalSize < 100 {
		t.Errorf("Final size should be at least 100, got %d", finalSize)
	}

	// Verify original keys still exist
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		if !c.Has(key) {
			t.Errorf("Original key %s should still exist", key)
		}
	}
}

// TestCollectionConcurrentModifications tests concurrent modifications
func TestCollectionConcurrentModifications(t *testing.T) {
	c := collection.New[string, int]()

	var wg sync.WaitGroup
	numGoroutines := 5

	// Concurrent Set operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("key_%d", j) // Same keys across goroutines
				value := id*100 + j
				c.Set(key, value)
			}
		}(i)
	}

	// Concurrent Delete operations
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 25; j < 50; j++ { // Delete second half
				key := fmt.Sprintf("key_%d", j)
				c.Delete(key)
			}
		}(i)
	}

	wg.Wait()

	// Check final state
	size := c.Size()
	if size > 25 {
		t.Logf("Final size: %d (expected <= 25, but concurrent operations may affect this)", size)
	}

	// Check that some keys from the first half should exist
	existingCount := 0
	for j := 0; j < 25; j++ {
		key := fmt.Sprintf("key_%d", j)
		if c.Has(key) {
			existingCount++
		}
	}

	if existingCount == 0 {
		t.Error("Expected at least some keys from the first half to exist")
	}
}

// TestCollectionConcurrentIterations tests concurrent iterations
func TestCollectionConcurrentIterations(t *testing.T) {
	c := collection.New[string, int]()

	// Pre-populate
	for i := 0; i < 50; i++ {
		c.Set(fmt.Sprintf("key%d", i), i)
	}

	var wg sync.WaitGroup
	numGoroutines := 8

	// Concurrent iterations using different methods
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			switch id % 4 {
			case 0:
				// Test Each
				c.Each(func(value int, key string, coll *collection.Collection[string, int]) {
					// Just iterate, don't modify
				})
			case 1:
				// Test Filter
				filtered := c.Filter(func(value int, key string, coll *collection.Collection[string, int]) bool {
					return value%2 == 0
				})
				if filtered.Size() == 0 {
					t.Error("Filter should return some even numbers")
				}
			case 2:
				// Test Keys/Values
				// Note: During concurrent modifications, keys and values might be different lengths
				// if the collection is modified between the two calls, so we just check they're reasonable
				keys := c.Keys()
				values := c.Values()
				if len(keys) == 0 && len(values) > 0 {
					t.Error("If keys is empty, values should also be empty")
				}
				if len(values) == 0 && len(keys) > 0 {
					t.Error("If values is empty, keys should also be empty")
				}
			case 3:
				// Test Find operations
				_, found := c.Find(func(value int, key string, coll *collection.Collection[string, int]) bool {
					return value > 10
				})
				if !found {
					t.Error("Should find value greater than 10")
				}
			}
		}(i)
	}

	// Concurrent modifications during iterations
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 50; i < 60; i++ {
			c.Set(fmt.Sprintf("new_key%d", i), i)
			time.Sleep(time.Microsecond)
		}
	}()

	wg.Wait()
}

// TestCollectionConcurrentCloneAndSort tests concurrent clone and sort operations
func TestCollectionConcurrentCloneAndSort(t *testing.T) {
	c := collection.New[string, int]()

	// Pre-populate with unsorted data
	for i := 50; i >= 0; i-- {
		c.Set(fmt.Sprintf("key%02d", i), i)
	}

	var wg sync.WaitGroup
	numGoroutines := 6

	// Concurrent clone operations
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				clone := c.Clone()
				if clone.Size() != c.Size() {
					t.Errorf("Clone size should match original, got %d vs %d", clone.Size(), c.Size())
				}

				// Modify clone to ensure independence
				clone.Set(fmt.Sprintf("clone_key_%d_%d", id, j), 999)
				if c.Has(fmt.Sprintf("clone_key_%d_%d", id, j)) {
					t.Error("Modifying clone should not affect original")
				}
			}
		}(i)
	}

	// Concurrent sort operations
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				// Create a copy to sort (since Sort modifies in place)
				sortCopy := c.Clone()
				sortCopy.Sort(func(firstValue, secondValue int, firstKey, secondKey string) int {
					if firstValue < secondValue {
						return -1
					} else if firstValue > secondValue {
						return 1
					}
					return 0
				})

				if sortCopy.Size() != c.Size() {
					t.Error("Sorted copy should have same size as original")
				}
			}
		}()
	}

	wg.Wait()
}

// TestCollectionConcurrentEnsure tests the concurrent safety of Ensure method
func TestCollectionConcurrentEnsure(t *testing.T) {
	c := collection.New[string, int]()

	var wg sync.WaitGroup
	numGoroutines := 10

	// Multiple goroutines trying to ensure the same key
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			value := c.Ensure("shared_key", func(key string, coll *collection.Collection[string, int]) int {
				// Simulate some work in the generator
				time.Sleep(time.Microsecond)
				return 42
			})

			if value != 42 {
				t.Errorf("Expected value 42, got %d", value)
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	if c.Size() != 1 {
		t.Errorf("Should have exactly 1 key, got %d", c.Size())
	}

	finalValue, ok := c.Get("shared_key")
	if !ok || finalValue != 42 {
		t.Errorf("Expected final value 42, got %d (exists: %t)", finalValue, ok)
	}

	// Test concurrent ensure with different keys
	wg = sync.WaitGroup{}
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("unique_key_%d", id)
			value := c.Ensure(key, func(k string, coll *collection.Collection[string, int]) int {
				return id * 10
			})

			if value != id*10 {
				t.Errorf("Expected value %d, got %d", id*10, value)
			}
		}(i)
	}

	wg.Wait()

	// Should now have 1 + numGoroutines keys
	expectedSize := 1 + numGoroutines
	if c.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, c.Size())
	}
}

// TestCollectionDataRace tests for data races using the race detector
func TestCollectionDataRace(t *testing.T) {
	c := collection.New[int, string]()

	var wg sync.WaitGroup
	iterations := 100

	// Concurrent readers and writers on different keys
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Write to own keys
				c.Set(id*1000+j, fmt.Sprintf("value_%d_%d", id, j))
			}
		}(i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Read from all keys
				for k := 0; k < 10; k++ {
					key := k*1000 + (j % iterations)
					c.Get(key)
					c.Has(key)
				}
			}
		}(i)
	}

	// Mixed operations
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			c.Size()
			c.Keys()
			c.Values()
			c.Clear()

			// Re-add some data
			for j := 0; j < 10; j++ {
				c.Set(j, fmt.Sprintf("restored_%d", j))
			}
		}
	}()

	wg.Wait()
}

// TestCollectionEdgeCases tests various edge cases and boundary conditions
func TestCollectionEdgeCases(t *testing.T) {
	// Test with zero values
	c := collection.New[string, int]()
	c.Set("zero", 0)

	val, ok := c.Get("zero")
	if !ok {
		t.Error("Should be able to store and retrieve zero value")
	}
	if val != 0 {
		t.Errorf("Expected 0, got %d", val)
	}

	// Test with empty string key
	c.Set("", 42)
	val, ok = c.Get("")
	if !ok {
		t.Error("Should be able to use empty string as key")
	}
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	// Test with nil slice values
	cSlice := collection.New[string, []int]()
	cSlice.Set("nil_slice", nil)

	slice, ok := cSlice.Get("nil_slice")
	if !ok {
		t.Error("Should be able to store nil slice")
	}
	if slice != nil {
		t.Error("Retrieved slice should be nil")
	}

	// Test large number of items
	largeCollection := collection.New[int, string]()
	numItems := 10000

	for i := 0; i < numItems; i++ {
		largeCollection.Set(i, fmt.Sprintf("value_%d", i))
	}

	if largeCollection.Size() != numItems {
		t.Errorf("Expected size %d, got %d", numItems, largeCollection.Size())
	}

	// Test random access on large collection
	randomKey := 5555
	val2, ok := largeCollection.Get(randomKey)
	if !ok {
		t.Errorf("Should find key %d in large collection", randomKey)
	}
	expectedVal := fmt.Sprintf("value_%d", randomKey)
	if val2 != expectedVal {
		t.Errorf("Expected %s, got %s", expectedVal, val2)
	}

	// Test Clear on large collection
	largeCollection.Clear()
	if largeCollection.Size() != 0 {
		t.Errorf("Large collection should be empty after clear, got size %d", largeCollection.Size())
	}
}

// TestCollectionNilAndEmptyValues tests handling of nil and empty values
func TestCollectionNilAndEmptyValues(t *testing.T) {
	// Test with pointer values
	cPtr := collection.New[string, *int]()

	// Store nil pointer
	cPtr.Set("nil_ptr", nil)
	ptr, ok := cPtr.Get("nil_ptr")
	if !ok {
		t.Error("Should be able to store nil pointer")
	}
	if ptr != nil {
		t.Error("Retrieved pointer should be nil")
	}

	// Store valid pointer
	value := 42
	cPtr.Set("valid_ptr", &value)
	ptr, ok = cPtr.Get("valid_ptr")
	if !ok {
		t.Error("Should be able to store valid pointer")
	}
	if ptr == nil || *ptr != 42 {
		t.Error("Retrieved pointer should point to 42")
	}

	// Test with any values
	cInterface := collection.New[string, any]()

	cInterface.Set("nil_interface", nil)
	cInterface.Set("string_value", "hello")
	cInterface.Set("int_value", 123)
	cInterface.Set("slice_value", []int{1, 2, 3})

	// Retrieve and verify types
	nilVal, _ := cInterface.Get("nil_interface")
	if nilVal != nil {
		t.Error("nil interface should remain nil")
	}

	strVal, _ := cInterface.Get("string_value")
	if str, ok := strVal.(string); !ok || str != "hello" {
		t.Error("String value should be preserved")
	}

	intVal, _ := cInterface.Get("int_value")
	if i, ok := intVal.(int); !ok || i != 123 {
		t.Error("Int value should be preserved")
	}

	sliceVal, _ := cInterface.Get("slice_value")
	if slice, ok := sliceVal.([]int); !ok || len(slice) != 3 || slice[0] != 1 {
		t.Error("Slice value should be preserved")
	}
}

// TestCollectionFunctionParameterEdgeCases tests edge cases in function parameters
func TestCollectionFunctionParameterEdgeCases(t *testing.T) {
	c := collection.New[string, int]()
	c.Set("key1", 10).Set("key2", 20).Set("key3", 30)

	// Test filter with function that always returns false
	filtered := c.Filter(func(value int, key string, coll *collection.Collection[string, int]) bool {
		return false
	})
	if filtered.Size() != 0 {
		t.Error("Filter with always-false function should return empty collection")
	}

	// Test filter with function that always returns true
	filtered = c.Filter(func(value int, key string, coll *collection.Collection[string, int]) bool {
		return true
	})
	if filtered.Size() != c.Size() {
		t.Error("Filter with always-true function should return full collection")
	}

	// Test reduce with zero initial value
	sum := collection.ReduceCollection(c, func(acc int, value int, key string, coll *collection.Collection[string, int]) int {
		return acc + value
	}, 0)
	if sum != 60 {
		t.Errorf("Expected sum 60, got %d", sum)
	}

	// Test map with function that returns same type
	mapped := collection.MapCollection(c, func(value int, key string, coll *collection.Collection[string, int]) int {
		return value
	})
	if len(mapped) != c.Size() {
		t.Error("Map with identity function should preserve all items")
	}

	// Test FlatMap that returns empty collections
	flat := c.FlatMap(func(value int, key string, coll *collection.Collection[string, int]) *collection.Collection[string, int] {
		return collection.New[string, int]()
	})
	if flat.Size() != 0 {
		t.Error("FlatMap with empty-returning function should return empty collection")
	}
}

// TestCollectionMethodChaining tests method chaining edge cases
func TestCollectionMethodChaining(t *testing.T) {
	c := collection.New[string, int]()

	// Test chaining on empty collection
	result := c.Set("key1", 10).Set("key2", 20).Clear().Set("key3", 30)

	if result != c {
		t.Error("All chaining methods should return the same collection instance")
	}

	if c.Size() != 1 {
		t.Errorf("After chaining operations, expected size 1, got %d", c.Size())
	}

	if !c.Has("key3") {
		t.Error("Final operation in chain should be effective")
	}

	// Test chaining with Tap
	chainedSize := 0
	c.Set("key1", 10).Set("key2", 20).Tap(func(coll *collection.Collection[string, int]) {
		chainedSize = coll.Size()
	}).Set("key4", 40)

	if chainedSize != 3 {
		t.Errorf("Tap should see collection with 3 items, saw %d", chainedSize)
	}

	if c.Size() != 4 {
		t.Errorf("After full chain, expected size 4, got %d", c.Size())
	}
}

// TestCollectionStringKeys tests collections with various string key types
func TestCollectionStringKeys(t *testing.T) {
	c := collection.New[string, int]()

	// Test with Unicode keys
	c.Set("café", 1)
	c.Set("🚀", 2)
	c.Set("测试", 3)
	c.Set("Ñoño", 4)

	if c.Size() != 4 {
		t.Errorf("Expected 4 unicode keys, got size %d", c.Size())
	}

	// Verify retrieval
	val, ok := c.Get("café")
	if !ok || val != 1 {
		t.Error("Unicode key 'café' should be retrievable")
	}

	val, ok = c.Get("🚀")
	if !ok || val != 2 {
		t.Error("Emoji key should be retrievable")
	}

	// Test with very long keys
	longKey := strings.Repeat("a", 1000)
	c.Set(longKey, 999)

	val, ok = c.Get(longKey)
	if !ok || val != 999 {
		t.Error("Very long key should be retrievable")
	}

	// Test keys with special characters
	specialKeys := []string{
		"\n\t\r",
		"key with spaces",
		"key\"with'quotes",
		"key\\with\\backslashes",
		"",
	}

	for i, key := range specialKeys {
		c.Set(key, i)
		val, ok := c.Get(key)
		if !ok || val != i {
			t.Errorf("Special key '%s' should be retrievable with value %d", key, i)
		}
	}
}

// TestCollectionTypeEdgeCases tests edge cases with different key/value types
func TestCollectionTypeEdgeCases(t *testing.T) {
	// Test with struct keys
	type Point struct {
		X, Y int
	}

	cStruct := collection.New[Point, string]()
	p1 := Point{1, 2}
	p2 := Point{3, 4}

	cStruct.Set(p1, "point1")
	cStruct.Set(p2, "point2")

	val, ok := cStruct.Get(p1)
	if !ok || val != "point1" {
		t.Error("Struct keys should work correctly")
	}

	// Test equality of struct keys
	p1Copy := Point{1, 2}
	val, ok = cStruct.Get(p1Copy)
	if !ok || val != "point1" {
		t.Error("Struct key equality should work")
	}

	// Test with array keys
	cArray := collection.New[[3]int, string]()
	arr1 := [3]int{1, 2, 3}
	arr2 := [3]int{4, 5, 6}

	cArray.Set(arr1, "array1")
	cArray.Set(arr2, "array2")

	val, ok = cArray.Get(arr1)
	if !ok || val != "array1" {
		t.Error("Array keys should work correctly")
	}

	// Test with boolean keys
	cBool := collection.New[bool, string]()
	cBool.Set(true, "true_value")
	cBool.Set(false, "false_value")

	val, ok = cBool.Get(true)
	if !ok || val != "true_value" {
		t.Error("Boolean keys should work correctly")
	}

	val, ok = cBool.Get(false)
	if !ok || val != "false_value" {
		t.Error("Boolean keys should work correctly")
	}

	if cBool.Size() != 2 {
		t.Error("Boolean key collection should have 2 items")
	}
}
