package list

import (
	"strings"
)

var exclude map[string]bool

func init() {
	exclude = make(map[string]bool)
	for _, word := range strings.Split(excluded, "\n") {
		exclude[word] = true
	}
}

// IsExclude returns true if w is excluded.
func IsExclude(w string) bool {
	_, ok := exclude[w]
	return ok
}
