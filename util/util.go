package util

import (
	"regexp"
	"strings"
)

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

// SelectLines returns a specific lines in the text
func SelectLines(str string, startLine int, endLine int) string {
	lines := strings.Split(str, "\n")

	startLine = clipNumber(startLine, 1, len(lines))
	endLine = clipNumber(endLine, 1, len(lines))

	if startLine > endLine {
		startLine, endLine = endLine, startLine
	}

	return strings.Join(lines[startLine-1:endLine], "\n")
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

// FormatMarkdownForSlack returns the formatted text of markdown for slack
func FormatMarkdownForSlack(str string) string {
	str = sanitizeEmbedImage(str)
	str = sanitizeEmptyLink(str)
	str = toSlackLink(str)
	return str
}

// ![text](url) -> text
func sanitizeEmbedImage(str string) string {
	re := regexp.MustCompile("\\!\\[(.*?)\\]\\(.*?\\)")
	return re.ReplaceAllString(str, "$1")
}

// [text](url) -> <url|text>
func toSlackLink(str string) string {
	re := regexp.MustCompile("\\[(\\S+?)\\]\\((\\S+?)\\)")
	return re.ReplaceAllString(str, "<$2|$1>")
}

func sanitizeEmptyLink(str string) string {
	// [](url) -> url
	re1 := regexp.MustCompile("\\[\\s*?\\]\\((\\S+?)\\)")
	str = re1.ReplaceAllString(str, "$1")

	// [text]() -> text
	re2 := regexp.MustCompile("\\[(\\S+?)\\]\\(\\s*?\\)")
	str = re2.ReplaceAllString(str, "$1")

	return str
}
