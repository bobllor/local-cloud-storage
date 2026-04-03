package utils

import "time"

// ConvertToAny converts a slice to an []any. This is only used for
// for query arguments.
func ConvertToAny[S ~[]T, T comparable](v S) []any {
	conv := []any{}

	for _, e := range v {
		conv = append(conv, e)
	}

	return conv
}

// FormatTime formats the given time to the format YYYY-MM-DD HH:MM:SS
// as a string.
func FormatTime(date time.Time) string {
	return date.Format(time.DateTime)
}
