package test

import (
	"strings"
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestHR(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Horizontal Rule Testing
		{
			name:     "Hyphen HR",
			input:    []string{"---"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n",
		},
		{
			name:     "Asterisk HR",
			input:    []string{"***"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n",
		},
		{
			name:     "Underscore HR",
			input:    []string{"___"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n",
		},
		{
			name:     "Mixed HR (should render as paragraph)",
			input:    []string{"_*-"},
			expected: "_*-",
		},
		{
			name:     "Overkill hyphen HR",
			input:    []string{"----------"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n",
		},
		{
			name:     "Overkill asterisk HR",
			input:    []string{"**********"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n",
		},
		{
			name:     "Overkill underscore HR",
			input:    []string{"__________"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := bean.RenderMarkdown(tt.input, terminalWidth)
			if output != tt.expected {
				bufferFailure(t, output, tt.expected)
			}
		})
	}
}
