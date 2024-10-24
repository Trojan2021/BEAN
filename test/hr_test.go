package test

import (
	"os"
	"strings"
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
	"golang.org/x/term"
)

func TestHR(t *testing.T) {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Horizontal Rule Testing
		{
			name:     "Hyphen HR",
			input:    []string{"---"},
			expected: strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Asterisk HR",
			input:    []string{"***"},
			expected: strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Underscore HR",
			input:    []string{"___"},
			expected: strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Mixed HR (should render as paragraph)",
			input:    []string{"_*-"},
			expected: "_*- ",
		},
		{
			name:     "Overkill hyphen HR",
			input:    []string{"----------"},
			expected: strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Overkill asterisk HR",
			input:    []string{"**********"},
			expected: strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Overkill underscore HR",
			input:    []string{"__________"},
			expected: strings.Repeat("─", width) + "\n\n",
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
