package core

import "strings"

// NormalizeComment normalizes a comment string by trimming whitespace,
// removing specific prefixes and suffixes, and returning the cleaned-up string.
// It takes a comment as input and returns the normalized comment.
func NormalizeComment(comment string) string {
	comment = strings.TrimSpace(comment)
	comment = strings.TrimPrefix(comment, "=*=")
	comment = strings.TrimSuffix(comment, "=*=")
	return strings.TrimSpace(comment)
}
