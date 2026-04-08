package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// FlexibleTime is a custom type that can parse multiple time formats
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface
// It tries multiple time formats to handle API inconsistencies
func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	s := strings.Trim(string(data), "\"")

	// If it's null or empty, return nil
	if s == "null" || s == "" {
		return nil
	}

	// List of time formats to try, in order of preference
	formats := []string{
		time.RFC3339,             // "2006-01-02T15:04:05Z07:00" (standard)
		"2006-01-02T15:04Z07:00", // "2025-10-26T14:05Z" (missing seconds)
		"2006-01-02T15:04Z",      // "2025-10-26T14:05Z" (missing seconds, UTC)
		"2006-01-02T15:04:05Z",   // "2025-10-26T14:05:00Z" (UTC)
		time.RFC3339Nano,         // "2006-01-02T15:04:05.999999999Z07:00" (with nanoseconds)
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			ft.Time = t
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("unable to parse time %q: %w", s, lastErr)
}

// MarshalJSON implements the json.Marshaler interface
// It always outputs in RFC3339 format for consistency
func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ft.Time.Format(time.RFC3339))
}

// ToTimePtr converts FlexibleTime to *time.Time
func (ft *FlexibleTime) ToTimePtr() *time.Time {
	if ft == nil || ft.IsZero() {
		return nil
	}
	t := ft.Time
	return &t
}
