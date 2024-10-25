package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestRenderStrikethroughText(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Basic strikethrough text
		{
			name:     "Strikethrough text using double tildes",
			input:    []string{"~~strikethrough~~"},
			expected: "\033[9mstrikethrough\033[0m",
		},
		// Strikethrough within a sentence
		{
			name:     "Strikethrough in the middle of sentence",
			input:    []string{"A sentence with ~~strikethrough~~ in the middle."},
			expected: "A sentence with \033[9mstrikethrough\033[0m in the middle.",
		},
		{
			name:     "Strikethrough at the start and end of sentence",
			input:    []string{"~~Strikethrough~~ at the start, and at the ~~end~~"},
			expected: "\033[9mStrikethrough\033[0m at the start, and at the \033[9mend\033[0m",
		},
		{
			name:     "Multiple strikethrough segments in one line (not at start or end)",
			input:    []string{"This is ~~strikethrough1~~, and this is ~~strikethrough2~~ as well."},
			expected: "This is \033[9mstrikethrough1\033[0m, and this is \033[9mstrikethrough2\033[0m as well.",
		},
		// Edge cases
		{
			name:     "Strikethrough with incomplete tildes (should render paragraph)",
			input:    []string{"This is ~not fully strikethrough~~ text."},
			expected: "This is ~not fully strikethrough~~ text.",
		},
		// Nested strikethrough with other formatting (italic/bold)
		{
			name:     "Strikethrough nested inside italic",
			input:    []string{"This is _italic ~~and strikethrough~~ inside_."},
			expected: "This is \033[3mitalic \033[9mand strikethrough\033[0m inside\033[0m.",
		},
		{
			name:     "Strikethrough nested inside bold",
			input:    []string{"This is __bold ~~and strikethrough~~ inside__."},
			expected: "This is \033[1mbold \033[9mand strikethrough\033[0m inside\033[0m.",
		},
		// Strikethrough with special characters
		{
			name:     "Strikethrough with special characters",
			input:    []string{"This is ~~strikethrough & special <chars>~~ text."},
			expected: "This is \033[9mstrikethrough & special <chars>\033[0m text.",
		},
		// Strikethrough tilde
		{
			name:     "Strikethrough tilde",
			input:    []string{"Strikethrough tildes (~~~~~) should be possible."},
			expected: "Strikethrough tildes (\033[9m~\033[0m) should be possible.",
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
