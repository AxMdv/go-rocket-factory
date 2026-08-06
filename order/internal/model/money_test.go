package model

import "testing"

func TestPriceToCents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		price float64
		want  int64
	}{
		{
			name:  "whole units",
			price: 10,
			want:  1000,
		},
		{
			name:  "fractional units",
			price: 123.45,
			want:  12345,
		},
		{
			name:  "rounds to nearest cent",
			price: 10.129,
			want:  1013,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := PriceToCents(tt.price); got != tt.want {
				t.Fatalf("PriceToCents() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCentsToPrice(t *testing.T) {
	t.Parallel()

	if got := CentsToPrice(12345); got != 123.45 {
		t.Fatalf("CentsToPrice() = %f, want %f", got, 123.45)
	}
}
