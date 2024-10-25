package test

import (
	"strings"
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestNewline(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// New line testing (relationship between headers/HRs/lists/paragraphs)
		{
			name:     "Minimal spacing + <br> line break + standard middle paragraph",
			input:    []string{"- First item", "\t- Sub item 1", "\t- Sub item 2", "- Second item", "# Heading 1", "Paragraph text", "that spans two lines.", "", "New paragraph.", "# Header 1", "__*Fancy*__ paragraph<br>", "with a ~~line break~~."},
			expected: "• First item\n    • Sub item 1\n    • Sub item 2\n• Second item\n\n\033[1m─Heading 1─\033[0m\nParagraph text that spans two lines.\n\nNew paragraph.\n\n\033[1m─Header 1─\033[0m\n\033[1m\033[3mFancy\033[0m\033[0m paragraph\nwith a \033[9mline break\033[0m.",
		},
		{
			name:     "Standard spacing + double-space line break + emphasized middle paragraph",
			input:    []string{"- First item", "\t- Sub item 1", "\t- Sub item 2", "- Second item", "", "# Heading 1", "Paragraph text", "that spans two lines.", "", "*New* paragraph.", "", "# Header 1", "__*Fancy*__ paragraph  ", "with a ~~line break~~."},
			expected: "• First item\n    • Sub item 1\n    • Sub item 2\n• Second item\n\n\033[1m─Heading 1─\033[0m\nParagraph text that spans two lines.\n\n\033[3mNew\033[0m paragraph.\n\n\033[1m─Header 1─\033[0m\n\033[1m\033[3mFancy\033[0m\033[0m paragraph\nwith a \033[9mline break\033[0m.",
		},
		{
			name:     "Exaggerated spacing + double-space line break + standard middle paragraph",
			input:    []string{"- First item", "\t- Sub item 1", "\t- Sub item 2", "- Second item", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "# Heading 1", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "Paragraph text", "that spans two lines.", "", "New paragraph.", "", "", "", "", "", "", "", "# Header 1", "__*Fancy*__ paragraph  ", "with a ~~line break~~."},
			expected: "• First item\n    • Sub item 1\n    • Sub item 2\n• Second item\n\n\033[1m─Heading 1─\033[0m\nParagraph text that spans two lines.\n\nNew paragraph.\n\n\033[1m─Header 1─\033[0m\n\033[1m\033[3mFancy\033[0m\033[0m paragraph\nwith a \033[9mline break\033[0m.",
		},
		{
			name:     "Paragraph followed by list (minimal spacing)",
			input:    []string{"This is a happy little paragraph", "- First item", "- Second item", "\t- Sub item 1", "- Third item"},
			expected: "This is a happy little paragraph\n• First item\n• Second item\n    • Sub item 1\n• Third item\n",
		},
		{
			name:     "Paragraph followed by list (standard spacing; the empty line should be preserved in this case)",
			input:    []string{"This is a happy little paragraph", "", "- First item", "- Second item", "\t- Sub item 1", "- Third item"},
			expected: "This is a happy little paragraph\n\n• First item\n• Second item\n    • Sub item 1\n• Third item\n",
		},
		{
			name:     "Paragraph followed by list (exaggerated spacing; one empty line should be preserved in this case)",
			input:    []string{"This is a happy little paragraph", "", "", "", "", "", "", "", "", "- First item", "- Second item", "\t- Sub item 1", "- Third item"},
			expected: "This is a happy little paragraph\n\n• First item\n• Second item\n    • Sub item 1\n• Third item\n",
		},
		{
			name:     "List followed by a paragraph (minimal spacing)",
			input:    []string{"- First item", "- Second item", "\t- Sub item 1", "- Third item", "This is a happy little paragraph"},
			expected: "• First item\n• Second item\n    • Sub item 1\n• Third item\n\nThis is a happy little paragraph",
		},
		{
			name:     "List followed by a paragraph (standard spacing)",
			input:    []string{"- First item", "- Second item", "\t- Sub item 1", "- Third item", "", "This is a happy little paragraph"},
			expected: "• First item\n• Second item\n    • Sub item 1\n• Third item\n\nThis is a happy little paragraph",
		},
		{
			name:     "List followed by a paragraph (exaggerated spacing)",
			input:    []string{"- First item", "- Second item", "\t- Sub item 1", "- Third item", "", "", "", "", "", "", "", "", "This is a happy little paragraph"},
			expected: "• First item\n• Second item\n    • Sub item 1\n• Third item\n\nThis is a happy little paragraph",
		},
		{
			name:     "Multiple HRs",
			input:    []string{"___", "---", "***", "***", "", "", "", "___"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n" + strings.Repeat("─", terminalWidth) + "\n\n" + strings.Repeat("─", terminalWidth) + "\n\n" + strings.Repeat("─", terminalWidth) + "\n\n" + strings.Repeat("─", terminalWidth) + "\n\n",
		},
		{
			name:     "HR after paragraph",
			input:    []string{"This is a paragraph.", "---"},
			expected: "This is a paragraph.\n\n" + strings.Repeat("─", terminalWidth) + "\n\n",
		},
		{
			name:     "Paragraph after HR",
			input:    []string{"---", "This is a paragraph."},
			expected: strings.Repeat("─", terminalWidth) + "\n\n" + "This is a paragraph.",
		},
		{
			name:     "Paragraph after HR (w/blank lines)",
			input:    []string{"---", "", "", "This is a paragraph."},
			expected: strings.Repeat("─", terminalWidth) + "\n\n" + "This is a paragraph.",
		},
		{
			name:     "HR after header",
			input:    []string{"# Heading 1", "---"},
			expected: "\033[1m─Heading 1─\033[0m\n\n" + strings.Repeat("─", terminalWidth) + "\n\n",
		},
		{
			name:     "Header after HR",
			input:    []string{"---", "# Heading 1"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n" + "\033[1m─Heading 1─\033[0m\n",
		},
		{
			name:     "Header after HR (w/blank lines)",
			input:    []string{"---", "", "# Heading 1"},
			expected: strings.Repeat("─", terminalWidth) + "\n\n" + "\033[1m─Heading 1─\033[0m\n",
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
