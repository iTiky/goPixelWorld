package pkg

// FilterSlice filters the slice by criteria (if false, do not include in the result).
func FilterSlice[T any](values []T, filterPositive func(T) bool) []T {
	n := 0
	for _, value := range values {
		if filterPositive(value) {
			values[n] = value
			n++
		}
	}
	values = values[:n]

	return values
}

// SliceHasValue checks if the slice has the value.
func SliceHasValue[T comparable](values []T, expectedValue T) bool {
	for _, value := range values {
		if value == expectedValue {
			return true
		}
	}

	return false
}

// ValuePtr return the value pointer.
func ValuePtr[T any](v T) *T {
	return &v
}
