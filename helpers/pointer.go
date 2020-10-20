package helpers

import (
	"time"
)

func BoolP(b bool) *bool {
	return &b
}

func StringP(s string) *string {
	return &s
}

func IntP(i int) *int {
	return &i
}

func TimeP(t time.Time) *time.Time {
	return &t
}

func FloatP(f float64) *float64 {
	return &f
}
