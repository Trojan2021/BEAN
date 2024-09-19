package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Header Testing
		{
			name:     "Header level 1",
			input:    []string{"# Heading 1"},
			expected: "\033[1m─Heading 1─\033[0m\n",
		},
		{
			name:     "Header level 2",
			input:    []string{"## Heading 2"},
			expected: "\033[1m──Heading 2──\033[0m\n",
		},
		{
			name:     "Header level 3",
			input:    []string{"### Heading 2"},
			expected: "\033[1m───Heading 2───\033[0m\n",
		},
		{
			name:     "Header level 4",
			input:    []string{"#### Heading 2"},
			expected: "\033[1m────Heading 2────\033[0m\n",
		},
		{
			name:     "Header level 5",
			input:    []string{"##### Heading 2"},
			expected: "\033[1m─────Heading 2─────\033[0m\n",
		},
		{
			name:     "Header level 6",
			input:    []string{"###### Heading 2"},
			expected: "\033[1m──────Heading 2──────\033[0m\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := bean.RenderMarkdown(tt.input)
			if output != tt.expected {
				bufferFailure(t, output, tt.expected)
			}
		})
	}
}
