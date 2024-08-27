package bean

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "# "):
			output.WriteString("\033[1m\033[4m" + strings.TrimPrefix(line, "# ") + "\033[0m\n")
		case strings.HasPrefix(line, "## "):
			output.WriteString("\033[1m" + strings.TrimPrefix(line, "## ") + "\033[0m\n")
		case strings.HasPrefix(line, "- "):
			output.WriteString("• " + strings.TrimPrefix(line, "- ") + "\n")
		default:
			output.WriteString(line + "\n")
		}
	}

	return output.String()
}
