package dckr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContainerName(t *testing.T) {
	testCases := []struct {
		name     string
		names    []string
		expected string
	}{
		{
			name:     "Single name",
			names:    []string{"/container1"},
			expected: "container1",
		},
		{
			name:     "Multiple names",
			names:    []string{"/container1", "/alias1", "/alias2"},
			expected: "container1",
		},
		{
			name:     "Empty names",
			names:    []string{},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetContainerName(tc.names)
			assert.Equal(t, tc.expected, result, "GetContainerName should return the expected name")
		})
	}
}

// Add more tests for other existing functions in utils.go here
