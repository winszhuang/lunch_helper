package util

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRegexQuery(t *testing.T) {
	regex := regexp.MustCompile(`lat=([-+]?[\d.]+),lng=([-+]?[\d.]+),radius=(\d+),pageIndex=(\d+)`)

	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "valid input",
			input:    "lat=24.123,lng=120.456,radius=500,pageIndex=1",
			expected: []string{"24.123", "120.456", "500", "1"},
		},
		{
			name:     "invalid input",
			input:    "invalid input",
			expected: []string{},
		},
		{
			name:     "missing values",
			input:    "lat=,lng=",
			expected: []string{},
		},
		{
			name:     "negative values",
			input:    "lat=-24.123,lng=-120.456,radius=500,pageIndex=1",
			expected: []string{"-24.123", "-120.456", "500", "1"},
		},
		{
			name:     "floating-point values",
			input:    "lat=24.123456,lng=120.654321,radius=500,pageIndex=1",
			expected: []string{"24.123456", "120.654321", "500", "1"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := ParseRegexQuery(testCase.input, regex)
			require.Equal(t, testCase.expected, result)
		})
	}
}

func TestParseQuery(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				input: "/searchnext?lat=24.17829&lng=120.68012&radius=500&pageIndex=2",
			},
			want:    []string{"24.17829", "120.68012", "500", "2"},
			wantErr: false,
		},
		{
			name: "miss some args",
			args: args{
				input: "/searchnext?lat&lng=120.68012&radius=500&pageIndex=2",
			},
			want:    []string{"", "120.68012", "500", "2"},
			wantErr: false,
		},
		{
			name: "miss arg",
			args: args{
				input: "/userrestaurantnext?pageIndex=",
			},
			want:    []string{""},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQuery(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
