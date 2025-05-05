package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " another   test case   ",
			expected: []string{"another", "test", "case"},
		},
		{
			input:    "LaSt OnE   ",
			expected: []string{"last", "one"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("%s does not match with %s", word, expectedWord)
			}
		}
	}
}
