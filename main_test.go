package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "L O L",
			expected: []string{"l", "o", "l"},
		},
		{
			input: "wwvvv a VAa ",
			expected: []string{"wwvvv", "a", "vaa"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check if lengths match
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d",
				c.input, len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			if word != expectedWord {
				t.Errorf("cleanInput(%q) word %d: got %q, expected %q",
					c.input, i, word, expectedWord)
			}
		}
	}
}

