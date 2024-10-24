package render

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	wrap "github.com/mitchellh/go-wordwrap"
	"golang.org/x/term"
)

// TODO General:
// Wrap text to terminal width (or a specified percentage of it); always wrap lists with hanging indentation
// Optionally support auto-detection of tab (space) width; if compiled to do this, replace indentSpaces with a variable holding the detected value

// ReadFile reads the markdown file and returns its lines as a slice of strings.
func ReadFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", fileName, err)
	}
	defer func(file *os.File) {
		_ = file.Close() // error ignored; if the file could be opened, it can probably be closed
	}(file)

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
	var oBuffer strings.Builder                         // stores the work-in-progress final output string
	var prevElements [2]uint8                           // stores the integer representation of the last rendered element [0] and the last rendered element different from the most recent [1] (0 = paragraph, 1 = header, 2 = hr, 10 = list item, 255 = none)
	var matchedSomething bool                           // indicates whether the current line matched any Markdown syntax
	var width, _, _ = term.GetSize(int(os.Stdout.Fd())) // stores the terminal width
	/// PARAGRAPHS
	var pBuffer strings.Builder // stores work-in-progress paragraph text
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
		prevIndentMultiplier = 0
		orderedIterator = 1
		orderedIteratorHistory = nil
		prevListWasOrdered = false
	}

	// mergeBuffers merges the contents of pBuffer (word-wrapped) into oBuffer and resets pBuffer.
	mergeBuffers := func() {
		if pBuffer.Len() > 0 {
			oBuffer.WriteString(wrap.WrapString(pBuffer.String(), uint(width)))
			pBuffer.Reset()
		}
	}

	// renderParagraph renders the current line as a paragraph by managing pBuffer and oBuffer.
	renderParagraph := func(lineNumber int, lines *[]string, lineInProgress string) {
		// trim spaces from current and previous line (later used to determine if they are empty)
		currentLineTrimmed := strings.TrimSpace(lineInProgress)
		if currentLineTrimmed == "" {
			// do not render empty lines; indicate that the previous element did not match any Markdown syntax
			updatePrevElements(255)
			resetListVariables() // break lists on empty lines

			// since the current line is empty, a new paragraph is being started and buffers should be merged
			mergeBuffers()
			return
		}

		// determine if newline characters should be added before the current line
		if prevElements[0] == 255 && prevElements[1] == 0 {
			// begin line with two newline characters if starting a new paragraph following another paragraph
			pBuffer.WriteString("\n\n")
		} else if prevElements[0] == 10 || (prevElements[0] == 255 && prevElements[1] == 10) {
			// begin line with a newline character if starting a new paragraph following a list
			pBuffer.WriteString("\n")
		}

		if strings.TrimRight(lineInProgress, "  ") != lineInProgress {
			// if line ends in two+ spaces, write the line with a newline character
			pBuffer.WriteString(strings.TrimRight(lineInProgress, " ") + "\n")
		} else if strings.HasSuffix(lineInProgress, "<br>") {
			// if line ends in <br>, write the line with a newline character and strip the <br> tag
			pBuffer.WriteString(lineInProgress[:len(lineInProgress)-4] + "\n")
		} else {
			// if line does not contain a manual break, write the line with a space at the end (for paragraph formatting)
			pBuffer.WriteString(lineInProgress + " ")
		}

		resetListVariables()

		updatePrevElements(0) // indicate that the current line is a paragraph (it passed the blank line check)
	}

	// containsMultipleUniqueChars returns true if the input string contains multiple unique characters.
	// It is used to avoid matching horizontal rules as bold/italic paragraphs.
	containsMultipleUniqueChars := func(s string) bool {
		charMap := make(map[rune]bool)
		for _, char := range s {
			charMap[char] = true
		}
		if len(charMap) == 1 {
			return false
		} else {
			return true
		}
	}

	// REGEX DICTIONARY
	// level 1-6 header (1)
	header := regexp.MustCompile(`^\s*(#{1,6}) (.*)`)
	// horizontal rule (2)
	hr := regexp.MustCompile(`^(?:-{3,}|\*{3,}|_{3,})$`)
	// in-line code (paragraph code) (0)
	pCode := regexp.MustCompile("^(.*)`(.+?)`(.*)")
	// bold text (0)
	bold := regexp.MustCompile(`^(.*)(\*\*.+?\*\*|__.+?__)(.*)`)
	// italic text (0)
	italic := regexp.MustCompile(`^(.*)(\*.+?\*|_.+?_)(.*)`)
	// strikethrough text (0)
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

		// headers require no newline characters if the previous element is a horizontal rule
		// this is because horizontal rules always end in two newline characters
		if prevElements[0] == 2 || (prevElements[0] == 255 && prevElements[1] == 2) {
			return ""
		}

		// headers require only one newline character if the previous element is a list or a header
		// this is because all list items and headers end in a newline character
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
		var renderedCodeBlocks []string

		// match elements that may be embedded within a paragraph
		for {
			if pCode.MatchString(internalOutput) {
				substrings := pCode.FindStringSubmatch(internalOutput)

				// temporarily replace in-line code block with \x1f (will be restored later to avoid processing Markdown elements within the code block)
				internalOutput = substrings[1] + "\x1f" + substrings[3]

				// save rendered in-line code block for later restoration
				renderedCodeBlocks = append(renderedCodeBlocks, "\033[48;5;238;38;5;1m"+substrings[2]+"\033[0m")

				// do not update prevElements since pCode/bold/italic/strikethrough is part of a paragraph (consider it unmatched so that renderParagraph can handle it properly)
			} else {
				break
			}
		}
		for {
			if bold.MatchString(internalOutput) { // bold must be rendered first to avoid matching as italic
				// avoid processing horizontal rules as bold text
				if !containsMultipleUniqueChars(internalOutput) {
					break
				}

				substrings := bold.FindStringSubmatch(internalOutput)
				internalOutput = substrings[1] + "\033[1m" + substrings[2][2:len(substrings[2])-2] + "\033[0m" + substrings[3]
			} else {
				break
			}
		}
		for {
			if italic.MatchString(internalOutput) {
				// avoid processing horizontal rules as italic text
				if !containsMultipleUniqueChars(internalOutput) {
					break
				}

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

		// restore in-line code blocks
		for i := len(renderedCodeBlocks) - 1; i >= 0; i-- {
			internalOutput = strings.Replace(internalOutput, "\x1f", renderedCodeBlocks[i], 1)
		}

		// match mutually exclusive elements that may NOT be embedded within a paragraph
		if header.MatchString(internalOutput) {
			// headers
			// determine header level (number of "#" characters at the beginning of the line)
			// and create a visual representation of the header
			substrings := header.FindStringSubmatch(internalOutput)
			headerLevel := len(substrings[1])
			visual := strings.Repeat("─", headerLevel)

			internalOutput = getHeaderBeginning(i) + "\033[1m" + visual + substrings[2] + visual + "\033[0m\n"
			updatePrevElements(1)
		} else if hr.MatchString(internalOutput) {
			// horizontal rule
			internalOutput = getHeaderBeginning(i) + strings.Repeat("─", width) + "\n\n"
			updatePrevElements(2)
		} else if list.MatchString(internalOutput) {
			// lists
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

			// if the previous line is a paragraph OR if the previous line is blank but the last matched element was not a header
			// (lists with a gap between them should be treated as separate lists), precede the list item with a newline character
			var lineBeginning string
			if i != 0 && ((prevElements[0] == 0) || (prevElements[0] == 255 && prevElements[1] != 1)) {
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
			renderParagraph(i, &lines, internalOutput)
		} else {
			// since a non-paragraph element was matched, merge pBuffer into oBuffer and reset pBuffer
			mergeBuffers()
			oBuffer.WriteString(internalOutput)
		}
	}

	// in case anything is left in pBuffer at the end, merge it into oBuffer
	mergeBuffers()
	return oBuffer.String()
}
