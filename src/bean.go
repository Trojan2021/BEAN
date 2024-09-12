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
// Fix misnumbering of ordered lists when they follow an unordered list at the same indentation level (only affects child lists)
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

// RenderMarkdown renders the input lines as Markdown formatted for the CLI using ANSI escape codes.
func RenderMarkdown(lines []string) string {
	// variables to keep value between iterations
	/// ALL
	var output strings.Builder // stores the work-in-progress final output string
	/// LISTS
	var prevIndentMultiplier int // stores the value of the previous indentation multiplier
	var bullet string            // stores the bullet character for lists
	var linesLength = len(lines) // stores the length of the input lines slice
	/// LISTS: ORDERED
	var orderedIterator int          // stores the current number of the ordered list item
	var orderedIteratorHistory []int // stores the history of ordered list items

	// calcIndentMultiplier calculates the visual indentation level of a line and returns the result as an integer.
	// It also returns a boolean value indicating whether the line is valid Markdown.
	calcIndentMultiplier := func(indentSubstring string) (bool, int) {
		var indentMultiplier int
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
		return true, indentMultiplier
	}

	// updateOrderedIteratorHistory updates the history of ordered list items based on changes in indentation level.
	updateOrderedIteratorHistory := func(indentMultiplier int) {
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

	// renderParagraph returns the input line as a Markdown paragraph-formatted string.
	renderParagraph := func(line string) string {
		var outputString string
		if strings.TrimSpace(line) == "" {
			// if line is empty, skip a line
			outputString = "\n\n"
		} else if strings.TrimRight(line, "  ") != line {
			// if line ends in two+ spaces, write the line with a newline character
			outputString = line + "\n"
		} else if strings.HasSuffix(line, "<br>") {
			// if line ends in <br>, write the line with a newline character and strip the <br> tag
			outputString = line[:len(line)-4] + "\n"
		} else {
			// if line is not empty, write the line with a space at the end (for paragraph formatting)
			outputString = line + " "
		}

		// reset list variables (safe since lists are never rendered as paragraphs)
		prevIndentMultiplier = 0
		orderedIterator = 0
		orderedIteratorHistory = nil

		return outputString
	}

	// REGEX DICTIONARY
	// level 1 header (1)
	h1 := regexp.MustCompile(`^\s*# (.*)`)
	// level 2 header (2)
	h2 := regexp.MustCompile(`^\s*## (.*)`)
	// bold text (7)
	bold := regexp.MustCompile(`^(.*)(\*\*.+?\*\*|__.+?__)(.*)`)
	// italic text (8)
	italic := regexp.MustCompile(`^(.*)(\*.+?\*|_.+?_)(.*)`)
	// strikethrough text (9)
	strikethrough := regexp.MustCompile(`^(.*)~~(.+?)~~(.*)`)
	// (un)ordered list item (10)
	list := regexp.MustCompile(fmt.Sprintf(`^((?:\s{%d})*|\t+)([-+*] |\d+\. )(.*)`, indentSpaces))

	// iterate over lines
	for i, line := range lines {

		// render each matched Markdown element in the current line
		var internalOutput = line // stores the work-in-progress output for the current line
		var doNotRenderParagraph bool
		if h1.MatchString(internalOutput) {
			internalOutput = "\033[1m\033[4m" + h1.FindStringSubmatch(internalOutput)[1] + "\033[0m\n"
			doNotRenderParagraph = true
		}
		if h2.MatchString(internalOutput) {
			internalOutput = "\033[1m" + h2.FindStringSubmatch(internalOutput)[1] + "\033[0m\n"
			doNotRenderParagraph = true
		}
		for m := 0; m < 1; {
			if bold.MatchString(internalOutput) { // bold must be rendered first to avoid matching as italic
				substrings := bold.FindStringSubmatch(internalOutput)
				internalOutput = substrings[1] + "\033[1m" + substrings[2][2:len(substrings[2])-2] + "\033[0m" + substrings[3]
			} else {
				m++
			}
		}
		for m := 0; m < 1; {
			if italic.MatchString(internalOutput) {
				substrings := italic.FindStringSubmatch(internalOutput)
				internalOutput = substrings[1] + "\033[3m" + substrings[2][1:len(substrings[2])-1] + "\033[0m" + substrings[3]
			} else {
				m++
			}
		}
		for m := 0; m < 1; {
			if strikethrough.MatchString(internalOutput) {
				substrings := strikethrough.FindStringSubmatch(internalOutput)
				internalOutput = substrings[1] + "\033[9m" + substrings[2] + "\033[0m" + substrings[3]
			} else {
				m++
			}
		}
		if list.MatchString(internalOutput) {
			substrings := list.FindStringSubmatch(internalOutput)

			validMarkdown, indentMultiplier := calcIndentMultiplier(substrings[1])
			if !validMarkdown {
				// do nothing (do not process as list item)
				break
			}
			doNotRenderParagraph = true

			switch substrings[2][0] {
			case '-', '+', '*':
				// operations to take for unordered lists
				bullet = "• "

				if substrings[1] == "" {
					// if the item is an unordered list parent, reset the orderedIterator and its history
					orderedIterator = 0
					orderedIteratorHistory = nil
				} else if indentMultiplier != prevIndentMultiplier {
					// otherwise, if changing the indentation level, update the history of ordered list iterators
					// must be done for compatibility with mixed ordered/unordered lists
					updateOrderedIteratorHistory(indentMultiplier)
				}
			default:
				// operations to take for ordered lists
				if indentMultiplier == prevIndentMultiplier {
					// if not changing the indentation level, increment the iterator
					orderedIterator++
				} else {
					// otherwise, update the history of ordered list iterators
					updateOrderedIteratorHistory(indentMultiplier)
				}
				bullet = strconv.Itoa(orderedIterator) + ". "
			}

			// determine if the next line is a list item (to determine if a newline character is needed after the current list item)
			var lineEnding string
			if linesLength >= i+2 && strings.TrimSpace(lines[i+1]) != "" {
				lineEnding = "\n"
			} else {
				lineEnding = ""
			}

			// write the list item with the appropriate indentation
			internalOutput = strings.Repeat(" ", indentMultiplier*4) + bullet + substrings[3] + lineEnding

			// supply information for next line iteration
			prevIndentMultiplier = indentMultiplier
		}

		if !doNotRenderParagraph {
			output.WriteString(renderParagraph(internalOutput))
		} else {
			output.WriteString(internalOutput)
		}
	}
	return output.String()
}
