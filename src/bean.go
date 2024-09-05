package bean

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TODO General:
// Optionally support auto-detection of tab (space) width; if compiled to do this, replace indentSpaces with a variable holding the detected value
// Wrap text to terminal width (or a specified percentage of it); always wrap lists with hanging indentation

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
	list := regexp.MustCompile(fmt.Sprintf(`^((\s{%d})*|\t+)[-+*] (.*)`, indentSpaces))

	var lastLineType uint8    // 0 = unimportant, 1 = list item
	var visualIndentLevel int // on a new iteration, stores the value of the last visual indentation multiplier
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
			// Convert tabs to groups of spaces

			// save substrings matched by regex for later reference
			substrings := list.FindStringSubmatch(line)

			// calculate the visual indentation level
			if lastLineType == 1 { // if line is not list parent...
				// count tabs used for indentation
				tabCount := strings.Count(substrings[1], "\t")
				if tabCount > 0 {
					// if tabs were used for indentation, set visualIndentLevel to the number of tabs
					visualIndentLevel = tabCount
				} else {
					// if spaces were used for indentation, set visualIndentLevel to the number of spaces divided by the indentSpaces constant
					visualIndentLevel = len(substrings[1]) / indentSpaces
				}
			} else { // if line is list parent...
				// do not allow first list item to be indented (print as-is)
				// TODO fallthrough for readability
				if substrings[1] != "" {
					output.WriteString(line + "\n")
					lastLineType = 0
					break
				} else { // reset visual indentation level
					visualIndentLevel = 0
				}
			}

			// write the list item with the appropriate indentation
			output.WriteString(strings.Repeat(" ", visualIndentLevel*4) + "• " + substrings[3] + "\n")

			// indicate that the last line was part of an unordered list
			lastLineType = 1
		default:
			output.WriteString(line + "\n")
			lastLineType = 0
		}
	}

	return output.String()
}
