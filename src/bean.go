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

// RenderMarkdown converts Markdown lines to a CLI-friendly format and returns the result as a string.
func RenderMarkdown(lines []string) string {
	var output strings.Builder

	// REGEX DICTIONARY
	// level 1 header
	h1 := regexp.MustCompile(`^\s*# (.*)`)
	// level 2 header
	h2 := regexp.MustCompile(`^\s*## (.*)`)
	// unordered list
	list := regexp.MustCompile(fmt.Sprintf(`^((\s{%d})*|\t+)[-+*] (.*)`, indentSpaces))

	var lastLineType uint8       // 0 = unimportant, 1 = list item (more to be implemented)
	var indentMultiplier int     //stores the value of the indentation multiplier
	var prevIndentMultiplier int // stores the value of the previous indentation multiplier
	for _, line := range lines {

		switch {
		case h1.MatchString(line):
			output.WriteString("\033[1m\033[4m" + h1.FindStringSubmatch(line)[1] + "\033[0m\n")
			lastLineType = 0
		case h2.MatchString(line):
			output.WriteString("\033[1m" + h2.FindStringSubmatch(line)[1] + "\033[0m\n")
			lastLineType = 0
		case list.MatchString(line):
			// save substrings matched by regex for later reference
			substrings := list.FindStringSubmatch(line)

			// calculate the visual indentation level
			if lastLineType == 1 { // if line is not list parent...
				// count tabs used for indentation
				tabCount := strings.Count(substrings[1], "\t")
				// store the visual indentation level
				if tabCount > 0 {
					// if tabs were used for indentation, set indentMultiplier to the number of tabs
					indentMultiplier = tabCount
				} else {
					// if spaces were used for indentation, set indentMultiplier to the number of spaces divided by the indentSpaces constant
					indentMultiplier = len(substrings[1]) / indentSpaces
				}
				// if line is indented by more than one level past the previous line, it is not valid Markdown and should be printed as-is
				if prevIndentMultiplier+1 < indentMultiplier {
					// TODO reduce re-used code
					output.WriteString(line + "\n")
					lastLineType = 0
					break
				}
			} else { // if line is list parent...
				// if first list item is indented, it is not valid Markdown and should be printed as-is
				if substrings[1] != "" {
					// TODO reduce re-used code
					output.WriteString(line + "\n")
					lastLineType = 0
					break
				} else { // reset visual indentation level
					indentMultiplier = 0
				}
			}

			// write the list item with the appropriate indentation
			output.WriteString(strings.Repeat(" ", indentMultiplier*4) + "• " + substrings[3] + "\n")

			// supply information for next line iteration
			prevIndentMultiplier = indentMultiplier
			lastLineType = 1 // set lastLineType to 1 (list item)
		default:
			output.WriteString(line + "\n")
			lastLineType = 0
		}
	}

	return output.String()
}
