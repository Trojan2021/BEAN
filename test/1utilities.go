package test

import (
	"strings"
	"testing"
)

// specify a consistent terminal width for testing
const terminalWidth = 80

// bufferFailure adds a literal and an ANSI-interpreted representation of a failed test case to the log buffer.
func bufferFailure(t *testing.T, got, want string) {
	divider := strings.Repeat("=", 40)
	t.Errorf("\nLITERAL:\nGOT:\n%q\nWANT:\n%q\n\nVISUAL:\nGOT:\n%s\n%s\nWANT:\n%s\n%s", got, want, strings.ReplaceAll(got, " ", "\u2592"), divider, strings.ReplaceAll(want, " ", "\u2592"), divider)
}
