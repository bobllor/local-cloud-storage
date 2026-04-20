package sqlquery

import (
	"fmt"
	"strings"
)

// BuildPlaceholder builds the placeholder string for parameters that are passed into
// queries.
// It will return the string base on the placeholder and repeat counts: "(?,?,...),(...)..."
// This does not include VALUES.
//
// placholderCount is the amount of placeholders that are being used. This is expected
// to be the number of parameters being used.
//
// repeat is how many times to repeat the final placeholder string. This is used for batch operations.
// If batch operations are not needed then it is expected to be 1. If less than 1 is given,
// that it will automatically be converted into 1.
func BuildPlaceholder(placeholderCount int, repeat int) string {
	questions := []string{}
	out := []string{}

	if repeat < 1 {
		repeat = 1
	}

	for range placeholderCount {
		questions = append(questions, "?")
	}

	param := "(" + strings.Join(questions, ",") + ")"

	for range repeat {
		out = append(out, param)
	}

	return strings.Join(out, ",")
}

// BuildSetPlaceholder builds the strings for updating columns in a table.
// The output will be in the form of: "Column1 = value,Column2 = value, ..."
func BuildSetPlaceholder(columns []string) string {
	placeholders := []string{}

	for _, column := range columns {
		placeholder := fmt.Sprintf("%s = ?", column)

		placeholders = append(placeholders, placeholder)
	}

	return strings.Join(placeholders, ",")
}
