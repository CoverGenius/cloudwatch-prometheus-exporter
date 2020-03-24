package helpers

import (
	"strings"
)

func FormatString(format string, args ...string) string {
	boilerplate := strings.NewReplacer(args...)
	out := boilerplate.Replace(format)
	return out
}

func AWSStrings(strings ...string) []*string {
	aws_strings := make([]*string, len(strings))
	for i, _ := range aws_strings {
		aws_strings[i] = &strings[i]
	}
	return aws_strings
}
