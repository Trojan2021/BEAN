package test

import (
	"testing"

	bean "github.com/Trojan2021/BEAN/render"
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
			expected: "• List item\n",
		},
		{
			name:     "Deep unordered list",
			input:    []string{"- List item", "    - Sub item 1", "        - Sub item 2", "            - Sub item 3", "                - Sub item 4", "                    - Sub item 5", "                        - Sub item 6", "                            - Sub item 7", "                                - Sub item 8", "                                    - Sub item 9", "                                    - Sub item 9", "                                    - Sub item 9", "                                    - Sub item 9", "                                        - Sub item 10", "                                        - Sub item 10", "                                        - Sub item 10", "                                        - Sub item 10"},
			expected: "• List item\n    ‣ Sub item 1\n        • Sub item 2\n            ‣ Sub item 3\n                • Sub item 4\n                    ‣ Sub item 5\n                        • Sub item 6\n                            ‣ Sub item 7\n                                • Sub item 8\n                                    ‣ Sub item 9\n                                    ‣ Sub item 9\n                                    ‣ Sub item 9\n                                    ‣ Sub item 9\n                                        • Sub item 10\n                                        • Sub item 10\n                                        • Sub item 10\n                                        • Sub item 10\n",
		},
		{
			name:     "Ordered list",
			input:    []string{"1. List item"},
			expected: "1. List item\n",
		},
		{
			name:     "Deep ordered list",
			input:    []string{"1. List item", "    1. Sub item 1", "        1. Sub item 2", "            1. Sub item 3", "                1. Sub item 4", "                    1. Sub item 5", "                        1. Sub item 6", "                            1. Sub item 7", "                                1. Sub item 8", "                                    1. Sub item 9", "                                    1. Sub item 9", "                                    1. Sub item 9", "                                    1. Sub item 9", "                                        1. Sub item 10", "                                        1. Sub item 10", "                                        1. Sub item 10", "                                        1. Sub item 10"},
			expected: "1. List item\n    i. Sub item 1\n        1. Sub item 2\n            i. Sub item 3\n                1. Sub item 4\n                    i. Sub item 5\n                        1. Sub item 6\n                            i. Sub item 7\n                                1. Sub item 8\n                                    i. Sub item 9\n                                    ii. Sub item 9\n                                    iii. Sub item 9\n                                    iv. Sub item 9\n                                        1. Sub item 10\n                                        2. Sub item 10\n                                        3. Sub item 10\n                                        4. Sub item 10\n",
		},
		{
			name:     "Unordered list following header w/different spacing",
			input:    []string{"# Header 1", "- List item", "# Header 1", "", "- List item", "# Header 1", "", "", "- List item"},
			expected: "\033[1m─Header 1─\033[0m\n• List item\n\n\033[1m─Header 1─\033[0m\n• List item\n\n\033[1m─Header 1─\033[0m\n• List item\n",
		},
		{
			name:     "Unordered list with multiple items",
			input:    []string{"- First item", "    - Sub item 1", "    - Sub item 2", "- Second item"},
			expected: "• First item\n    ‣ Sub item 1\n    ‣ Sub item 2\n• Second item\n",
		},
		{
			name:     "Unordered list with multiple items (tab input)",
			input:    []string{"- First item", "\t- Sub item 1", "\t- Sub item 2", "- Second item"},
			expected: "• First item\n    ‣ Sub item 1\n    ‣ Sub item 2\n• Second item\n",
		},
		{
			name:     "Ordered list with multiple items",
			input:    []string{"1. First item", "2. Second item", "3. Third item"},
			expected: "1. First item\n2. Second item\n3. Third item\n",
		},
		{
			name:     "Ordered list with sub-items",
			input:    []string{"1. First item", "    1. Sub item 1", "    2. Sub item 2", "2. Second item"},
			expected: "1. First item\n    1. Sub item 1\n    2. Sub item 2\n2. Second item\n",
		},
		{
			name:     "Ordered list with sub-items (tab input)",
			input:    []string{"1. First item", "\t1. Sub item 1", "\t2. Sub item 2", "2. Second item"},
			expected: "1. First item\n    1. Sub item 1\n    2. Sub item 2\n2. Second item\n",
		},
		{
			name:     "Mixed list (unordered inside ordered)",
			input:    []string{"1. First item", "    - Sub item 1", "    - Sub item 2", "2. Second item"},
			expected: "1. First item\n    ‣ Sub item 1\n    ‣ Sub item 2\n2. Second item\n",
		},
		{
			name:     "Mixed list (unordered inside ordered, tab input)",
			input:    []string{"1. First item", "\t- Sub item 1", "\t- Sub item 2", "2. Second item"},
			expected: "1. First item\n    ‣ Sub item 1\n    ‣ Sub item 2\n2. Second item\n",
		},
		{
			name:     "Mixed list (unordered AND ordered inside ordered, tab input)",
			input:    []string{"1. First item", "\t- Sub item 1", "\t1. Sub item 2", "\t1. Sub item 3"},
			expected: "1. First item\n    ‣ Sub item 1\n    1. Sub item 2\n    2. Sub item 3\n",
		},
		{
			name:     "Mixed list(s) (unordered AND ordered inside ordered, tab input, separated by headers and spaces)",
			input:    []string{"1. First item", "\t- Sub item 1", "\t1. Sub item 2", "\t1. Sub item 3", "# Header 1", "1. First item", "\t- Sub item 1", "\t1. Sub item 2", "\t1. Sub item 3", "", "1. First item", "\t- Sub item 1", "\t1. Sub item 2", "\t1. Sub item 3", "1. Second item", "\t- Sub item 1", "\t1. Sub item 2", "\t1. Sub item 3"},
			expected: "1. First item\n    ‣ Sub item 1\n    1. Sub item 2\n    2. Sub item 3\n\n\033[1m─Header 1─\033[0m\n1. First item\n    ‣ Sub item 1\n    1. Sub item 2\n    2. Sub item 3\n\n1. First item\n    ‣ Sub item 1\n    1. Sub item 2\n    2. Sub item 3\n2. Second item\n    ‣ Sub item 1\n    1. Sub item 2\n    2. Sub item 3\n",
		},
		{
			name:     "Mixed list (ordered inside unordered)",
			input:    []string{"- First item", "    1. Sub item 1", "    2. Sub item 2", "- Second item"},
			expected: "• First item\n    1. Sub item 1\n    2. Sub item 2\n• Second item\n",
		},
		{
			name:     "Mixed list (ordered inside unordered, tab input)",
			input:    []string{"- First item", "\t1. Sub item 1", "\t2. Sub item 2", "- Second item"},
			expected: "• First item\n    1. Sub item 1\n    2. Sub item 2\n• Second item\n",
		},
		// Lists with Bold, Italic, and Strikethrough
		{
			name:     "Unordered list with bold text",
			input:    []string{"- **Bold item**", "    - **Bold sub-item**"},
			expected: "• \033[1mBold item\033[0m\n    ‣ \033[1mBold sub-item\033[0m\n",
		},
		{
			name:     "Ordered list with italic text",
			input:    []string{"1. _Italic item_", "    2. _Italic sub-item_"},
			expected: "1. \033[3mItalic item\033[0m\n    1. \033[3mItalic sub-item\033[0m\n",
		},
		{
			name:     "Mixed list with bold and italic",
			input:    []string{"- **Bold item**", "    1. _Italic sub-item_", "    2. **Bold sub-item**"},
			expected: "• \033[1mBold item\033[0m\n    1. \033[3mItalic sub-item\033[0m\n    2. \033[1mBold sub-item\033[0m\n",
		},
		{
			name:     "Unordered list with strikethrough",
			input:    []string{"- ~~Strikethrough item~~", "    - ~~Strikethrough sub-item~~"},
			expected: "• \033[9mStrikethrough item\033[0m\n    ‣ \033[9mStrikethrough sub-item\033[0m\n",
		},
		{
			name: "Mixed list with bold, italic, and strikethrough",
			input: []string{
				"- **Bold item**",
				"    - _Italic sub-item_",
				"    - ~~Strikethrough sub-item~~",
				"- Mixed **_Bold and Italic_** item"},
			expected: "• \033[1mBold item\033[0m\n    ‣ \033[3mItalic sub-item\033[0m\n    ‣ \033[9mStrikethrough sub-item\033[0m\n• Mixed \033[1m\033[3mBold and Italic\033[0m\033[0m item\n",
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
