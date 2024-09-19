package http

import (
	"strconv"
	"time"
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