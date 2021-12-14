package helpers

import (
	"errors"
	"time"
)

// Average returns the mean of a slice of float64 pointers
//
// If the input slice is empty returns 0
func Average(items []*float64) (float64, error) {
	var sum float64
	for _, item := range items {
		sum += *item
	}
	average := sum / float64(len(items))
	return average, nil
}

// Sum returns the sum of a slice of float64 pointers
func Sum(items []*float64) (float64, error) {
	var sum float64
	for _, item := range items {
		sum += *item
	}
	return sum, nil
}

// Min returns the smallest value from a slice of float64 pointers
//
// returns an error if the input slice is empty
func Min(items []*float64) (float64, error) {
	if len(items) < 1 {
		return 0.0, errors.New("cannot calculate minimum of empty list")
	}

	var min = *items[0]
	for _, item := range items {
		if *item < min {
			min = *item
		}
	}
	return min, nil
}

// Max returns the largest value from a slice of float64 pointers
func Max(items []*float64) (float64, error) {
	var max float64
	for _, item := range items {
		if *item > max {
			max = *item
		}
	}
	return max, nil
}

// NewValues filters the slice of values to remove any which are not newer than the input threshold.
//
// The timestamp for value[x] is taken to be times[x]
func NewValues(values []*float64, times []*time.Time, threshold time.Time) []*float64 {
	newValues := []*float64{}
	for i, value := range values {
		if times[i].After(threshold) {
			newValues = append(newValues, value)
		}
	}
	return newValues
}

func IntPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}
