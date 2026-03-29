package utils

// ConvertToAny converts a slice to an []any. This is only used for
// for query arguments.
func ConvertToAny[S ~[]T, T comparable](v S) []any {
	conv := []any{}

	for _, e := range v {
		conv = append(conv, e)
	}

	return conv
}
