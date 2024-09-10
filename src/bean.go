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
// Allow combining multiple Markdown elements in a single line (bold and italic text, text within a list item, etc.)
// Wrap text to terminal width (or a specified percentage of it); always wrap lists with hanging indentation
// Optionally support auto-detection of tab (space) width; if compiled to do this, replace indentSpaces with a variable holding the detected value

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
	// variables to keep value between iterations
	/// ALL
	var output strings.Builder
	var prevLineType uint8 // 0 = unimportant, 1 = list item
	/// LISTS
	var prevIndentMultiplier int // stores the value of the previous indentation multiplier
	var bullet string            // stores the bullet character for lists
	/// LISTS: ORDERED
	var orderedIterator int          // stores the current number of the ordered list item
	var orderedIteratorHistory []int // stores the history of ordered list items

	// calcIndentMultiplier calculates the visual indentation level of a line and returns the result as an integer.
	// It also returns a boolean value indicating whether the line is valid Markdown.
	calcIndentMultiplier := func(currentLineType uint8, indentSubstring string) (bool, int) {
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

	manageOrderedIteratorHistory := func(indentMultiplier int) {
		if indentMultiplier > prevIndentMultiplier {
			// if indenting in, add the current orderedIterator to the history and reset the iterator
			orderedIteratorHistory = append(orderedIteratorHistory, orderedIterator)
			orderedIterator = 1
		} else {
			// if indenting out, determine how many levels
			outLevels := prevIndentMultiplier - indentMultiplier
			// pick up from the proper element in the history
			orderedIterator = orderedIteratorHistory[len(orderedIteratorHistory)-outLevels] + 1
			// remove the used elements from the history
			orderedIteratorHistory = orderedIteratorHistory[:len(orderedIteratorHistory)-outLevels]
		}
	}

	resetPreloopVariables := func(lineType uint8) {
		prevLineType = lineType
		prevIndentMultiplier = 0
		orderedIterator = 0
		orderedIteratorHistory = nil
	}

	renderParagraph := func(line string) {
		if strings.TrimSpace(line) == "" {
			// if line is empty, skip a line
			output.WriteString("\n\n")
		} else if strings.TrimRight(line, "  ") != line {
			// if line ends in two+ spaces, write the line with a newline character
			output.WriteString(line + "\n")
		} else if strings.HasSuffix(line, "<br>") {
			// if line ends in <br>, write the line with a newline character and strip the <br> tag
			output.WriteString(line[:len(line)-4] + "\n")
		} else {
			// if line is not empty, write the line with a space at the end (for paragraph formatting)
			output.WriteString(line + " ")
		}
		resetPreloopVariables(0)
	}

	// REGEX DICTIONARY
	// level 1 header
	h1 := regexp.MustCompile(`^\s*# (.*)`)
	// level 2 header
	h2 := regexp.MustCompile(`^\s*## (.*)`)
	// bold text
	bold := regexp.MustCompile(`^(.*)([^*]\*\*[^*].+?[^*]\*\*[^*]|[^_]__[^_].+?[^_]__[^_])(.*)`)
	// italic text
	italic := regexp.MustCompile(`^(.*)([^*]\*[^*].+?[^*]\*[^*]|[^_]_[^_].+?[^_]_[^_])(.*)`)
	// strikethrough text
	strikethrough := regexp.MustCompile(`^(.*)~~(.+?)~~(.*)`)
	// (un)ordered list item
	list := regexp.MustCompile(fmt.Sprintf(`^((?:\s{%d})*|\t+)([-+*] |\d+\. )(.*)`, indentSpaces))

	// loop over matched regex
	for _, line := range lines {
		switch {

		case h1.MatchString(line):
			output.WriteString("\033[1m\033[4m" + h1.FindStringSubmatch(line)[1] + "\033[0m\n")
			resetPreloopVariables(0)

		case h2.MatchString(line):
			output.WriteString("\033[1m" + h2.FindStringSubmatch(line)[1] + "\033[0m\n")
			resetPreloopVariables(0)

		case bold.MatchString(line):
			substrings := bold.FindStringSubmatch(line)
			renderParagraph(substrings[1] + "\033[1m " + substrings[2][3:len(substrings[2])-3] + " \033[0m" + substrings[3])

		case italic.MatchString(line):
			substrings := italic.FindStringSubmatch(line)
			renderParagraph(substrings[1] + "\033[3m " + substrings[2][2:len(substrings[2])-2] + " \033[0m" + substrings[3])

		case strikethrough.MatchString(line):
			substrings := strikethrough.FindStringSubmatch(line)
			renderParagraph(substrings[1] + "\033[9m" + substrings[2] + "\033[0m" + substrings[3])

		case list.MatchString(line):
			// save substrings matched by regex for later reference
			substrings := list.FindStringSubmatch(line)

			validMarkdown, indentMultiplier := calcIndentMultiplier(1, substrings[1])
			if !validMarkdown {
				renderParagraph(line)
				break
			}

			switch substrings[2][0] {
			case '-', '+', '*':
				// operations to take for unordered lists
				bullet = "• "

				if substrings[1] == "" {
					// if the item is an unordered list parent, reset the orderedIterator and its history
					orderedIterator = 0
					orderedIteratorHistory = nil
				} else if indentMultiplier != prevIndentMultiplier {
					// otherwise, if changing the indentation level, manage the history of ordered list iterators
					// must be done for compatibility with mixed ordered/unordered lists
					manageOrderedIteratorHistory(indentMultiplier)
				}
			default:
				// operations to take for ordered lists
				if indentMultiplier == prevIndentMultiplier {
					// if not changing the indentation level, increment the iterator
					orderedIterator++
				} else {
					// otherwise, manage the history of ordered list iterators
					manageOrderedIteratorHistory(indentMultiplier)
				}
				bullet = strconv.Itoa(orderedIterator) + ". "
			}

			// write the list item with the appropriate indentation
			output.WriteString(strings.Repeat(" ", indentMultiplier*4) + bullet + substrings[3] + "\n")

			// supply information for next line iteration
			prevIndentMultiplier = indentMultiplier
			prevLineType = 1 // set prevLineType to 1 (list item)

		default:
			renderParagraph(line)
		}
	}

	return output.String()
}
