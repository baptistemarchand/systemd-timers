package main

import (
	"testing"
)

func TestFormatExecutionTime(t *testing.T) {
	testcases := []struct {
		executionTime uint64
		expected      string
	}{
		{
			executionTime: 0,
			expected:      "",
		},
		{
			executionTime: 100,
			expected:      "0s",
		},
		{
			executionTime: 1 * 1000 * 1000,
			expected:      "1s",
		},
		{
			executionTime: 20 * 1000 * 1000,
			expected:      "20s",
		},
		{
			executionTime: 60 * 1000 * 1000,
			expected:      "<fg 1>1m 0s<reset>",
		},
		{
			executionTime: 61 * 1000 * 1000,
			expected:      "<fg 1>1m 1s<reset>",
		},
		{
			executionTime: 62 * 1000 * 1000,
			expected:      "<fg 1>1m 2s<reset>",
		},
		{
			executionTime: 120 * 1000 * 1000,
			expected:     "<fg 1>2m 0s<reset>",
		},
		{
			executionTime: 121 * 1000 * 1000,
			expected:      "<fg 1>2m 1s<reset>",
		},
		{
			executionTime: 142 * 1000 * 1000,
			expected:      "<fg 1>2m 22s<reset>",
		},
		{
			executionTime: 192 * 60 * 1000 * 1000,
			expected:      "<fg 1>192m 0s<reset>",
		},
	}

	for _, testcase := range testcases {
		got := formatExecutionTime(testcase.executionTime)
		if got != testcase.expected {
			t.Errorf("got: %q, expected: %q", got, testcase.expected)
		}
	}
}
