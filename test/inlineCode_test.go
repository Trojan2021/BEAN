package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestInlineCode(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Inline Code Testing
		{
			name:     "Plain in-line code block",
			input:    []string{"`code`"},
			expected: "\033[48;5;238;38;5;1mcode\033[0m",
		},
		{
			name:     "Plain in-line code block (emphasis within code block)",
			input:    []string{"`**_~~code~~_**`"},
			expected: "\033[48;5;238;38;5;1m**_~~code~~_**\033[0m",
		},
		{
			name:     "Used in a paragraph (single-line)",
			input:    []string{"Paragraph with `code` in it."},
			expected: "Paragraph with \033[48;5;238;38;5;1mcode\033[0m in it.",
		},
		{
			name:     "Used in a paragraph (multi-line)",
			input:    []string{"Paragraph with", "`code` in it."},
			expected: "Paragraph with \033[48;5;238;38;5;1mcode\033[0m in it.",
		},
		{
			name:     "Used in a paragraph (multi-line; all forms of emphasis)",
			input:    []string{"**Paragraph** _with_ `code` ~~in it~~ `**and such**`.", "Code blocks can be **_~~`emphasized`~~_** externally."},
			expected: "\033[1mParagraph\033[0m \033[3mwith\033[0m \033[48;5;238;38;5;1mcode\033[0m \033[9min it\033[0m \033[48;5;238;38;5;1m**and such**\033[0m. Code blocks can be \033[1m\033[3m\033[9m\033[48;5;238;38;5;1memphasized\033[0m\033[0m\033[0m\033[0m externally.",
		},
		{
			name:     "Used in a list",
			input:    []string{"- List parent", "\t- List child with `code` in it."},
			expected: "• List parent\n    • List child with \033[48;5;238;38;5;1mcode\033[0m in it.\n",
		},
		{
			name:     "Used in a header",
			input:    []string{"# Header With `Code` In It"},
			expected: "\033[1m─Header With \033[48;5;238;38;5;1mCode\033[0m In It─\033[0m\n",
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
