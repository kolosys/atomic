package collection_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/kolosys/atomic/collection"
)

// TestMapCollection tests the MapCollection function
func TestMapCollection(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := collection.MapCollection(c, func(value int, key string, collection *collection.Collection[string, int]) string {
		return key + ":" + string(rune(value))
	})
	if len(result) != 0 {
		t.Errorf("MapCollection on empty collection should return empty slice, got %d items", len(result))
	}

	// Test with single item
	c.Set("key1", 10)
	result = collection.MapCollection(c, func(value int, key string, collection *collection.Collection[string, int]) string {
		return key + ":" + string(rune(value+'0'))
	})
	if len(result) != 1 {
		t.Errorf("MapCollection should return slice with 1 item, got %d", len(result))
	}
	expected := "key1:" + string(rune(10+'0'))
	if result[0] != expected {
		t.Errorf("Expected %s, got %s", expected, result[0])
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)
	intResult := collection.MapCollection(c, func(value int, key string, collection *collection.Collection[string, int]) int {
		return value * 2
	})
	if len(intResult) != 3 {
		t.Errorf("MapCollection should return slice with 3 items, got %d", len(intResult))
	}

	// Check that all mapped values are present (order may vary due to map iteration)
	expectedValues := map[int]bool{20: true, 40: true, 60: true}
	for _, val := range intResult {
		if !expectedValues[val] {
			t.Errorf("Unexpected mapped value %d", val)
		}
		delete(expectedValues, val)
	}
	if len(expectedValues) != 0 {
		t.Errorf("Not all expected values were found: %v", expectedValues)
	}

	// Test that function receives correct parameters
	collection.MapCollection(c, func(value int, key string, coll *collection.Collection[string, int]) int {
		if coll.Size() != 3 {
			t.Errorf("Function should receive collection with 3 items, got %d", coll.Size())
		}
		if !coll.Has(key) {
			t.Errorf("Collection should contain key %s", key)
		}
		return value
	})
}

// TestMapCollectionValues tests the MapCollectionValues function
func TestMapCollectionValues(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := collection.MapCollectionValues(c, func(value int, key string, collection *collection.Collection[string, int]) string {
		return key + ":" + string(rune(value))
	})
	if result.Size() != 0 {
		t.Errorf("MapCollectionValues on empty collection should return empty collection, got size %d", result.Size())
	}

	// Test with single item
	c.Set("key1", 10)
	result = collection.MapCollectionValues(c, func(value int, key string, collection *collection.Collection[string, int]) string {
		return strings.ToUpper(key)
	})
	if result.Size() != 1 {
		t.Errorf("MapCollectionValues should return collection with 1 item, got %d", result.Size())
	}
	if !result.Has("key1") {
		t.Error("Result collection should have key1")
	}
	val, _ := result.Get("key1")
	if val != "KEY1" {
		t.Errorf("Expected 'KEY1', got %s", val)
	}

	// Test with multiple items
	c.Set("key2", 20).Set("key3", 30)
	intResult := collection.MapCollectionValues(c, func(value int, key string, collection *collection.Collection[string, int]) int {
		return value * 2
	})
	if intResult.Size() != 3 {
		t.Errorf("MapCollectionValues should return collection with 3 items, got %d", intResult.Size())
	}

	// Check that all keys are preserved with mapped values
	if !intResult.Has("key1") || !intResult.Has("key2") || !intResult.Has("key3") {
		t.Error("Result collection should contain all original keys")
	}

	val1, _ := intResult.Get("key1")
	val2, _ := intResult.Get("key2")
	val3, _ := intResult.Get("key3")

	if val1 != 20 || val2 != 40 || val3 != 60 {
		t.Errorf("Expected mapped values 20, 40, 60, got %d, %d, %d", val1, val2, val3)
	}

	// Test that original collection is unchanged
	if c.Size() != 3 {
		t.Errorf("Original collection should remain unchanged, expected size 3, got %d", c.Size())
	}
	origVal1, _ := c.Get("key1")
	if origVal1 != 10 {
		t.Errorf("Original collection should have unchanged values, expected 10, got %d", origVal1)
	}
}

// TestReduceCollection tests the ReduceCollection function
func TestReduceCollection(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := collection.ReduceCollection(c, func(acc int, value int, key string, collection *collection.Collection[string, int]) int {
		return acc + value
	}, 0)
	if result != 0 {
		t.Errorf("ReduceCollection on empty collection should return initial value 0, got %d", result)
	}

	// Test with different initial value
	result = collection.ReduceCollection(c, func(acc int, value int, key string, collection *collection.Collection[string, int]) int {
		return acc + value
	}, 100)
	if result != 100 {
		t.Errorf("ReduceCollection on empty collection should return initial value 100, got %d", result)
	}

	// Test with single item
	c.Set("key1", 10)
	result = collection.ReduceCollection(c, func(acc int, value int, key string, collection *collection.Collection[string, int]) int {
		return acc + value
	}, 0)
	if result != 10 {
		t.Errorf("ReduceCollection should return 10, got %d", result)
	}

	// Test with multiple items - sum
	c.Set("key2", 20).Set("key3", 30)
	result = collection.ReduceCollection(c, func(acc int, value int, key string, collection *collection.Collection[string, int]) int {
		return acc + value
	}, 0)
	if result != 60 {
		t.Errorf("ReduceCollection sum should return 60, got %d", result)
	}

	// Test with different operation - concatenation
	strResult := collection.ReduceCollection(c, func(acc string, value int, key string, collection *collection.Collection[string, int]) string {
		if acc == "" {
			return key
		}
		return acc + "," + key
	}, "")

	// Result order may vary due to map iteration, so check that all keys are present
	if !strings.Contains(strResult, "key1") || !strings.Contains(strResult, "key2") || !strings.Contains(strResult, "key3") {
		t.Errorf("String result should contain all keys, got %s", strResult)
	}

	// Test with initial value
	result = collection.ReduceCollection(c, func(acc int, value int, key string, collection *collection.Collection[string, int]) int {
		return acc + value
	}, 100)
	if result != 160 {
		t.Errorf("ReduceCollection with initial value 100 should return 160, got %d", result)
	}

	// Test that function receives correct parameters
	collection.ReduceCollection(c, func(acc int, value int, key string, coll *collection.Collection[string, int]) int {
		if coll.Size() != 3 {
			t.Errorf("Function should receive collection with 3 items, got %d", coll.Size())
		}
		if !coll.Has(key) {
			t.Errorf("Collection should contain key %s", key)
		}
		return acc
	}, 0)
}

// TestReduceRightCollection tests the ReduceRightCollection function
func TestReduceRightCollection(t *testing.T) {
	c := collection.New[string, int]()

	// Test with empty collection
	result := collection.ReduceRightCollection(c, func(acc int, value int, key string, collection *collection.Collection[string, int]) int {
		return acc + value
	}, 0)
	if result != 0 {
		t.Errorf("ReduceRightCollection on empty collection should return initial value 0, got %d", result)
	}

	// Test with single item
	c.Set("key1", 10)
	result = collection.ReduceRightCollection(c, func(acc int, value int, key string, collection *collection.Collection[string, int]) int {
		return acc + value
	}, 0)
	if result != 10 {
		t.Errorf("ReduceRightCollection should return 10, got %d", result)
	}

	// Test with multiple items - sum (should be same as regular reduce for commutative operations)
	c.Set("key2", 20).Set("key3", 30)
	result = collection.ReduceRightCollection(c, func(acc int, value int, key string, collection *collection.Collection[string, int]) int {
		return acc + value
	}, 0)
	if result != 60 {
		t.Errorf("ReduceRightCollection sum should return 60, got %d", result)
	}

	// Test with string concatenation to show order difference (though map iteration order isn't guaranteed)
	strResult := collection.ReduceRightCollection(c, func(acc string, value int, key string, collection *collection.Collection[string, int]) string {
		if acc == "" {
			return key
		}
		return acc + "," + key
	}, "")

	// Result order may vary due to map iteration, so just check that all keys are present
	if !strings.Contains(strResult, "key1") || !strings.Contains(strResult, "key2") || !strings.Contains(strResult, "key3") {
		t.Errorf("String result should contain all keys, got %s", strResult)
	}

	// Test that function receives correct parameters
	collection.ReduceRightCollection(c, func(acc int, value int, key string, coll *collection.Collection[string, int]) int {
		if coll.Size() != 3 {
			t.Errorf("Function should receive collection with 3 items, got %d", coll.Size())
		}
		if !coll.Has(key) {
			t.Errorf("Collection should contain key %s", key)
		}
		return acc
	}, 0)
}

// TestMergeCollection tests the MergeCollection function
func TestMergeCollection(t *testing.T) {
	c1 := collection.New[string, int]()
	c2 := collection.New[string, int]()

	// Test with both empty collections
	result := collection.MergeCollection(c1, c2,
		func(value int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value}
		},
		func(valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: valueOther}
		},
		func(value int, valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value + valueOther}
		},
	)
	if result.Size() != 0 {
		t.Errorf("Merge of empty collections should return empty collection, got size %d", result.Size())
	}

	// Test with first collection empty
	c2.Set("key1", 10)
	result = collection.MergeCollection(c1, c2,
		func(value int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value}
		},
		func(valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: valueOther * 2}
		},
		func(value int, valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value + valueOther}
		},
	)
	if result.Size() != 1 {
		t.Errorf("Merge should return collection with 1 item, got %d", result.Size())
	}
	val, _ := result.Get("key1")
	if val != 20 { // valueOther * 2
		t.Errorf("Expected 20, got %d", val)
	}

	// Test with second collection empty
	c1.Clear() // Clear c1 first to start fresh
	c1.Set("key2", 15)
	c2.Clear()
	result = collection.MergeCollection(c1, c2,
		func(value int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value * 3}
		},
		func(valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: valueOther}
		},
		func(value int, valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value + valueOther}
		},
	)
	if result.Size() != 1 {
		t.Errorf("Merge should return collection with 1 item, got %d", result.Size())
	}
	val, _ = result.Get("key2")
	if val != 45 { // value * 3
		t.Errorf("Expected 45, got %d", val)
	}

	// Test with overlapping keys
	c1.Clear().Set("key1", 10).Set("key2", 15) // Ensure we have both keys in c1
	c2.Clear().Set("key1", 20)                 // Only key1 in c2
	result = collection.MergeCollection(c1, c2,
		func(value int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value * 3} // multiply by 3 for items only in c1
		},
		func(valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: valueOther}
		},
		func(value int, valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value + valueOther}
		},
	)
	if result.Size() != 2 {
		t.Errorf("Merge should return collection with 2 items, got %d", result.Size())
	}
	val1, _ := result.Get("key1")
	val2, _ := result.Get("key2")
	if val1 != 30 { // 10 + 20 (both in both collections)
		t.Errorf("Expected 30 for key1, got %d", val1)
	}
	if val2 != 45 { // 15 * 3 (only in c1)
		t.Errorf("Expected 45 for key2, got %d", val2)
	}

	// Test with Keep.Keep = false
	result = collection.MergeCollection(c1, c2,
		func(value int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: false, Value: 0}
		},
		func(valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: key != "key1", Value: valueOther}
		},
		func(value int, valueOther int, key string) collection.Keep[int] {
			return collection.Keep[int]{Keep: true, Value: value * valueOther}
		},
	)
	if result.Size() != 1 {
		t.Errorf("Merge with selective keep should return collection with 1 item, got %d", result.Size())
	}
	val, _ = result.Get("key1")
	if val != 200 { // 10 * 20
		t.Errorf("Expected 200 for key1, got %d", val)
	}
}

// TestCombineEntries tests the CombineEntries function
func TestCombineEntries(t *testing.T) {
	// Test with empty entries
	entries := [][2]any{}
	result := collection.CombineEntries(entries, func(firstValue, secondValue int, key string) int {
		return firstValue + secondValue
	})
	if result.Size() != 0 {
		t.Errorf("CombineEntries with empty entries should return empty collection, got size %d", result.Size())
	}

	// Test with single entry
	entries = [][2]any{{"key1", 10}}
	result = collection.CombineEntries(entries, func(firstValue, secondValue int, key string) int {
		return firstValue + secondValue
	})
	if result.Size() != 1 {
		t.Errorf("CombineEntries should return collection with 1 item, got %d", result.Size())
	}
	val, _ := result.Get("key1")
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	// Test with multiple unique entries
	entries = [][2]any{{"key1", 10}, {"key2", 20}, {"key3", 30}}
	result = collection.CombineEntries(entries, func(firstValue, secondValue int, key string) int {
		return firstValue + secondValue
	})
	if result.Size() != 3 {
		t.Errorf("CombineEntries should return collection with 3 items, got %d", result.Size())
	}

	val1, _ := result.Get("key1")
	val2, _ := result.Get("key2")
	val3, _ := result.Get("key3")
	if val1 != 10 || val2 != 20 || val3 != 30 {
		t.Errorf("Expected values 10, 20, 30, got %d, %d, %d", val1, val2, val3)
	}

	// Test with duplicate keys - combine function should be called
	entries = [][2]any{{"key1", 10}, {"key2", 20}, {"key1", 15}}
	result = collection.CombineEntries(entries, func(firstValue, secondValue int, key string) int {
		return firstValue + secondValue
	})
	if result.Size() != 2 {
		t.Errorf("CombineEntries should return collection with 2 items, got %d", result.Size())
	}

	val1, _ = result.Get("key1")
	val2, _ = result.Get("key2")
	if val1 != 25 { // 10 + 15
		t.Errorf("Expected 25 for key1 (combined), got %d", val1)
	}
	if val2 != 20 {
		t.Errorf("Expected 20 for key2, got %d", val2)
	}

	// Test with multiple duplicates
	entries = [][2]any{{"key1", 5}, {"key1", 10}, {"key1", 15}}
	result = collection.CombineEntries(entries, func(firstValue, secondValue int, key string) int {
		return firstValue * secondValue
	})
	if result.Size() != 1 {
		t.Errorf("CombineEntries should return collection with 1 item, got %d", result.Size())
	}

	val, _ = result.Get("key1")
	if val != 750 { // 5 * 10 * 15 = 750
		t.Errorf("Expected 750 for key1 (combined), got %d", val)
	}
}

// TestGroupBy tests the GroupBy function
func TestGroupBy(t *testing.T) {
	// Test with empty slice
	items := []int{}
	result := collection.GroupBy(items, func(item int, index int) string {
		if item%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if result.Size() != 0 {
		t.Errorf("GroupBy with empty slice should return empty collection, got size %d", result.Size())
	}

	// Test with single item
	items = []int{10}
	result = collection.GroupBy(items, func(item int, index int) string {
		if item%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if result.Size() != 1 {
		t.Errorf("GroupBy should return collection with 1 group, got %d", result.Size())
	}

	evenGroup, _ := result.Get("even")
	if len(evenGroup) != 1 || evenGroup[0] != 10 {
		t.Errorf("Even group should contain [10], got %v", evenGroup)
	}

	// Test with multiple items
	items = []int{1, 2, 3, 4, 5, 6}
	result = collection.GroupBy(items, func(item int, index int) string {
		if item%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if result.Size() != 2 {
		t.Errorf("GroupBy should return collection with 2 groups, got %d", result.Size())
	}

	evenGroup, _ = result.Get("even")
	oddGroup, _ := result.Get("odd")

	expectedEven := []int{2, 4, 6}
	expectedOdd := []int{1, 3, 5}

	if !reflect.DeepEqual(evenGroup, expectedEven) {
		t.Errorf("Even group should be %v, got %v", expectedEven, evenGroup)
	}
	if !reflect.DeepEqual(oddGroup, expectedOdd) {
		t.Errorf("Odd group should be %v, got %v", expectedOdd, oddGroup)
	}

	// Test with custom grouping - group by tens
	items = []int{15, 23, 31, 17, 24, 35}
	result = collection.GroupBy(items, func(item int, index int) string {
		tens := item / 10
		if tens == 1 {
			return "teens"
		} else if tens == 2 {
			return "twenties"
		} else if tens == 3 {
			return "thirties"
		}
		return "other"
	})
	if result.Size() != 3 {
		t.Errorf("GroupBy should return collection with 3 groups, got %d", result.Size())
	}

	teensGroup, _ := result.Get("teens")
	twentiesGroup, _ := result.Get("twenties")
	thirtiesGroup, _ := result.Get("thirties")

	expectedTeens := []int{15, 17}
	expectedTwenties := []int{23, 24}
	expectedThirties := []int{31, 35}

	if !reflect.DeepEqual(teensGroup, expectedTeens) {
		t.Errorf("Teens group should be %v, got %v", expectedTeens, teensGroup)
	}
	if !reflect.DeepEqual(twentiesGroup, expectedTwenties) {
		t.Errorf("Twenties group should be %v, got %v", expectedTwenties, twentiesGroup)
	}
	if !reflect.DeepEqual(thirtiesGroup, expectedThirties) {
		t.Errorf("Thirties group should be %v, got %v", expectedThirties, thirtiesGroup)
	}

	// Test that index parameter is correct
	items = []int{100, 200, 300}
	result = collection.GroupBy(items, func(item int, index int) string {
		return "index_" + string(rune(index+'0'))
	})
	if result.Size() != 3 {
		t.Errorf("GroupBy should return collection with 3 groups, got %d", result.Size())
	}

	group0, _ := result.Get("index_0")
	group1, _ := result.Get("index_1")
	group2, _ := result.Get("index_2")

	if len(group0) != 1 || group0[0] != 100 {
		t.Errorf("Index 0 group should contain [100], got %v", group0)
	}
	if len(group1) != 1 || group1[0] != 200 {
		t.Errorf("Index 1 group should contain [200], got %v", group1)
	}
	if len(group2) != 1 || group2[0] != 300 {
		t.Errorf("Index 2 group should contain [300], got %v", group2)
	}
}
