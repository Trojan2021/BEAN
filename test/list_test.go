package bean

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/src"
)

func TestRenderLists(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
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
		// Lists with Bold, Italic, and Strikethrough
		{
			name:     "Unordered list with bold text",
			input:    []string{"- **Bold item**", "    - **Bold sub-item**"},
			expected: "• \033[1mBold item\033[0m\n    • \033[1mBold sub-item\033[0m",
		},
		{
			name:     "Ordered list with italic text",
			input:    []string{"1. _Italic item_", "    2. _Italic sub-item_"},
			expected: "1. \033[3mItalic item\033[0m\n    2. \033[3mItalic sub-item\033[0m",
		},
		{
			name:     "Mixed list with bold and italic",
			input:    []string{"- **Bold item**", "    1. _Italic sub-item_", "    2. **Bold sub-item**"},
			expected: "• \033[1mBold item\033[0m\n    1. \033[3mItalic sub-item\033[0m\n    2. \033[1mBold sub-item\033[0m",
		},
		{
			name:     "Unordered list with strikethrough",
			input:    []string{"- ~~Strikethrough item~~", "    - ~~Strikethrough sub-item~~"},
			expected: "• \033[9mStrikethrough item\033[0m\n    • \033[9mStrikethrough sub-item\033[0m",
		},
		{
			name: "Mixed list with bold, italic, and strikethrough",
			input: []string{
				"- **Bold item**",
				"    - _Italic sub-item_",
				"    - ~~Strikethrough sub-item~~",
				"- Mixed **_Bold and Italic_** item"},
			expected: "• \033[1mBold item\033[0m\n    • \033[3mItalic sub-item\033[0m\n    • \033[9mStrikethrough sub-item\033[0m\n• Mixed \033[1m\033[3mBold and Italic\033[0m item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := bean.RenderMarkdown(tt.input)
			if output != tt.expected {
				t.Errorf("got %q, want %q", output, tt.expected)
			}
		})
	}
}
