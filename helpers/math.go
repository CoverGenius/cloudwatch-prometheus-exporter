package helpers

import "errors"

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

func Max(items []*float64) (float64, error) {
	var max float64 = 0
	for _, item := range items {
		if *item > max {
			max = *item
		}
	}
	return max, nil
}
