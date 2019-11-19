package helpers

import (
	"strings"
)

func FormatString(format string, args ...string) string {
	boilerplate := strings.NewReplacer(args...)
	out := boilerplate.Replace(format)
	return out
}
