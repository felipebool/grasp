package main

import (
	"reflect"
	"testing"
)

func Test_normalizeAddresses(t *testing.T) {
	cases := map[string]struct{
		input []string
		expected []string
	}{
		"empty list": {
			input: []string{},
			expected: []string{},
		},
		"contain invalid address": {
			input: []string{
				"www.google.com",
			},
			expected: []string{
				"http://www.google.com",
			},
		},
		"does not contain invalid address": {
			input: []string{
				"http://www.google.com",
				"https://www.google.com",
			},
			expected: []string{
				"http://www.google.com",
				"https://www.google.com",
			},
		},
	}

	for label := range cases {
		tc := cases[label]
		t.Run(label, func(t *testing.T) {
			t.Parallel()
			got := normalizeAddresses(tc.input)
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("got: %v; expected: %v", got, tc.expected)
			}
		})
	}
}
