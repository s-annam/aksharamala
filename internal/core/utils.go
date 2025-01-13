package core

import "strings"

func NormalizeComment(comment string) string {
	comment = strings.TrimSpace(comment)
	comment = strings.TrimPrefix(comment, "=*=")
	comment = strings.TrimSuffix(comment, "=*=")
	return strings.TrimSpace(comment)
}
