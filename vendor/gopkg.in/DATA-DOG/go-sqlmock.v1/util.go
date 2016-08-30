package sqlmock

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile("\\s+")

// strip out new lines and trim spaces
func stripQuery(q string) (s string) {
	return strings.TrimSpace(re.ReplaceAllString(q, " "))
}
