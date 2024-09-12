package test

import (
	"strings"
	"testing"
)

// bufferFailure adds a literal and an ANSI-interpreted representation of a failed test case to the log buffer.
func bufferFailure(t *testing.T, got, want string) {
	t.Errorf("\nLITERAL:\nGOT:\n%q\nWANT:\n%q\n\nVISUAL:\nGOT:\n%s\nWANT:\n%s\n<END>", got, want, strings.ReplaceAll(got, " ", "\u2592"), strings.ReplaceAll(want, " ", "\u2592"))
}
