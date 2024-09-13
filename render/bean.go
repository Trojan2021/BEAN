package render

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
	var prevElements [2]uint8  // stores the integer representation of the last rendered element [0] and the last rendered element different from the most recent [1] (0 = paragraph, 1 = header, 10 = list item, 255 = none)
	var matchedSomething bool  // indicates whether the current line matched any Markdown syntax
	/// LISTS
	var prevIndentMultiplier int // stores the value of the previous indentation multiplier
	var prevListWasOrdered bool  // stores whether the previous list was ordered
	var bullet string            // stores the bullet character for lists
	/// LISTS: ORDERED
	var orderedIterator = 1          // stores the current number of the ordered list item
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

	// updatePrevElements updates the prevElements slice with the most recent element.
	// It also sets matchedSomething to true to indicate that the current line matched some Markdown syntax.
	updatePrevElements := func(element uint8) {
		if prevElements[0] != element {
			prevElements[1] = prevElements[0]
			prevElements[0] = element
		}
		matchedSomething = true
	}

	// resetListVariables resets the variables used for list rendering.
	// It should be called whenever a list is not being rendered (currently headers and paragraphs).
	resetListVariables := func() {
		// TODO once lists and paragraphs are made mutually exclusive, simply perform this action whenever not rendering a list
		prevIndentMultiplier = 0
		orderedIterator = 1
		orderedIteratorHistory = nil
		prevListWasOrdered = false
	}

	// renderParagraph returns the input line as a Markdown paragraph-formatted string.
	renderParagraph := func(lineNumber int, lines *[]string, lineInProgress string) string {
		// trim spaces from current and previous line (later used to determine if they are empty)
		currentLineTrimmed := strings.TrimSpace(lineInProgress)
		if currentLineTrimmed == "" {
			// do not render empty lines; indicate that the previous element did not match any Markdown syntax
			updatePrevElements(255)
			resetListVariables() // break lists on empty lines
			return ""
		}

		// delcare variables to store work-in-progress strings
		var lineBeginning string
		var outputString string

		// headers require an extra newline character to break away from paragraphs
		// determine if following a paragraph

		// begin line with two newline characters if starting a new paragraph following another paragraph
		if currentLineTrimmed != "" {
			if prevElements[0] == 255 && prevElements[1] == 0 {
				lineBeginning = "\n\n"
			} else {
				// do not begin line with newline character if the previous line is part of the current paragraph
				lineBeginning = ""
			}
		} else {
			// do not begin line with newline character if the previous line is part of the current paragraph
			// (the previous line is not empty)
			// or if the current line is empty
			lineBeginning = ""
		}

		if strings.TrimRight(lineInProgress, "  ") != lineInProgress {
			// if line ends in two+ spaces, write the line with a newline character
			outputString = strings.TrimRight(lineInProgress, " ") + "\n"
		} else if strings.HasSuffix(lineInProgress, "<br>") {
			// if line ends in <br>, write the line with a newline character and strip the <br> tag
			outputString = lineInProgress[:len(lineInProgress)-4] + "\n"
		} else {
			// if line is not empty, write the line with a space at the end (for paragraph formatting)
			outputString = lineInProgress + " "
		}

		resetListVariables()

		updatePrevElements(0) // indicate that the current line is a paragraph (it passed the blank line check)

		return lineBeginning + outputString
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

	// getHeaderBeginning returns the appropriate number of newline characters to begin a header.
	getHeaderBeginning := func(lineNumber int) string {
		// since this function is run whenever a header is being rendered, reset list variables
		resetListVariables()

		// do not begin line with newline character if it is the first line
		if lineNumber == 0 {
			return ""
		}

		// headers require only one newline character if the previous element is a list or a header
		// this is because all list items and headers end in a newline character]
		if prevElements[0] == 10 || prevElements[0] == 1 || (prevElements[0] == 255 && (prevElements[1] == 10 || prevElements[1] == 1)) {
			return "\n"
		}

		// headers require an extra newline character to break away from paragraphs
		// determine if following a paragraph
		if prevElements[0] == 0 || (prevElements[0] == 255 && prevElements[1] == 0) {
			return "\n\n"
		}

		return ""
	}

	// iterate over lines
	for i, line := range lines {

		// render each matched Markdown element in the current line
		matchedSomething = false
		var internalOutput = line // stores the work-in-progress output for the current line
		if h1.MatchString(internalOutput) {
			internalOutput = getHeaderBeginning(i) + "\033[1m\033[4m" + h1.FindStringSubmatch(internalOutput)[1] + "\033[0m\n"
			updatePrevElements(1)
		}
		if h2.MatchString(internalOutput) {
			internalOutput = getHeaderBeginning(i) + "\033[1m" + h2.FindStringSubmatch(internalOutput)[1] + "\033[0m\n"
			updatePrevElements(1)
		}
		for {
			if bold.MatchString(internalOutput) { // bold must be rendered first to avoid matching as italic
				substrings := bold.FindStringSubmatch(internalOutput)
				internalOutput = substrings[1] + "\033[1m" + substrings[2][2:len(substrings[2])-2] + "\033[0m" + substrings[3]
				// do not update prevElements since bold/italic/strikethrough is part of a paragraph (consider it unmatched so that renderParagraph can handle it properly)
			} else {
				break
			}
		}
		for {
			if italic.MatchString(internalOutput) {
				substrings := italic.FindStringSubmatch(internalOutput)
				internalOutput = substrings[1] + "\033[3m" + substrings[2][1:len(substrings[2])-1] + "\033[0m" + substrings[3]
			} else {
				break
			}
		}
		for {
			if strikethrough.MatchString(internalOutput) {
				substrings := strikethrough.FindStringSubmatch(internalOutput)
				internalOutput = substrings[1] + "\033[9m" + substrings[2] + "\033[0m" + substrings[3]
			} else {
				break
			}
		}
		if list.MatchString(internalOutput) {
			substrings := list.FindStringSubmatch(internalOutput)

			validMarkdown, indentMultiplier := calcIndentMultiplier(substrings[1])
			if !validMarkdown {
				// do nothing (do not process as list item)
				break
			}

			switch substrings[2][0] {
			case '-', '+', '*':
				// operations to take for unordered lists
				bullet = "• "

				if substrings[1] == "" {
					// if the item is an unordered list parent, reset the orderedIterator and its history
					orderedIterator = 1
					orderedIteratorHistory = nil
				} else if indentMultiplier != prevIndentMultiplier {
					// otherwise, if changing the indentation level, update the history of ordered list iterators
					// must be done for compatibility with mixed ordered/unordered lists
					updateOrderedIteratorHistory(indentMultiplier)
				}

				prevListWasOrdered = false
			default:
				// operations to take for ordered lists
				if indentMultiplier == prevIndentMultiplier {
					// if not changing the indentation level, increment the iterator
					if prevListWasOrdered {
						orderedIterator++
					}
				} else {
					// otherwise, update the history of ordered list iterators
					updateOrderedIteratorHistory(indentMultiplier)
				}
				bullet = strconv.Itoa(orderedIterator) + ". "

				prevListWasOrdered = true
			}

			// if the previous line is not a list item/header OR if the previous line is blank but the last matched element was not a list item/header,
			// preceed the list item with a newline character
			var lineBeginning string
			if i != 0 && (prevElements[0] != 10 && prevElements[0] != 1) || (prevElements[0] == 255 && (prevElements[1] != 10 && prevElements[1] != 1)) {
				lineBeginning = "\n"
			}

			// write the list item with the appropriate indentation
			internalOutput = lineBeginning + strings.Repeat(" ", indentMultiplier*4) + bullet + substrings[3] + "\n"

			// supply information for next line iteration
			prevIndentMultiplier = indentMultiplier
			updatePrevElements(10)
		}

		// determine whether to render line as paragraph
		if !matchedSomething {
			// render as paragraph if no Markdown was matched or if a paragraph was explicitly matched
			output.WriteString(renderParagraph(i, &lines, internalOutput))
		} else {
			output.WriteString(internalOutput)
		}
	}

	return output.String()
}
