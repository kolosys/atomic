package timestamp

import (
	"encoding/json"
	"time"
)

// Timestamp represents an ISO8601/RFC3339 timestamp with JSON marshaling.
type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*t = Timestamp(parsed)
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(time.Time(t).Format(time.RFC3339))
}

func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

func (t Timestamp) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t Timestamp) Unix() int64 {
	return time.Time(t).Unix()
}

func (t Timestamp) UnixMilli() int64 {
	return time.Time(t).UnixMilli()
}

// Now returns the current time as a Timestamp.
func Now() Timestamp {
	return Timestamp(time.Now())
}

// FromTime creates a Timestamp from a time.Time.
func FromTime(t time.Time) Timestamp {
	return Timestamp(t)
}

// FromUnix creates a Timestamp from a Unix timestamp (seconds).
func FromUnix(sec int64) Timestamp {
	return Timestamp(time.Unix(sec, 0))
}

// FromUnixMilli creates a Timestamp from a Unix timestamp (milliseconds).
func FromUnixMilli(msec int64) Timestamp {
	return Timestamp(time.UnixMilli(msec))
}

// Parse parses an RFC3339 string into a Timestamp.
func Parse(s string) (Timestamp, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp(t), nil
}

// Nullable represents a timestamp that can be null in JSON.
type Nullable struct {
	Time  Timestamp
	Valid bool
}

func (n *Nullable) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		n.Time = Timestamp{}
		return nil
	}
	if err := n.Time.UnmarshalJSON(data); err != nil {
		return err
	}
	n.Valid = true
	return nil
}

func (n Nullable) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return n.Time.MarshalJSON()
}

// NewNullable creates a valid Nullable from a Timestamp.
func NewNullable(t Timestamp) Nullable {
	return Nullable{Time: t, Valid: true}
}

// NullableFromTime creates a valid Nullable from a time.Time.
func NullableFromTime(t time.Time) Nullable {
	return Nullable{Time: Timestamp(t), Valid: true}
}
