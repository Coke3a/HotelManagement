package http

import (
	"strconv"
	"time"
	"fmt"
)

// toMap is a helper function to add meta and data to a map
func toMap(m meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}

func convertStringToUint64(str string) (i uint64, err error) {
	i, err = strconv.ParseUint(str, 10, 64)
    if err != nil {
        return 0, err
    }
	return i, nil
}

func zeroValueOrDefault[T comparable](val *T, defaultValue T) T {
	if val != nil {
		return *val
	}
	return defaultValue
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	return t
}

func parseTime(dateStr string) (time.Time, error) {
    // Try parsing as RFC3339 first
    t, err := time.Parse(time.RFC3339, dateStr)
    if err == nil {
        return t, nil
    }

    // If that fails, try parsing as date only
    t, err = time.Parse("2006-01-02", dateStr)
    if err == nil {
        return t, nil
    }

    // If both fail, return an error
    return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}