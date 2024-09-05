package bean

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"regexp"
)

// ReadFile reads the markdown file and returns its lines as a slice of strings.
func ReadFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", fileName, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return lines, nil
}

// RenderMarkdown converts markdown lines to a CLI-friendly format.
func RenderMarkdown(lines []string) string {
	var output strings.Builder

	// regex dictionary

	// level 1 header
	h1 := regexp.MustCompile(`^\s*# (.*)`)
	// level 2 header
	h2 := regexp.MustCompile(`^\s*## (.*)`)
	// unordered list
	list := regexp.MustCompile(`^((\s\s\s\s)*|\t+)[-+*] (.*)`)

	var lastLineType uint8 // 0 = unimportant, 1 = list item
	for _, line := range lines {

		switch {
		case h1.MatchString(line):
			output.WriteString("\033[1m\033[4m" + h1.FindStringSubmatch(line)[1] + "\033[0m\n")
			lastLineType = 0
		case h2.MatchString(line):
			output.WriteString("\033[1m" + h2.FindStringSubmatch(line)[1] + "\033[0m\n")
			lastLineType = 0
		case list.MatchString(line):
			// TODO Unordered list rendering:
			// Disallow indenting a list item by more than one level
			// Convert tabs to groups of 4 spaces
			// Wrap lists with handing indentation
			// Optionally support detecting how many spaces equal a tab

			// save substrings matched by regex for later reference
			substrings := list.FindStringSubmatch(line)

			// determine whether to allow indented/nested list items based on if the previous line was a list item
			var visualIndent string
			if lastLineType == 1 {
				visualIndent = substrings[1]
			} else {
				// do not allow first list item to be indented (print as-is)
				if substrings[1] != "" {
					output.WriteString(line + "\n")
					break
				}
			}

			// write the list item with the appropriate indentation
			output.WriteString(visualIndent + "• " + substrings[3] + "\n")

			// indicate that the last line was part of an unordered list
			lastLineType = 1
		default:
			output.WriteString(line + "\n")
			lastLineType = 0
		}
	}

	return output.String()
}
