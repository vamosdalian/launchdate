package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFlexibleTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		jsonValue  string
		wantError  bool
		wantResult string // Expected time in RFC3339 format
	}{
		{
			name:       "RFC3339 with seconds and timezone",
			jsonValue:  `"2025-10-26T14:05:00Z"`,
			wantError:  false,
			wantResult: "2025-10-26T14:05:00Z",
		},
		{
			name:       "Missing seconds (the problematic format)",
			jsonValue:  `"2025-10-26T14:05Z"`,
			wantError:  false,
			wantResult: "2025-10-26T14:05:00Z",
		},
		{
			name:       "RFC3339 with timezone offset",
			jsonValue:  `"2025-10-26T14:05:00+08:00"`,
			wantError:  false,
			wantResult: "2025-10-26T06:05:00Z",
		},
		{
			name:       "RFC3339Nano with nanoseconds",
			jsonValue:  `"2025-10-26T14:05:00.123456789Z"`,
			wantError:  false,
			wantResult: "2025-10-26T14:05:00.123456789Z",
		},
		{
			name:       "Null value",
			jsonValue:  `null`,
			wantError:  false,
			wantResult: "",
		},
		{
			name:       "Empty string",
			jsonValue:  `""`,
			wantError:  false,
			wantResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ft FlexibleTime
			err := json.Unmarshal([]byte(tt.jsonValue), &ft)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.wantResult == "" {
				// For null/empty values, check if time is zero
				if !ft.IsZero() {
					t.Errorf("Expected zero time for null/empty value, got %v", ft.Time)
				}
			} else {
				// Parse expected result for comparison
				expectedTime, _ := time.Parse(time.RFC3339Nano, tt.wantResult)
				if !ft.Time.Equal(expectedTime) {
					t.Errorf("Expected time %v, got %v", expectedTime, ft.Time)
				}
			}
		})
	}
}

func TestFlexibleTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		time       time.Time
		wantResult string
	}{
		{
			name:       "Normal time",
			time:       time.Date(2025, 10, 26, 14, 5, 0, 0, time.UTC),
			wantResult: `"2025-10-26T14:05:00Z"`,
		},
		{
			name:       "Zero time",
			time:       time.Time{},
			wantResult: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ft := FlexibleTime{Time: tt.time}
			result, err := json.Marshal(ft)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if string(result) != tt.wantResult {
				t.Errorf("Expected %s, got %s", tt.wantResult, string(result))
			}
		})
	}
}

func TestFlexibleTime_ToTimePtr(t *testing.T) {
	t.Run("Non-zero time", func(t *testing.T) {
		now := time.Now()
		ft := FlexibleTime{Time: now}
		ptr := ft.ToTimePtr()

		if ptr == nil {
			t.Error("Expected non-nil pointer")
			return
		}

		if !ptr.Equal(now) {
			t.Errorf("Expected time %v, got %v", now, *ptr)
		}
	})

	t.Run("Zero time", func(t *testing.T) {
		ft := FlexibleTime{Time: time.Time{}}
		ptr := ft.ToTimePtr()

		if ptr != nil {
			t.Error("Expected nil pointer for zero time")
		}
	})

	t.Run("Nil FlexibleTime", func(t *testing.T) {
		var ft *FlexibleTime
		ptr := ft.ToTimePtr()

		if ptr != nil {
			t.Error("Expected nil pointer for nil FlexibleTime")
		}
	})
}

func TestFlexibleTime_InJSON(t *testing.T) {
	// Test that FlexibleTime works correctly when embedded in a struct
	type TestStruct struct {
		Timestamp *FlexibleTime `json:"timestamp,omitempty"`
	}

	tests := []struct {
		name      string
		jsonInput string
		wantError bool
	}{
		{
			name:      "Valid timestamp with missing seconds",
			jsonInput: `{"timestamp":"2025-10-26T14:05Z"}`,
			wantError: false,
		},
		{
			name:      "Valid timestamp with seconds",
			jsonInput: `{"timestamp":"2025-10-26T14:05:00Z"}`,
			wantError: false,
		},
		{
			name:      "Null timestamp",
			jsonInput: `{"timestamp":null}`,
			wantError: false,
		},
		{
			name:      "Missing timestamp field",
			jsonInput: `{}`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts TestStruct
			err := json.Unmarshal([]byte(tt.jsonInput), &ts)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
