package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var array []*float64

func init() {
	array = floatPointers(2, 1, 4, 3)
}

func floatPointers(values ...float64) []*float64 {
	fp := make([]*float64, len(values))
	for i, _ := range fp {
		fp[i] = &values[i]
	}
	return fp
}

func timePointers(times ...time.Time) []*time.Time {
	tp := make([]*time.Time, len(times))
	for i, _ := range tp {
		tp[i] = &times[i]
	}
	return tp
}

func TestMax(t *testing.T) {
	got, err := Max(array)
	if err != nil {
		t.Errorf("Got err %s", err)
	}
	if got != 4 {
		t.Errorf("Max(%v) = %f; want 4", array, got)
	}
}

func TestMin(t *testing.T) {
	got, err := Min(array)
	if err != nil {
		t.Errorf("Got err %s", err)
	}
	if got != 1 {
		t.Errorf("Min(%v) = %f; want 1", array, got)
	}
}

func TestSum(t *testing.T) {
	got, err := Sum(array)
	if err != nil {
		t.Errorf("Got err %s", err)
	}
	if got != 10 {
		t.Errorf("Sum(%v) = %f; want 10", array, got)
	}
}

func TestAverage(t *testing.T) {
	got, err := Average(array)
	if err != nil {
		t.Errorf("Got err %s", err)
	}
	if got != 2.5 {
		t.Errorf("Average(%v) = %f; want 2.5", array, got)
	}
}

var newValuesTests = []struct {
	values    []*float64
	times     []*time.Time
	threshold time.Time
	expected  []*float64
}{
	{floatPointers(), timePointers(), time.Now(), floatPointers()},
	{floatPointers(1), timePointers(time.Now().Add(-time.Hour)), time.Now(), floatPointers()},
	{floatPointers(1), timePointers(time.Now()), time.Now().Add(-time.Hour), floatPointers(1)},
}

func TestNewValues(t *testing.T) {
	for _, v := range newValuesTests {
		got := NewValues(v.values, v.times, v.threshold)
		assert.Equal(t, v.expected, got)
	}
}
