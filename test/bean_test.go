package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
)

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		// Header Testing
		{
			name:     "Header level 1",
			input:    []string{"# Heading 1"},
			expected: "\033[1m\033[4mHeading 1\033[0m\n",
		},
		{
			name:     "Header level 2",
			input:    []string{"## Heading 2"},
			expected: "\033[1mHeading 2\033[0m\n",
		},
		{
			name:     "Plain text",
			input:    []string{"Just a line of text."},
			expected: "Just a line of text. ",
		},

		// List Testing
		{
			name:     "Unordered list",
			input:    []string{"- List item"},
			expected: "• List item",
		},
		{
			name:     "Ordered list",
			input:    []string{"1. List item"},
			expected: "1. List item",
		},
		{
			name:     "Unordered list with multiple items",
			input:    []string{"- First item", "    - Sub item 1", "    - Sub item 2", "- Second item"},
			expected: "• First item\n    • Sub item 1\n    • Sub item 2\n• Second item",
		},
		{
			name:     "Unordered list with multiple items (tab input)",
			input:    []string{"- First item", "\t- Sub item 1", "\t- Sub item 2", "- Second item"},
			expected: "• First item\n    • Sub item 1\n    • Sub item 2\n• Second item",
		},
		{
			name:     "Ordered list with multiple items",
			input:    []string{"1. First item", "2. Second item", "3. Third item"},
			expected: "1. First item\n2. Second item\n3. Third item",
		},
		{
			name:     "Ordered list with sub-items",
			input:    []string{"1. First item", "    1. Sub item 1", "    2. Sub item 2", "2. Second item"},
			expected: "1. First item\n    1. Sub item 1\n    2. Sub item 2\n2. Second item",
		},
		{
			name:     "Ordered list with sub-items (tab input)",
			input:    []string{"1. First item", "\t1. Sub item 1", "\t2. Sub item 2", "2. Second item"},
			expected: "1. First item\n    1. Sub item 1\n    2. Sub item 2\n2. Second item",
		},
		{
			name:     "Mixed list (unordered inside ordered)",
			input:    []string{"1. First item", "    - Sub item 1", "    - Sub item 2", "2. Second item"},
			expected: "1. First item\n    • Sub item 1\n    • Sub item 2\n2. Second item",
		},
		{
			name:     "Mixed list (unordered inside ordered, tab input)",
			input:    []string{"1. First item", "\t- Sub item 1", "\t- Sub item 2", "2. Second item"},
			expected: "1. First item\n    • Sub item 1\n    • Sub item 2\n2. Second item",
		},
		{
			name:     "Mixed list (ordered inside unordered)",
			input:    []string{"- First item", "    1. Sub item 1", "    2. Sub item 2", "- Second item"},
			expected: "• First item\n    1. Sub item 1\n    2. Sub item 2\n• Second item",
		},
		{
			name:     "Mixed list (ordered inside unordered, tab input)",
			input:    []string{"- First item", "\t1. Sub item 1", "\t2. Sub item 2", "- Second item"},
			expected: "• First item\n    1. Sub item 1\n    2. Sub item 2\n• Second item",
		},
		// Strikethroughs
		{
			name:     "Simple strikethrough",
			input:    []string{"~~Strikethrough text~~"},
			expected: "\033[9mStrikethrough text\033[0m ",
		},
		{
			name:     "Strikethrough with surrounding spaces",
			input:    []string{"~~ Strikethrough with spaces ~~"},
			expected: "\033[9m Strikethrough with spaces \033[0m ",
		},
		{
			name:     "Strikethrough at start of line",
			input:    []string{"~~Strikethrough at start~~ regular text"},
			expected: "\033[9mStrikethrough at start\033[0m regular text ",
		},
		{
			name:     "Strikethrough at end of line",
			input:    []string{"Regular text ~~strikethrough at end~~"},
			expected: "Regular text \033[9mstrikethrough at end\033[0m ",
		},
		{
			name:     "Strikethrough in middle of text",
			input:    []string{"Text with ~~strikethrough~~ in the middle"},
			expected: "Text with \033[9mstrikethrough\033[0m in the middle ",
		},
		{
			name:     "Multiple strikethrough sections",
			input:    []string{"~~First section~~ and ~~second section~~"},
			expected: "\033[9mFirst section\033[0m and \033[9msecond section\033[0m ",
		},
		{
			name:     "Nested strikethrough and bold",
			input:    []string{"~~**Bold and strikethrough**~~"},
			expected: "\033[9m\033[1mBold and strikethrough\033[0m\033[0m ",
		},
		{
			name:     "Strikethrough with inline code",
			input:    []string{"~~Strikethrough `inline code`~~"},
			expected: "\033[9mStrikethrough `inline code`\033[0m ",
		},
		{
			name:     "Incomplete strikethrough start",
			input:    []string{"~~Incomplete strikethrough"},
			expected: "~~Incomplete strikethrough ", // Not rendered as strikethrough
		},
		{
			name:     "Incomplete strikethrough end",
			input:    []string{"Incomplete strikethrough~~"},
			expected: "Incomplete strikethrough~~ ", // Not rendered as strikethrough
		},
		{
			name:     "Strikethrough with mixed formatting",
			input:    []string{"**Bold** and ~~strikethrough~~ and _italic_"},
			expected: "\033[1mBold\033[0m and \033[9mstrikethrough\033[0m and \033[3mitalic\033[0m ",
		},
		{
			name:     "Strikethrough inside a list",
			input:    []string{"- ~~Strikethrough list item~~"},
			expected: "• \033[9mStrikethrough list item\033[0m",
		},
		{
			name:     "Strikethrough with numbers",
			input:    []string{"~~12345~~"},
			expected: "\033[9m12345\033[0m ",
		},
		{
			name:     "Strikethrough with numbers and repeated",
			input:    []string{"~~12345~~~~12345~~"},
			expected: "\033[9m12345\033[0m\033[9m12345\033[0m ",
		},
		{
			name:     "Strikethrough with special characters",
			input:    []string{"~~!@#$%^&*()~~"},
			expected: "\033[9m!@#$%^&*()\033[0m ",
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
