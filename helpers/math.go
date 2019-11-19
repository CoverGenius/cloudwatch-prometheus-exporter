package helpers

func CountAverage(items []*float64, size *float64) (*float64, error) {
	var sum float64 = 0
	for _, item := range items {
		sum += *item
	}
	average := sum / (*size)
	return &average, nil
}
