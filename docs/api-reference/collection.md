# collection API

Complete API documentation for the collection package.

**Import Path:** `github.com/kolosys/atomic/collection`

## Package Documentation



## Types

### Collection
Collection is a generic map-like structure with additional utility methods. It is safe for concurrent use.

#### Example Usage

```go
// Create a new Collection
collection := Collection{

}
```

#### Type Definition

```go
type Collection struct {
}
```

### Constructor Functions

### CombineEntries

CombineEntries creates a Collection from a list of entries.

```go
func CombineEntries(entries [][*ast.BasicLit]any, combine func(firstValue, secondValue V, key K) V) **ast.IndexListExpr
```

**Parameters:**
- `entries` ([][*ast.BasicLit]any)
- `combine` (func(firstValue, secondValue V, key K) V)

**Returns:**
- **ast.IndexListExpr

### GroupBy

GroupBy groups items by a key selector function.

```go
func GroupBy(items []Item, keySelector func(item Item, index int) K) **ast.IndexListExpr
```

**Parameters:**
- `items` ([]Item)
- `keySelector` (func(item Item, index int) K)

**Returns:**
- **ast.IndexListExpr

### MapCollectionValues

MapValues returns a new collection with the same keys but values mapped by fn.

```go
func MapCollectionValues(c **ast.IndexListExpr, fn func(value V, key K, collection **ast.IndexListExpr) R) **ast.IndexListExpr
```

**Parameters:**
- `c` (**ast.IndexListExpr)
- `fn` (func(value V, key K, collection **ast.IndexListExpr) R)

**Returns:**
- **ast.IndexListExpr

### MergeCollection

Merge merges two collections together into a new collection.

```go
func MergeCollection(c **ast.IndexListExpr, other **ast.IndexListExpr, whenInSelf func(value V, key K) *ast.IndexExpr, whenInOther func(valueOther O, key K) *ast.IndexExpr, whenInBoth func(value V, valueOther O, key K) *ast.IndexExpr) **ast.IndexListExpr
```

**Parameters:**
- `c` (**ast.IndexListExpr)
- `other` (**ast.IndexListExpr)
- `whenInSelf` (func(value V, key K) *ast.IndexExpr)
- `whenInOther` (func(valueOther O, key K) *ast.IndexExpr)
- `whenInBoth` (func(value V, valueOther O, key K) *ast.IndexExpr)

**Returns:**
- **ast.IndexListExpr

### New

New creates a new Collection.

```go
func New() **ast.IndexListExpr
```

**Parameters:**
  None

**Returns:**
- **ast.IndexListExpr

## Methods

### At

At returns the value at a given index, allowing for positive and negative integers.

```go
func (**ast.IndexListExpr) At(index int) (V, bool)
```

**Parameters:**
- `index` (int)

**Returns:**
- V
- bool

### Clear

Clear removes all items from the collection.

```go
func (**ast.IndexListExpr) Clear() **ast.IndexListExpr
```

**Parameters:**
  None

**Returns:**
- **ast.IndexListExpr

### Clone

Clone creates a shallow copy of the collection.

```go
func (**ast.IndexListExpr) Clone() **ast.IndexListExpr
```

**Parameters:**
  None

**Returns:**
- **ast.IndexListExpr

### Concat

Concat combines this collection with others into a new collection.

```go
func (**ast.IndexListExpr) Concat(collections ...**ast.IndexListExpr) **ast.IndexListExpr
```

**Parameters:**
- `collections` (...**ast.IndexListExpr)

**Returns:**
- **ast.IndexListExpr

### Delete

Delete removes an item from the collection.

```go
func (**ast.IndexListExpr) Delete(key K) bool
```

**Parameters:**
- `key` (K)

**Returns:**
- bool

### Difference

Difference returns a new collection containing the items where the key is present in this collection but not the other.

```go
func (**ast.IndexListExpr) Difference(other **ast.IndexListExpr) **ast.IndexListExpr
```

**Parameters:**
- `other` (**ast.IndexListExpr)

**Returns:**
- **ast.IndexListExpr

### Each

Each executes fn for each element and returns the collection.

```go
func (**ast.IndexListExpr) Each(fn func(value V, key K, collection **ast.IndexListExpr)) **ast.IndexListExpr
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr))

**Returns:**
- **ast.IndexListExpr

### Ensure

Ensure obtains the value for the given key if it exists, otherwise sets and returns the value provided by the default value generator.

```go
func (**ast.IndexListExpr) Ensure(key K, defaultValueGenerator func(key K, collection **ast.IndexListExpr) V) V
```

**Parameters:**
- `key` (K)
- `defaultValueGenerator` (func(key K, collection **ast.IndexListExpr) V)

**Returns:**
- V

### Entries

Entries returns all key-value pairs in the collection.

```go
func (**ast.IndexListExpr) Entries() [][*ast.BasicLit]any
```

**Parameters:**
  None

**Returns:**
- [][*ast.BasicLit]any

### Equals

Equals checks if this collection shares identical items with another.

```go
func (**ast.IndexListExpr) Equals(other **ast.IndexListExpr) bool
```

**Parameters:**
- `other` (**ast.IndexListExpr)

**Returns:**
- bool

### Every

Every returns true if all items pass the test.

```go
func (**ast.IndexListExpr) Every(fn func(value V, key K, collection **ast.IndexListExpr) bool) bool
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- bool

### Filter

Filter returns a new collection containing only the items for which fn returns true.

```go
func (**ast.IndexListExpr) Filter(fn func(value V, key K, collection **ast.IndexListExpr) bool) **ast.IndexListExpr
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- **ast.IndexListExpr

### Find

Find returns the first value for which fn returns true.

```go
func (**ast.IndexListExpr) Find(fn func(value V, key K, collection **ast.IndexListExpr) bool) (V, bool)
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- V
- bool

### FindKey

FindKey returns the first key for which fn returns true.

```go
func (**ast.IndexListExpr) FindKey(fn func(value V, key K, collection **ast.IndexListExpr) bool) (K, bool)
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- K
- bool

### FindLast

FindLast returns the last value for which fn returns true.

```go
func (**ast.IndexListExpr) FindLast(fn func(value V, key K, collection **ast.IndexListExpr) bool) (V, bool)
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- V
- bool

### FindLastKey

FindLastKey returns the last key for which fn returns true.

```go
func (**ast.IndexListExpr) FindLastKey(fn func(value V, key K, collection **ast.IndexListExpr) bool) (K, bool)
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- K
- bool

### First

First returns the first value(s) in the collection. If amount is 0, returns nil. If amount < 0, returns Last(-amount).

```go
func (**ast.IndexListExpr) First(amount ...int) any
```

**Parameters:**
- `amount` (...int)

**Returns:**
- any

### FirstKey

FirstKey returns the first key(s) in the collection.

```go
func (**ast.IndexListExpr) FirstKey(amount ...int) any
```

**Parameters:**
- `amount` (...int)

**Returns:**
- any

### FlatMap

FlatMap maps each item into a collection, then joins the results into a single collection.

```go
func (**ast.IndexListExpr) FlatMap(fn func(value V, key K, collection **ast.IndexListExpr) **ast.IndexListExpr) **ast.IndexListExpr
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) **ast.IndexListExpr)

**Returns:**
- **ast.IndexListExpr

### Get

Get retrieves an item from the collection.

```go
func (**ast.IndexListExpr) Get(key K) (V, bool)
```

**Parameters:**
- `key` (K)

**Returns:**
- V
- bool

### Has

Has checks if a key exists in the collection.

```go
func (**ast.IndexListExpr) Has(key K) bool
```

**Parameters:**
- `key` (K)

**Returns:**
- bool

### HasAll

HasAll checks if all of the provided keys exist in the collection.

```go
func (**ast.IndexListExpr) HasAll(keys ...K) bool
```

**Parameters:**
- `keys` (...K)

**Returns:**
- bool

### HasAny

HasAny checks if any of the provided keys exist in the collection.

```go
func (**ast.IndexListExpr) HasAny(keys ...K) bool
```

**Parameters:**
- `keys` (...K)

**Returns:**
- bool

### Intersection

Intersection returns a new collection containing the items where the key is present in both collections.

```go
func (**ast.IndexListExpr) Intersection(other **ast.IndexListExpr) **ast.IndexListExpr
```

**Parameters:**
- `other` (**ast.IndexListExpr)

**Returns:**
- **ast.IndexListExpr

### KeyAt

KeyAt returns the key at a given index, allowing for positive and negative integers.

```go
func (**ast.IndexListExpr) KeyAt(index int) (K, bool)
```

**Parameters:**
- `index` (int)

**Returns:**
- K
- bool

### Keys

Keys returns all keys in the collection.

```go
func (**ast.IndexListExpr) Keys() []K
```

**Parameters:**
  None

**Returns:**
- []K

### Last

Last returns the last value(s) in the collection.

```go
func (**ast.IndexListExpr) Last(amount ...int) any
```

**Parameters:**
- `amount` (...int)

**Returns:**
- any

### LastKey

LastKey returns the last key(s) in the collection.

```go
func (**ast.IndexListExpr) LastKey(amount ...int) any
```

**Parameters:**
- `amount` (...int)

**Returns:**
- any

### Partition

Partition splits the collection into two collections: the first contains items that passed, the second those that failed.

```go
func (**ast.IndexListExpr) Partition(fn func(value V, key K, collection **ast.IndexListExpr) bool) (**ast.IndexListExpr, **ast.IndexListExpr)
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- **ast.IndexListExpr
- **ast.IndexListExpr

### Random

Random returns a random value or n unique random values from the collection.

```go
func (**ast.IndexListExpr) Random(amount ...int) any
```

**Parameters:**
- `amount` (...int)

**Returns:**
- any

### RandomKey

RandomKey returns a random key or n unique random keys from the collection.

```go
func (**ast.IndexListExpr) RandomKey(amount ...int) any
```

**Parameters:**
- `amount` (...int)

**Returns:**
- any

### Reverse

Reverse reverses the order of the collection in place.

```go
func (**ast.IndexListExpr) Reverse() **ast.IndexListExpr
```

**Parameters:**
  None

**Returns:**
- **ast.IndexListExpr

### Set

Set adds or updates an item in the collection.

```go
func (**ast.IndexListExpr) Set(key K, value V) **ast.IndexListExpr
```

**Parameters:**
- `key` (K)
- `value` (V)

**Returns:**
- **ast.IndexListExpr

### Size

Size returns the number of items in the collection.

```go
func (**ast.IndexListExpr) Size() int
```

**Parameters:**
  None

**Returns:**
- int

### Some

Some returns true if any item passes the test.

```go
func (**ast.IndexListExpr) Some(fn func(value V, key K, collection **ast.IndexListExpr) bool) bool
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- bool

### Sort

Sort sorts the items of a collection in place and returns it.

```go
func (**ast.IndexListExpr) Sort(compare *ast.IndexListExpr) **ast.IndexListExpr
```

**Parameters:**
- `compare` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexListExpr

### Sweep

Sweep removes items that satisfy the provided filter function. Returns the number of removed entries.

```go
func (**ast.IndexListExpr) Sweep(fn func(value V, key K, collection **ast.IndexListExpr) bool) int
```

**Parameters:**
- `fn` (func(value V, key K, collection **ast.IndexListExpr) bool)

**Returns:**
- int

### SymmetricDifference

SymmetricDifference returns a new collection containing only the items where the keys are present in either collection, but not both.

```go
func (**ast.IndexListExpr) SymmetricDifference(other **ast.IndexListExpr) **ast.IndexListExpr
```

**Parameters:**
- `other` (**ast.IndexListExpr)

**Returns:**
- **ast.IndexListExpr

### Tap

Tap runs a function on the collection and returns the collection.

```go
func (**ast.IndexListExpr) Tap(fn func(collection **ast.IndexListExpr)) **ast.IndexListExpr
```

**Parameters:**
- `fn` (func(collection **ast.IndexListExpr))

**Returns:**
- **ast.IndexListExpr

### ToJSON

ToJSON returns the collection as a JSON array of [key, value] pairs.

```go
func (**ast.IndexListExpr) ToJSON() ([]byte, error)
```

**Parameters:**
  None

**Returns:**
- []byte
- error

### ToReversed

ToReversed returns a new collection with the items in reverse order.

```go
func (**ast.IndexListExpr) ToReversed() **ast.IndexListExpr
```

**Parameters:**
  None

**Returns:**
- **ast.IndexListExpr

### ToSorted

ToSorted returns a shallow copy of the collection with the items sorted.

```go
func (**ast.IndexListExpr) ToSorted(compare *ast.IndexListExpr) **ast.IndexListExpr
```

**Parameters:**
- `compare` (*ast.IndexListExpr)

**Returns:**
- **ast.IndexListExpr

### Union

Union returns a new collection containing the items where the key is present in either of the collections.

```go
func (**ast.IndexListExpr) Union(other **ast.IndexListExpr) **ast.IndexListExpr
```

**Parameters:**
- `other` (**ast.IndexListExpr)

**Returns:**
- **ast.IndexListExpr

### Values

Values returns all values in the collection.

```go
func (**ast.IndexListExpr) Values() []V
```

**Parameters:**
  None

**Returns:**
- []V

### Comparator
Comparator is a function that compares two values and their keys, returning -1, 0, or 1.

#### Example Usage

```go
// Example usage of Comparator
var value Comparator
// Initialize with appropriate value
```

#### Type Definition

```go
type Comparator func(firstValue, secondValue V, firstKey, secondKey K) int
```

### Keep
Keep is used for merge operations to indicate whether to keep a value and what value to keep.

#### Example Usage

```go
// Create a new Keep
keep := Keep{
    Keep: true,
    Value: V{},
}
```

#### Type Definition

```go
type Keep struct {
    Keep bool
    Value V
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Keep | `bool` |  |
| Value | `V` |  |

## Functions

### DefaultSort
DefaultSort is the default sort comparison algorithm used in ECMAScript.

```go
func DefaultSort(firstValue, secondValue V, firstKey, secondKey K) int
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `firstValue` | `V` | |
| `secondValue` | `V` | |
| `firstKey` | `K` | |
| `secondKey` | `K` | |

**Returns:**
| Type | Description |
|------|-------------|
| `int` | |

**Example:**

```go
// Example usage of DefaultSort
result := DefaultSort(/* parameters */)
```

### MapCollection
Map returns a slice of values produced by applying fn to each item.

```go
func MapCollection(c **ast.IndexListExpr, fn func(value V, key K, collection **ast.IndexListExpr) R) []R
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `c` | `**ast.IndexListExpr` | |
| `fn` | `func(value V, key K, collection **ast.IndexListExpr) R` | |

**Returns:**
| Type | Description |
|------|-------------|
| `[]R` | |

**Example:**

```go
// Example usage of MapCollection
result := MapCollection(/* parameters */)
```

### ReduceCollection
Reduce applies a function to produce a single value.

```go
func ReduceCollection(c **ast.IndexListExpr, fn func(accumulator R, value V, key K, collection **ast.IndexListExpr) R, initialValue R) R
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `c` | `**ast.IndexListExpr` | |
| `fn` | `func(accumulator R, value V, key K, collection **ast.IndexListExpr) R` | |
| `initialValue` | `R` | |

**Returns:**
| Type | Description |
|------|-------------|
| `R` | |

**Example:**

```go
// Example usage of ReduceCollection
result := ReduceCollection(/* parameters */)
```

### ReduceRightCollection
ReduceRight applies a function to produce a single value, iterating from the end.

```go
func ReduceRightCollection(c **ast.IndexListExpr, fn func(accumulator R, value V, key K, collection **ast.IndexListExpr) R, initialValue R) R
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `c` | `**ast.IndexListExpr` | |
| `fn` | `func(accumulator R, value V, key K, collection **ast.IndexListExpr) R` | |
| `initialValue` | `R` | |

**Returns:**
| Type | Description |
|------|-------------|
| `R` | |

**Example:**

```go
// Example usage of ReduceRightCollection
result := ReduceRightCollection(/* parameters */)
```

## External Links

- [Package Overview](../packages/collection.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/atomic/collection)
- [Source Code](https://github.com/kolosys/atomic/tree/main/collection)
