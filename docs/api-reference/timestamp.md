# timestamp API

Complete API documentation for the timestamp package.

**Import Path:** `github.com/kolosys/atomic/timestamp`

## Package Documentation



## Types

### Nullable
Nullable represents a timestamp that can be null in JSON.

#### Example Usage

```go
// Create a new Nullable
nullable := Nullable{
    Time: Timestamp{},
    Valid: true,
}
```

#### Type Definition

```go
type Nullable struct {
    Time Timestamp
    Valid bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Time | `Timestamp` |  |
| Valid | `bool` |  |

### Constructor Functions

### NewNullable

NewNullable creates a valid Nullable from a Timestamp.

```go
func NewNullable(t Timestamp) Nullable
```

**Parameters:**
- `t` (Timestamp)

**Returns:**
- Nullable

### NullableFromTime

NullableFromTime creates a valid Nullable from a time.Time.

```go
func NullableFromTime(t time.Time) Nullable
```

**Parameters:**
- `t` (time.Time)

**Returns:**
- Nullable

## Methods

### MarshalJSON



```go
func (Nullable) MarshalJSON() ([]byte, error)
```

**Parameters:**
  None

**Returns:**
- []byte
- error

### UnmarshalJSON



```go
func (*Nullable) UnmarshalJSON(data []byte) error
```

**Parameters:**
- `data` ([]byte)

**Returns:**
- error

### Timestamp
Timestamp represents an ISO8601/RFC3339 timestamp with JSON marshaling.

#### Example Usage

```go
// Example usage of Timestamp
var value Timestamp
// Initialize with appropriate value
```

#### Type Definition

```go
type Timestamp time.Time
```

### Constructor Functions

### FromTime

FromTime creates a Timestamp from a time.Time.

```go
func FromTime(t time.Time) Timestamp
```

**Parameters:**
- `t` (time.Time)

**Returns:**
- Timestamp

### FromUnix

FromUnix creates a Timestamp from a Unix timestamp (seconds).

```go
func FromUnix(sec int64) Timestamp
```

**Parameters:**
- `sec` (int64)

**Returns:**
- Timestamp

### FromUnixMilli

FromUnixMilli creates a Timestamp from a Unix timestamp (milliseconds).

```go
func FromUnixMilli(msec int64) Timestamp
```

**Parameters:**
- `msec` (int64)

**Returns:**
- Timestamp

### Now

Now returns the current time as a Timestamp.

```go
func Now() Timestamp
```

**Parameters:**
  None

**Returns:**
- Timestamp

### Parse

Parse parses an RFC3339 string into a Timestamp.

```go
func Parse(s string) (Timestamp, error)
```

**Parameters:**
- `s` (string)

**Returns:**
- Timestamp
- error

## Methods

### IsZero



```go
func (Timestamp) IsZero() bool
```

**Parameters:**
  None

**Returns:**
- bool

### MarshalJSON



```go
func (Nullable) MarshalJSON() ([]byte, error)
```

**Parameters:**
  None

**Returns:**
- []byte
- error

### Time



```go
func (Timestamp) Time() time.Time
```

**Parameters:**
  None

**Returns:**
- time.Time

### Unix



```go
func (Timestamp) Unix() int64
```

**Parameters:**
  None

**Returns:**
- int64

### UnixMilli



```go
func (Timestamp) UnixMilli() int64
```

**Parameters:**
  None

**Returns:**
- int64

### UnmarshalJSON



```go
func (*Nullable) UnmarshalJSON(data []byte) error
```

**Parameters:**
- `data` ([]byte)

**Returns:**
- error

## External Links

- [Package Overview](../packages/timestamp.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/atomic/timestamp)
- [Source Code](https://github.com/kolosys/atomic/tree/main/timestamp)
