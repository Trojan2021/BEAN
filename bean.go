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
	// unordered list (hyphen)
	list := regexp.MustCompile(`^((\s\s\s\s)*|\t+)- (.*)`)

	for _, line := range lines {

		switch {
		case h1.MatchString(line):
			output.WriteString("\033[1m\033[4m" + h1.FindStringSubmatch(line)[1] + "\033[0m\n")
		case h2.MatchString(line):
			output.WriteString("\033[1m" + h2.FindStringSubmatch(line)[1] + "\033[0m\n")
		case list.MatchString(line):
			substrings := list.FindStringSubmatch(line)
			output.WriteString(substrings[1] + "• " + substrings[2] + "\n")
		default:
			output.WriteString(line + "\n")
		}
	}

	return output.String()
}
