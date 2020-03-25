package helpers

import "testing"

var array []*float64

func init() {
	one := 1.0
	two := 2.0
	three := 3.0
	four := 4.0
	array = []*float64{&two, &one, &four, &three}
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
