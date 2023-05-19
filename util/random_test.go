package util

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestRandomLineID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: "C0da0tzu4cbvjuswj4p4mrj94drmim0n2zxei3in0oahx5pdvyn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomLineID(); got != tt.want {
				t.Errorf("RandomLineID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomPicture(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: "https://loremflickr.com/320/240/hygclfnmkmpvbbzitr586vn9615cty90210txfhh",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomPicture(); got != tt.want {
				t.Errorf("RandomPicture() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomRating(t *testing.T) {
	tests := []struct {
		name string
		want decimal.NullDecimal
	}{
		{
			name: "test1",
			want: decimal.NullDecimal{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomRating(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RandomRating() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomInt32(t *testing.T) {
	type args struct {
		min int32
		max int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "test1",
			args: args{
				min: 50,
				max: 100,
			},
			want: 50,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomInt32(tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("RandomInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomPhoneNumber(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: "13000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomPhoneNumber(); got != tt.want {
				t.Errorf("RandomPhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
