package http

import "strconv"

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
