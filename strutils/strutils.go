package strutils

import "strings"

// Concat into a single string and return
func Concat(n int, s ...string) string {
	var b strings.Builder
	b.Grow(n)
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}
