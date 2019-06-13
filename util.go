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

// SelectLine returns a specific line in the text
func SelectLine(str string, line int) string {
	lines := strings.Split(str, "\n")

	line = clipNumber(line, 1, len(lines))

	return lines[line-1]
}

func clipNumber(number int, lower int, upper int) int {
	if number < lower {
		return lower
	}

	if number > upper {
		return upper
	}

	return number
}
