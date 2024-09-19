package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestRenderBoldText(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Basic bold text
		{
			name:     "Bold text using double asterisks",
			input:    []string{"This is **bold** text."},
			expected: "This is \033[1mbold\033[0m text. ",
		},
		{
			name:     "Bold text using double underscores",
			input:    []string{"This is __bold__ text."},
			expected: "This is \033[1mbold\033[0m text. ",
		},
		// Bold within a sentence
		{
			name:     "Bold in the middle of sentence",
			input:    []string{"A sentence with **bold** in the middle."},
			expected: "A sentence with \033[1mbold\033[0m in the middle. ",
		},
		// Bold at the start and end of a sentence
		{
			name:     "Bold at the start and end of sentence",
			input:    []string{"**Bold** at the start, and at the **end**."},
			expected: "\033[1mBold\033[0m at the start, and at the \033[1mend\033[0m. ",
		},
		// Edge cases
		{
			name:     "Bold with incomplete asterisks (should render italics)",
			input:    []string{"This is *not fully bold** text."},
			expected: "This is \033[3mnot fully bold\033[0m* text. ", // Should render as italics
		},
		{
			name:     "Bold with incomplete underscores (should render italics)",
			input:    []string{"This is _not fully bold__ text."},
			expected: "This is \033[3mnot fully bold\033[0m_ text. ", // Should render as italics
		},
		// Nested bold with other formatting (italic)
		{
			name:     "Bold nested inside italic",
			input:    []string{"This is _italic **and bold** inside_."},
			expected: "This is \033[3mitalic \033[1mand bold\033[0m inside\033[0m. ",
		},
		{
			name:     "Italic nested inside bold",
			input:    []string{"This is **bold _and italic_ inside**."},
			expected: "This is \033[1mbold \033[3mand italic\033[0m inside\033[0m. ",
		},
		// Bold with special characters
		{
			name:     "Bold with special characters",
			input:    []string{"This is **bold & special <chars>** text."},
			expected: "This is \033[1mbold & special <chars>\033[0m text. ",
		},
		// Multiple bold segments in one line
		{
			name:     "Multiple bold segments in one line",
			input:    []string{"This is **bold1**, and this is **bold2**."},
			expected: "This is \033[1mbold1\033[0m, and this is \033[1mbold2\033[0m. ",
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
