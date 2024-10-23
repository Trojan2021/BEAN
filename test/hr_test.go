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
			name:     "Overkill Hyphen HR",
			input:    []string{"----------"},
			expected: strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Overkill Asterisk HR",
			input:    []string{"**********"},
			expected: strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Overkill Underscore HR",
			input:    []string{"__________"},
			expected: strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Multiple HRs",
			input:    []string{"___", "---", "***", "***"},
			expected: strings.Repeat("─", width) + "\n\n" + strings.Repeat("─", width) + "\n\n" + strings.Repeat("─", width) + "\n\n" + strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "HR After Paragraph",
			input:    []string{"This is a paragraph.", "---"},
			expected: "This is a paragraph. \n\n" + strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Paragraph After HR",
			input:    []string{"---", "This is a paragraph."},
			expected: strings.Repeat("─", width) + "\n\n" + "This is a paragraph. ",
		},
		{
			name:     "Paragraph After HR (w/blank lines)",
			input:    []string{"---", "", "", "This is a paragraph."},
			expected: strings.Repeat("─", width) + "\n\n" + "This is a paragraph. ",
		},
		{
			name:     "HR After Header",
			input:    []string{"# Heading 1", "---"},
			expected: "\033[1m─Heading 1─\033[0m\n\n" + strings.Repeat("─", width) + "\n\n",
		},
		{
			name:     "Header After HR",
			input:    []string{"---", "# Heading 1"},
			expected: strings.Repeat("─", width) + "\n\n" + "\033[1m─Heading 1─\033[0m\n",
		},
		{
			name:     "Header After HR (w/blank lines)",
			input:    []string{"---", "", "# Heading 1"},
			expected: strings.Repeat("─", width) + "\n\n" + "\033[1m─Heading 1─\033[0m\n",
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
