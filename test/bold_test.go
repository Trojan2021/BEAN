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
			input:    []string{"**bold**"},
			expected: "\033[1mbold\033[0m",
		},
		{
			name:     "Bold text using double underscores",
			input:    []string{"__bold__"},
			expected: "\033[1mbold\033[0m",
		},
		// Bold within a sentence
		{
			name:     "Bold in the middle of sentence",
			input:    []string{"A sentence with **bold** in the middle."},
			expected: "A sentence with \033[1mbold\033[0m in the middle.",
		},
		{
			name:     "Bold at the start and end of sentence",
			input:    []string{"**Bold** at the start, and at the **end**"},
			expected: "\033[1mBold\033[0m at the start, and at the \033[1mend\033[0m",
		},
		{
			name:     "Multiple bold segments in one line (not at start or end)",
			input:    []string{"This is **bold1**, and this is **bold2** as well."},
			expected: "This is \033[1mbold1\033[0m, and this is \033[1mbold2\033[0m as well.",
		},
		// Edge cases
		{
			name:     "Bold with incomplete asterisks (should render italics)",
			input:    []string{"This is *not fully bold** text."},
			expected: "This is \033[3mnot fully bold\033[0m* text.",
		},
		{
			name:     "Bold with incomplete underscores (should render italics)",
			input:    []string{"This is _not fully bold__ text."},
			expected: "This is \033[3mnot fully bold\033[0m_ text.",
		},
		// Nested bold with other formatting (italic/strikethrough)
		{
			name:     "Bold nested inside italic",
			input:    []string{"This is _italic **and bold** inside_."},
			expected: "This is \033[3mitalic \033[1mand bold\033[0m inside\033[0m.",
		},
		{
			name:     "Bold nested inside strikethrough",
			input:    []string{"This is ~~strikethrough **and bold** inside~~."},
			expected: "This is \033[9mstrikethrough \033[1mand bold\033[0m inside\033[0m.",
		},
		// Bold with special characters
		{
			name:     "Bold with special characters",
			input:    []string{"This is **bold & special <chars>** text."},
			expected: "This is \033[1mbold & special <chars>\033[0m text.",
		},
		// Bold asterisks and underscores
		{
			name:     "Bold asterisks and underscores in one line",
			input:    []string{"Bold asterisks (*****) and bold underscores (_____) should be possible."},
			expected: "Bold asterisks (\033[1m*\033[0m) and bold underscores (\033[1m_\033[0m) should be possible.",
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
