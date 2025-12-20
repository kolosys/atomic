package snowflake_test

import (
	"testing"
	"time"

	. "github.com/kolosys/atomic/snowflake"
)

func TestSnowflakeStringAndInt64(t *testing.T) {
	val := int64(123456789012345678)
	sf := New(val, DefaultEpoch)
	if got := sf.String(); got != "123456789012345678" {
		t.Errorf("String() = %s, want %s", got, "123456789012345678")
	}
	gotInt, err := sf.Int64()
	if err != nil {
		t.Errorf("Int64() returned error: %v", err)
	}
	if gotInt != val {
		t.Errorf("Int64() = %d, want %d", gotInt, val)
	}
}

func TestSnowflakeNewFromString(t *testing.T) {
	val := "9876543210123456"
	sf := NewFromString(val, DefaultEpoch)
	expected, _ := NewFromString(val, DefaultEpoch).Int64()
	got, err := sf.Int64()
	if err != nil {
		t.Errorf("Int64() returned error: %v", err)
	}
	if got != expected {
		t.Errorf("NewFromString(%s) = %d, want %d", val, got, expected)
	}

	bad := NewFromString("notanumber", DefaultEpoch)
	gotBad, _ := bad.Int64()
	if gotBad != 0 {
		t.Errorf("NewFromString with bad string should be zero, got %d", gotBad)
	}
}

func TestNewRandom(t *testing.T) {
	sf1 := NewRandom(DefaultEpoch)
	sf2 := NewRandom(DefaultEpoch)
	if sf1.String() == sf2.String() {
		t.Errorf("NewRandom() generated two identical snowflakes: %s", sf1.String())
	}
}

func TestSnowflakeTime(t *testing.T) {
	// Compose a simple snowflake id with now timestamp in ms
	nowMs := time.Now().UnixMilli()
	// mimic snowflake: ((timestamp - epoch) << 22)
	id := (nowMs - DefaultEpoch) << 22
	sf := New(id, DefaultEpoch)
	tm, err := sf.Time()
	if err != nil {
		t.Fatalf("Time() returned error: %v", err)
	}
	// Allow for a few seconds difference
	timeout := time.Duration(5) * time.Second
	if tm.Sub(time.UnixMilli(nowMs)) > timeout {
		t.Errorf("Extracted time too far from now: got %v, want ~%v", tm, time.UnixMilli(nowMs))
	}
}

func TestIsValid(t *testing.T) {
	sf := New(42, DefaultEpoch)
	if !sf.IsValid() {
		t.Errorf("IsValid() = false, want true")
	}
	sf0 := New(0, DefaultEpoch)
	if sf0.IsValid() {
		t.Errorf("IsValid() = true, want false")
	}
}
