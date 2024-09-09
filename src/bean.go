package bean

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

// calcIndentMultiplier calculates the visual indentation level of a line and returns the result as an integer.
// It also returns a boolean value indicating whether the line is valid Markdown.
func calcIndentMultiplier(prevLineType, currentLineType uint8, prevIndentMultiplier int, indentSubstring string) (bool, int) {
	var indentMultiplier int
	if prevLineType == currentLineType { // if line is not list parent...
		// count tabs used for indentation
		tabCount := strings.Count(indentSubstring, "\t")
		// store the visual indentation level
		if tabCount > 0 {
			// if tabs were used for indentation, set indentMultiplier to the number of tabs
			indentMultiplier = tabCount
		} else {
			// if spaces were used for indentation, set indentMultiplier to the number of spaces divided by the indentSpaces constant
			indentMultiplier = len(indentSubstring) / indentSpaces
		}
		// if line is indented by more than one level past the previous line, it is not valid Markdown
		if prevIndentMultiplier+1 < indentMultiplier {
			return false, indentMultiplier
		}
	} else { // if line is list parent...
		// if first list item is indented, it is not valid Markdown
		if indentSubstring != "" {
			return false, indentMultiplier
		} else { // reset visual indentation level
			indentMultiplier = 0
		}
	}

	return true, indentMultiplier
}

// RenderMarkdown converts Markdown lines to a CLI-friendly format and returns the result as a string.
func RenderMarkdown(lines []string) string {
	var output strings.Builder

	// REGEX DICTIONARY
	// level 1 header
	h1 := regexp.MustCompile(`^\s*# (.*)`)
	// level 2 header
	h2 := regexp.MustCompile(`^\s*## (.*)`)
	// (un)ordered list
	list := regexp.MustCompile(fmt.Sprintf(`^((?:\s{%d})*|\t+)([-+*] |\d+\. )(.*)`, indentSpaces))

	// ALL
	var prevLineType uint8 // 0 = unimportant, 1 = list item
	// LISTS
	var prevIndentMultiplier int // stores the value of the previous indentation multiplier
	var bullet string            // stores the bullet character for lists
	// LISTS: ORDERED
	var orderedIterator int // stores the current number of the ordered list item
	for _, line := range lines {

		switch {
		case h1.MatchString(line):
			output.WriteString("\033[1m\033[4m" + h1.FindStringSubmatch(line)[1] + "\033[0m\n")
			prevLineType = 0
			orderedIterator = 0
		case h2.MatchString(line):
			output.WriteString("\033[1m" + h2.FindStringSubmatch(line)[1] + "\033[0m\n")
			prevLineType = 0
			orderedIterator = 0
		case list.MatchString(line):
			// save substrings matched by regex for later reference
			substrings := list.FindStringSubmatch(line)

			validMarkdown, indentMultiplier := calcIndentMultiplier(prevLineType, 1, prevIndentMultiplier, substrings[1])
			if !validMarkdown {
				output.WriteString(line + "\n")
				prevLineType = 0
				break
			}

			if substrings[2][0] != '-' && substrings[2][0] != '+' && substrings[2][0] != '*' {
				// operations to take for ordered lists
				if indentMultiplier == prevIndentMultiplier {
					orderedIterator++
				} else {
					orderedIterator = 1
				}
				bullet = strconv.Itoa(orderedIterator) + ". "
			} else {
				// operations to take for unordered lists
				bullet = "• "
				orderedIterator = 0
			}

			// write the list item with the appropriate indentation
			output.WriteString(strings.Repeat(" ", indentMultiplier*4) + bullet + substrings[3] + "\n")

			// supply information for next line iteration
			prevIndentMultiplier = indentMultiplier
			prevLineType = 1 // set prevLineType to 1 (list item)
		default:
			output.WriteString(line + "\n")
			prevLineType = 0
			orderedIterator = 0
		}
	}

	return output.String()
}
