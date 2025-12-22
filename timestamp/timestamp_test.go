package timestamp_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kolosys/atomic/timestamp"
)

func TestTimestamp_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		ts       timestamp.Timestamp
		expected string
	}{
		{
			name:     "zero value marshals to null",
			ts:       timestamp.Timestamp{},
			expected: "null",
		},
		{
			name:     "valid timestamp marshals to RFC3339",
			ts:       timestamp.FromTime(time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)),
			expected: `"2024-06-15T12:30:00Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.ts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("got %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantZero  bool
		wantError bool
	}{
		{
			name:     "null unmarshals to zero",
			input:    "null",
			wantZero: true,
		},
		{
			name:     "empty string unmarshals to zero",
			input:    `""`,
			wantZero: true,
		},
		{
			name:     "valid RFC3339 unmarshals correctly",
			input:    `"2024-06-15T12:30:00Z"`,
			wantZero: false,
		},
		{
			name:      "invalid format returns error",
			input:     `"not-a-date"`,
			wantError: true,
		},
		{
			name:      "invalid JSON returns error",
			input:     `{invalid}`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts timestamp.Timestamp
			err := json.Unmarshal([]byte(tt.input), &ts)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantZero && !ts.IsZero() {
				t.Error("expected zero timestamp")
			}
			if !tt.wantZero && ts.IsZero() {
				t.Error("expected non-zero timestamp")
			}
		})
	}
}

func TestTimestamp_RoundTrip(t *testing.T) {
	original := timestamp.FromTime(time.Date(2024, 12, 25, 10, 0, 0, 0, time.UTC))

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed timestamp.Timestamp
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !original.Time().Equal(parsed.Time()) {
		t.Errorf("round trip failed: got %v, want %v", parsed.Time(), original.Time())
	}
}

func TestTimestamp_Methods(t *testing.T) {
	ref := time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC)
	ts := timestamp.FromTime(ref)

	t.Run("Time", func(t *testing.T) {
		if !ts.Time().Equal(ref) {
			t.Errorf("got %v, want %v", ts.Time(), ref)
		}
	})

	t.Run("Unix", func(t *testing.T) {
		if ts.Unix() != ref.Unix() {
			t.Errorf("got %d, want %d", ts.Unix(), ref.Unix())
		}
	})

	t.Run("UnixMilli", func(t *testing.T) {
		if ts.UnixMilli() != ref.UnixMilli() {
			t.Errorf("got %d, want %d", ts.UnixMilli(), ref.UnixMilli())
		}
	})

	t.Run("IsZero", func(t *testing.T) {
		if ts.IsZero() {
			t.Error("expected non-zero")
		}
		var zero timestamp.Timestamp
		if !zero.IsZero() {
			t.Error("expected zero")
		}
	})
}

func TestTimestamp_Constructors(t *testing.T) {
	t.Run("Now", func(t *testing.T) {
		before := time.Now()
		ts := timestamp.Now()
		after := time.Now()

		if ts.Time().Before(before) || ts.Time().After(after) {
			t.Error("Now() returned time outside expected range")
		}
	})

	t.Run("FromUnix", func(t *testing.T) {
		ts := timestamp.FromUnix(1718451045)
		if ts.Unix() != 1718451045 {
			t.Errorf("got %d, want 1718451045", ts.Unix())
		}
	})

	t.Run("FromUnixMilli", func(t *testing.T) {
		ts := timestamp.FromUnixMilli(1718451045123)
		if ts.UnixMilli() != 1718451045123 {
			t.Errorf("got %d, want 1718451045123", ts.UnixMilli())
		}
	})

	t.Run("Parse valid", func(t *testing.T) {
		ts, err := timestamp.Parse("2024-06-15T12:30:45Z")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ts.IsZero() {
			t.Error("expected non-zero timestamp")
		}
	})

	t.Run("Parse invalid", func(t *testing.T) {
		_, err := timestamp.Parse("invalid")
		if err == nil {
			t.Error("expected error for invalid input")
		}
	})
}

func TestNullable_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		nullable timestamp.Nullable
		expected string
	}{
		{
			name:     "invalid marshals to null",
			nullable: timestamp.Nullable{Valid: false},
			expected: "null",
		},
		{
			name:     "valid marshals to RFC3339",
			nullable: timestamp.NewNullable(timestamp.FromTime(time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC))),
			expected: `"2024-06-15T12:30:00Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.nullable)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("got %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestNullable_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
		wantError bool
	}{
		{
			name:      "null unmarshals to invalid",
			input:     "null",
			wantValid: false,
		},
		{
			name:      "valid RFC3339 unmarshals to valid",
			input:     `"2024-06-15T12:30:00Z"`,
			wantValid: true,
		},
		{
			name:      "invalid format returns error",
			input:     `"not-a-date"`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n timestamp.Nullable
			err := json.Unmarshal([]byte(tt.input), &n)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", n.Valid, tt.wantValid)
			}
		})
	}
}

func TestNullable_RoundTrip(t *testing.T) {
	t.Run("valid value", func(t *testing.T) {
		original := timestamp.NewNullable(timestamp.FromTime(time.Date(2024, 12, 25, 10, 0, 0, 0, time.UTC)))

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var parsed timestamp.Nullable
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !parsed.Valid {
			t.Error("expected Valid = true")
		}
		if !original.Time.Time().Equal(parsed.Time.Time()) {
			t.Errorf("round trip failed: got %v, want %v", parsed.Time.Time(), original.Time.Time())
		}
	})

	t.Run("null value", func(t *testing.T) {
		original := timestamp.Nullable{Valid: false}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var parsed timestamp.Nullable
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if parsed.Valid {
			t.Error("expected Valid = false")
		}
	})
}

func TestNullable_Constructors(t *testing.T) {
	t.Run("NewNullable", func(t *testing.T) {
		ts := timestamp.FromTime(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC))
		n := timestamp.NewNullable(ts)

		if !n.Valid {
			t.Error("expected Valid = true")
		}
		if !n.Time.Time().Equal(ts.Time()) {
			t.Error("time mismatch")
		}
	})

	t.Run("NullableFromTime", func(t *testing.T) {
		ref := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
		n := timestamp.NullableFromTime(ref)

		if !n.Valid {
			t.Error("expected Valid = true")
		}
		if !n.Time.Time().Equal(ref) {
			t.Error("time mismatch")
		}
	})
}
