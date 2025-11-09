# Collection

A powerful, generic, thread-safe map-like data structure for Go with rich utility methods inspired by JavaScript's Map and Array methods.

## Features

- **Thread-Safe**: All operations are protected by read-write mutexes for concurrent use
- **Generic**: Works with any comparable key type and any value type
- **Functional Programming**: Includes map, filter, reduce, and other functional utilities
- **Set Operations**: Union, intersection, difference, and symmetric difference
- **Chainable API**: Many methods return the collection for method chaining
- **Rich Query Methods**: Find, filter, partition, and test operations
- **Array-Like Access**: Get items by index with positive/negative indexing
- **Sorting & Ordering**: Sort, reverse, and randomize collection items

## Installation

```bash
go get github.com/kolosys/atomic/collection
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/kolosys/atomic/collection"
)

func main() {
    // Create a new collection
    users := collection.New[string, User]()

    // Add items
    users.Set("alice", User{Name: "Alice", Age: 30})
    users.Set("bob", User{Name: "Bob", Age: 25})

    // Get items
    if user, ok := users.Get("alice"); ok {
        fmt.Printf("Found user: %s\n", user.Name)
    }

    // Check existence
    if users.Has("alice") {
        fmt.Println("Alice exists!")
    }

    // Get collection size
    fmt.Printf("Total users: %d\n", users.Size())
}
```

## Basic Operations

### Creating a Collection

```go
// Create an empty collection
c := collection.New[string, int]()

// Add items
c.Set("one", 1)
c.Set("two", 2)
c.Set("three", 3)
```

### Getting and Checking Values

```go
// Get a value
value, exists := c.Get("one")
if exists {
    fmt.Println(value) // 1
}

// Check if key exists
if c.Has("one") {
    fmt.Println("Key exists")
}

// Check multiple keys
if c.HasAll("one", "two") {
    fmt.Println("All keys exist")
}

if c.HasAny("one", "four") {
    fmt.Println("At least one key exists")
}
```

### Modifying Collections

```go
// Delete an item
existed := c.Delete("one") // returns true if key existed

// Clear all items
c.Clear()

// Ensure a value exists (get or set)
value := c.Ensure("key", func(key string, coll *collection.Collection[string, int]) int {
    return 42 // default value if key doesn't exist
})
```

## Collection Information

```go
// Get size
size := c.Size()

// Get all keys
keys := c.Keys() // []K

// Get all values
values := c.Values() // []V

// Get all entries as [key, value] pairs
entries := c.Entries() // [][2]any
```

## Iteration and Traversal

### Each

```go
// Execute function for each element
c.Each(func(value int, key string, coll *collection.Collection[string, int]) {
    fmt.Printf("%s: %d\n", key, value)
})
```

### Map

```go
// Map to a slice
doubled := collection.MapCollection(c, func(value int, key string, coll *collection.Collection[string, int]) int {
    return value * 2
})
// Result: []int with all values doubled
```

### MapValues

```go
// Map to a new collection with transformed values
squares := collection.MapCollectionValues(c, func(value int, key string, coll *collection.Collection[string, int]) int {
    return value * value
})
// Result: *Collection[string, int] with squared values
```

### Reduce

```go
// Reduce to a single value
sum := collection.ReduceCollection(c, func(acc int, value int, key string, coll *collection.Collection[string, int]) int {
    return acc + value
}, 0)
```

## Filtering and Searching

### Filter

```go
// Create new collection with filtered items
evens := c.Filter(func(value int, key string, coll *collection.Collection[string, int]) bool {
    return value%2 == 0
})
```

### Sweep

```go
// Remove items in place (returns count removed)
removed := c.Sweep(func(value int, key string, coll *collection.Collection[string, int]) bool {
    return value < 10 // remove values less than 10
})
```

### Find

```go
// Find first matching value
value, found := c.Find(func(value int, key string, coll *collection.Collection[string, int]) bool {
    return value > 50
})

// Find first matching key
key, found := c.FindKey(func(value int, key string, coll *collection.Collection[string, int]) bool {
    return value > 50
})

// Find last matching value/key
lastValue, found := c.FindLast(...)
lastKey, found := c.FindLastKey(...)
```

### Partition

```go
// Split into two collections based on predicate
pass, fail := c.Partition(func(value int, key string, coll *collection.Collection[string, int]) bool {
    return value%2 == 0
})
// pass: even values, fail: odd values
```

### Test Operations

```go
// Check if some items match
hasPositive := c.Some(func(value int, key string, coll *collection.Collection[string, int]) bool {
    return value > 0
})

// Check if all items match
allPositive := c.Every(func(value int, key string, coll *collection.Collection[string, int]) bool {
    return value > 0
})
```

## Array-Like Access

### Accessing by Index

```go
// Get first element(s)
first := c.First()          // Returns single value
firstThree := c.First(3)    // Returns []V with up to 3 values
firstKey := c.FirstKey()    // Returns single key
firstKeys := c.FirstKey(3)  // Returns []K with up to 3 keys

// Get last element(s)
last := c.Last()            // Returns single value
lastThree := c.Last(3)      // Returns []V with up to 3 values
lastKey := c.LastKey()      // Returns single key
lastKeys := c.LastKey(3)    // Returns []K with up to 3 keys

// Access by index (supports negative indices)
value, ok := c.At(0)        // First element
value, ok = c.At(-1)        // Last element
key, ok := c.KeyAt(2)       // Third key
```

### Random Selection

```go
// Get random value
random := c.Random()        // Returns single random value
randoms := c.Random(3)      // Returns []V with up to 3 unique random values

// Get random key
randomKey := c.RandomKey()  // Returns single random key
randomKeys := c.RandomKey(3) // Returns []K with up to 3 unique random keys
```

## Sorting and Ordering

### Sort

```go
// Sort in place
c.Sort(func(v1, v2 int, k1, k2 string) int {
    if v1 < v2 {
        return -1
    } else if v1 > v2 {
        return 1
    }
    return 0
})

// Or use the default sort (string comparison)
c.Sort(collection.DefaultSort[string, int])

// Create sorted copy
sorted := c.ToSorted(collection.DefaultSort[string, int])
```

### Reverse

```go
// Reverse in place
c.Reverse()

// Create reversed copy
reversed := c.ToReversed()
```

## Set Operations

### Union

```go
// Items present in either collection
union := c1.Union(c2)
```

### Intersection

```go
// Items with keys present in both collections
intersection := c1.Intersection(c2)
```

### Difference

```go
// Items in c1 but not in c2
difference := c1.Difference(c2)
```

### Symmetric Difference

```go
// Items in either collection but not both
symDiff := c1.SymmetricDifference(c2)
```

## Advanced Operations

### Clone

```go
// Create a shallow copy
clone := c.Clone()
```

### Concat

```go
// Combine multiple collections
combined := c1.Concat(c2, c3, c4)
```

### FlatMap

```go
// Map each item to a collection, then flatten
result := c.FlatMap(func(value int, key string, coll *collection.Collection[string, int]) *collection.Collection[string, int] {
    nested := collection.New[string, int]()
    nested.Set(key+"_1", value*1)
    nested.Set(key+"_2", value*2)
    return nested
})
```

### Merge

```go
// Advanced merge with control over which values to keep
merged := collection.MergeCollection(
    c1,
    c2,
    func(v1 int, key string) collection.Keep[int] {
        // When key only in c1
        return collection.Keep[int]{Keep: true, Value: v1}
    },
    func(v2 int, key string) collection.Keep[int] {
        // When key only in c2
        return collection.Keep[int]{Keep: true, Value: v2}
    },
    func(v1, v2 int, key string) collection.Keep[int] {
        // When key in both - keep the larger value
        if v1 > v2 {
            return collection.Keep[int]{Keep: true, Value: v1}
        }
        return collection.Keep[int]{Keep: true, Value: v2}
    },
)
```

### Equals

```go
// Check if two collections have identical items
areEqual := c1.Equals(c2)
```

### Tap

```go
// Execute a function on the collection and return it (useful for debugging)
c.Tap(func(coll *collection.Collection[string, int]) {
    fmt.Printf("Collection size: %d\n", coll.Size())
}).Set("key", 42)
```

## Utility Functions

### GroupBy

```go
// Group items by a key selector
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {Name: "Alice", Age: 30},
    {Name: "Bob", Age: 25},
    {Name: "Charlie", Age: 30},
}

grouped := collection.GroupBy(people, func(person Person, index int) int {
    return person.Age
})
// Result: *Collection[int, []Person] grouped by age
```

### CombineEntries

```go
// Create collection from entries with duplicate key handling
entries := [][2]any{
    {"key1", 10},
    {"key2", 20},
    {"key1", 5}, // duplicate key
}

c := collection.CombineEntries[string, int](entries, func(first, second int, key string) int {
    return first + second // combine values for duplicate keys
})
```

### ToJSON

```go
// Export as JSON array of [key, value] pairs
jsonData, err := c.ToJSON()
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(jsonData))
```

## Thread Safety

All Collection operations are thread-safe and can be used concurrently:

```go
c := collection.New[string, int]()

var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        c.Set(fmt.Sprintf("key%d", n), n)
    }(i)
}
wg.Wait()

fmt.Printf("Final size: %d\n", c.Size())
```

## Method Chaining

Many methods return the collection itself, allowing for fluent method chaining:

```go
result := collection.New[string, int]().
    Set("one", 1).
    Set("two", 2).
    Set("three", 3).
    Tap(func(c *collection.Collection[string, int]) {
        fmt.Printf("Size: %d\n", c.Size())
    }).
    Each(func(value int, key string, c *collection.Collection[string, int]) {
        fmt.Printf("%s: %d\n", key, value)
    })
```

## Performance Considerations

- **Read Operations**: Protected by `RWMutex.RLock()`, allowing concurrent reads
- **Write Operations**: Protected by `RWMutex.Lock()`, ensuring exclusive access
- **Memory**: Shallow copies are created by `Clone()` - the values themselves are not deep copied
- **Ordering**: Go maps are unordered, so iteration order is not guaranteed unless sorted

## License

This package is part of the Kolosys Atomic project.

## Contributing

Contributions are welcome! Please ensure all tests pass before submitting a pull request.

```bash
go test -v
```
