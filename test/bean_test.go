package bean

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/src"
)

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Header level 1",
			input:    []string{"# Heading 1"},
			expected: "\033[1m\033[4mHeading 1\033[0m\n",
		},
		{
			name:     "Header level 2",
			input:    []string{"## Heading 2"},
			expected: "\033[1mHeading 2\033[0m\n",
		},
		{
			name:     "Unordered list",
			input:    []string{"- List item"},
			expected: "• List item",
		},
		{
			name:     "Plain text",
			input:    []string{"Just a line of text."},
			expected: "Just a line of text. ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := bean.RenderMarkdown(tt.input)
			if output != tt.expected {
				t.Errorf("got %q, want %q", output, tt.expected)
			}
		})
	}
}
