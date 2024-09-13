package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestNewline(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// New line testing (relationship between headers/list/paragraphs)
		{
			name:     "Minimal spacing + <br> line break + standard middle paragraph",
			input:    []string{"- First item", "\t- Sub item 1", "\t- Sub item 2", "- Second item", "# Heading 1", "Paragraph text", "that spans two lines.", "", "New paragraph.", "# Header 1", "__*Fancy*__ paragraph<br>", "with a ~~line break~~."},
			expected: "• First item\n    • Sub item 1\n    • Sub item 2\n• Second item\n\n\033[1m\033[4mHeading 1\033[0m\nParagraph text that spans two lines. \n\nNew paragraph. \n\n\033[1m\033[4mHeader 1\033[0m\n\033[1m\033[3mFancy\033[0m\033[0m paragraph\nwith a \033[9mline break\033[0m. ",
		},
		{
			name:     "Standard spacing + double-space line break + emphasized middle paragraph",
			input:    []string{"- First item", "\t- Sub item 1", "\t- Sub item 2", "- Second item", "", "# Heading 1", "Paragraph text", "that spans two lines.", "", "*New* paragraph.", "", "# Header 1", "__*Fancy*__ paragraph  ", "with a ~~line break~~."},
			expected: "• First item\n    • Sub item 1\n    • Sub item 2\n• Second item\n\n\033[1m\033[4mHeading 1\033[0m\nParagraph text that spans two lines. \n\n\033[3mNew\033[0m paragraph. \n\n\033[1m\033[4mHeader 1\033[0m\n\033[1m\033[3mFancy\033[0m\033[0m paragraph\nwith a \033[9mline break\033[0m. ",
		},
		{
			name:     "Exaggerated spacing + double-space line break + standard middle paragraph",
			input:    []string{"- First item", "\t- Sub item 1", "\t- Sub item 2", "- Second item", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "# Heading 1", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "Paragraph text", "that spans two lines.", "", "New paragraph.", "", "", "", "", "", "", "", "# Header 1", "__*Fancy*__ paragraph  ", "with a ~~line break~~."},
			expected: "• First item\n    • Sub item 1\n    • Sub item 2\n• Second item\n\n\033[1m\033[4mHeading 1\033[0m\nParagraph text that spans two lines. \n\nNew paragraph. \n\n\033[1m\033[4mHeader 1\033[0m\n\033[1m\033[3mFancy\033[0m\033[0m paragraph\nwith a \033[9mline break\033[0m. ",
		},
		{
			name:     "Paragraph followed by list (minimal spacing)",
			input:    []string{"This is a happy little paragraph", "- First item", "Second item", "\t- Sub item 1", "- Third item"},
			expected: "This is a happy little paragraph\n• First item\n• Second item\n    • Sub item 1\n• Third item",
		},
		{
			name:     "Paragraph followed by list (standard spacing)",
			input:    []string{"This is a happy little paragraph", "", "- First item", "Second item", "\t- Sub item 1", "- Third item"},
			expected: "This is a happy little paragraph\n• First item\n• Second item\n    • Sub item 1\n• Third item",
		},
		{
			name:     "Paragraph followed by list (exaggerated spacing)",
			input:    []string{"This is a happy little paragraph", "", "", "", "", "", "", "", "", "- First item", "Second item", "\t- Sub item 1", "- Third item"},
			expected: "This is a happy little paragraph\n• First item\n• Second item\n    • Sub item 1\n• Third item",
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
