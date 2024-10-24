package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestRenderItalicText(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Basic italic text
		{
			name:     "Italic text using single asterisks",
			input:    []string{"*italic*"},
			expected: "\033[3mitalic\033[0m",
		},
		{
			name:     "Italic text using single underscores",
			input:    []string{"_italic_"},
			expected: "\033[3mitalic\033[0m",
		},
		// Italic within a sentence
		{
			name:     "Italic in the middle of sentence",
			input:    []string{"A sentence with *italic* in the middle."},
			expected: "A sentence with \033[3mitalic\033[0m in the middle.",
		},
		{
			name:     "Italic at the start and end of sentence",
			input:    []string{"*Italic* at the start, and at the *end*"},
			expected: "\033[3mItalic\033[0m at the start, and at the \033[3mend\033[0m",
		},
		{
			name:     "Multiple italic segments in one line (not at start or end)",
			input:    []string{"This is *italic1*, and this is *italic2* as well."},
			expected: "This is \033[3mitalic1\033[0m, and this is \033[3mitalic2\033[0m as well.",
		},
		// Nested italic with other formatting (bold/strikethrough)
		{
			name:     "Italic nested inside bold",
			input:    []string{"This is __bold *and italic* inside__."},
			expected: "This is \033[1mbold \033[3mand italic\033[0m inside\033[0m.",
		},
		{
			name:     "Italic nested inside strikethrough",
			input:    []string{"This is ~~strikethrough *and italic* inside~~."},
			expected: "This is \033[9mstrikethrough \033[3mand italic\033[0m inside\033[0m.",
		},
		// Italic with special characters
		{
			name:     "Italic with special characters",
			input:    []string{"This is *italic & special <chars>* text."},
			expected: "This is \033[3mitalic & special <chars>\033[0m text.",
		},
		// Italic asterisks and underscores
		{
			name:     "Italic asterisks and underscores in one line",
			input:    []string{"Italic asterisks (***) and italic underscores (___) should be possible."},
			expected: "Italic asterisks (\033[3m*\033[0m) and italic underscores (\033[3m_\033[0m) should be possible.",
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
