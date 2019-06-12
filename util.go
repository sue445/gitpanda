package main

import "strings"

// TruncateWithLine truncate string with line
func TruncateWithLine(str string, maxLines int) string {
	if maxLines < 1 {
		return str
	}

	lines := strings.Split(str, "\n")

	if len(lines) > maxLines {
		lines = lines[:maxLines]
		return strings.Join(lines, "\n")
	}

	return str
}
