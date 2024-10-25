package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestTextWrapping(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Text/Word Wrapping Testing
		{
			name:     "Single-line paragraph (should wrap once)",
			input:    []string{"This is a fairly long paragraph that was input by the user in a single line and will likely need to be wrapped."},
			expected: "This is a fairly long paragraph that was input by the user in a single line and\nwill likely need to be wrapped.",
		},
		{
			name:     "Multi-line paragraph (should wrap once)",
			input:    []string{"This is a fairly long paragraph", "that was input by the user in two separate lines and will likely need to be wrapped."},
			expected: "This is a fairly long paragraph that was input by the user in two separate lines\nand will likely need to be wrapped.",
		},
		{
			name:     "Single-line paragraph (should wrap twice)",
			input:    []string{"This is a fairly long paragraph that was input by the user in a single line and will likely need to be wrapped. It is longer than the last one and as such it should be wrapped twice."},
			expected: "This is a fairly long paragraph that was input by the user in a single line and\nwill likely need to be wrapped. It is longer than the last one and as such it\nshould be wrapped twice.",
		},
		{
			name:     "Multi-line paragraph (should wrap twice)",
			input:    []string{"This is a fairly long paragraph", "that was input by the user in two separate lines and will likely need to be wrapped.", "It is longer than the last one and as such it should be wrapped twice."},
			expected: "This is a fairly long paragraph that was input by the user in two separate lines\nand will likely need to be wrapped. It is longer than the last one and as such\nit should be wrapped twice.",
		},
		{
			name:     "Single-line paragraph w/bold text across lines (should wrap once)",
			input:    []string{"This is a fairly long paragraph that was input by the user in a single **line and will likely need to be** wrapped."},
			expected: "This is a fairly long paragraph that was input by the user in a single \033[1mline and\nwill likely need to be\033[0m wrapped.",
		},
		{
			name:     "Single-line paragraph w/italic text across lines (should wrap once)",
			input:    []string{"This is a fairly long paragraph that was input by the user in a single *line and will likely need to be* wrapped."},
			expected: "This is a fairly long paragraph that was input by the user in a single \033[3mline and\nwill likely need to be\033[0m wrapped.",
		},
		{
			name:     "Single-line paragraph w/strikethrough text across lines (should wrap once)",
			input:    []string{"This is a fairly long paragraph that was input by the user in a single ~~line and will likely need to be~~ wrapped."},
			expected: "This is a fairly long paragraph that was input by the user in a single \033[9mline and\nwill likely need to be\033[0m wrapped.",
		},
		{
			name:     "Single-line paragraph w/inline code across lines (should wrap once)",
			input:    []string{"This is a fairly long paragraph that was input by the user in a single `line and will likely need to be` wrapped."},
			expected: "This is a fairly long paragraph that was input by the user in a single \033[48;5;238;38;5;1mline and\nwill likely need to be\033[0m wrapped.",
		},
		{
			name:     "Single-line paragraph w/o spaces (should wrap once)",
			input:    []string{"Thisisafairlylongparagraphthatwasinputbytheuserinasinglelinewithoutanyspacesandwilllikelyneedtobewrapped."},
			expected: "Thisisafairlylongparagraphthatwasinputbytheuserinasinglelinewithoutanyspacesandw\nilllikelyneedtobewrapped.",
		},
		{
			name:     "Paragraph w/manual break via <br> (should wrap once)",
			input:    []string{"This is a fairly long paragraph that was input by<br>", "the user and should be wrapped uniquely since the manual break occurs before the terminal width."},
			expected: "This is a fairly long paragraph that was input by\nthe user and should be wrapped uniquely since the manual break occurs before the\nterminal width.",
		},
		{
			name:     "Paragraph w/manual break via double space (should wrap once)",
			input:    []string{"This is a fairly long paragraph that was input by  ", "the user and should be wrapped uniquely since the manual break occurs before the terminal width."},
			expected: "This is a fairly long paragraph that was input by\nthe user and should be wrapped uniquely since the manual break occurs before the\nterminal width.",
		},
		{
			name:     "Mixed list w/long items that should wrap",
			input:    []string{"- First item", "\t1. Sub item 1", "\t2. Sub item 2 is really quite long gee I wonder if this will wrap correctly to the next line it sure would be nice if it did", "- Second item is also really long to prove that unindented list items can also wrap in the same way as indented ones"},
			expected: "• First item\n    1. Sub item 1\n    2. Sub item 2\n• Second item\n",
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
