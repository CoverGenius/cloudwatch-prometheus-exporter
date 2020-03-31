package helpers

import (
	"errors"
	"time"
)

func Average(items []*float64) (float64, error) {
	var sum float64 = 0
	for _, item := range items {
		sum += *item
	}
	average := sum / float64(len(items))
	return average, nil
}

func Sum(items []*float64) (float64, error) {
	var sum float64 = 0
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
		return 0.0, errors.New("Cannot calculate minimum of empty list")
	}

	var min float64 = *items[0]
	for _, item := range items {
		if *item < min {
			min = *item
		}
	}
	return min, nil
}

// Max returns the largest value from a slice of float64 pointers
func Max(items []*float64) (float64, error) {
	var max float64 = 0
	for _, item := range items {
		if *item > max {
			max = *item
		}
	}
	return max, nil
}

// NewValues filters the slice of values to remove any which are not newer than the input timestamp.
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
