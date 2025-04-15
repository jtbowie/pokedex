package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "      hello  world     ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "    fizz    buzz fizzbub   ",
			expected: []string{"fizz", "buzz", "fizzbub"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("You failed the test, dummy!")
			}
		}
	}
}
