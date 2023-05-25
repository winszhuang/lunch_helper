package util

import (
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
