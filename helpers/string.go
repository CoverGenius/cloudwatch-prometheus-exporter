package helpers

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

// StringPointers converts a slice of string values into a slice of string pointers
//
// This function complements aws.StringSlice but works with variadic arguments so that an array literal is not required.
func StringPointers(s ...string) []*string {
	sp := make([]*string, len(s))
	for i := range sp {
		sp[i] = &s[i]
	}
	return sp
}

var matchUnsafeChars = regexp.MustCompile("[^a-zA-Z0-9:_]")
var badSplits = map[string]string{
	"un_healthy":   "unhealthy",
	"_lc_us":       "_lcus",
	"ec_2_":        "ec2_",
	"s_3_":         "s3_",
	"elasti_cache": "elasticache",
	"_i_ds":        "_ids",
}

// ToSnakeCase converts a CamelCaseString to snake_case
func ToSnakeCase(s string) string {
	s = strcase.ToSnake(s)
	for oldS, newS := range badSplits {
		s = strings.ReplaceAll(s, oldS, newS)
	}
	return s
}

func safeName(s string) string {
	return matchUnsafeChars.ReplaceAllString(s, "_")
}

// ToPromString process a string to ensure it is compatible with prometheus
//
// It converts the string to lower_snake_case, removes non-alphanumeric characters
// and applies a few special case rules.
func ToPromString(s string) string {
	s = ToSnakeCase(s)
	s = strings.ReplaceAll(s, "application_elb", "alb")
	s = strings.ReplaceAll(s, "network_elb", "nlb")
	s = safeName(s)
	return regexp.MustCompile("__+").ReplaceAllString(s, "_")
}
