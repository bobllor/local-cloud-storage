package utils

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

// GetFiles retrieves all files in a given path.
// It will return a map with the key being the lowercased file name
// and the value being the absolute path of the file.
//
// This does not recursively search through directories.
func GetFiles(root string) (map[string]string, error) {
	files := make(map[string]string)

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, ent := range entries {
		files[strings.ToLower(ent.Name())] = filepath.Join(root, ent.Name())
	}

	return files, nil
}
