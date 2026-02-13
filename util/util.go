package util

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
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

// FormatMarkdownForSlack returns the formatted text of markdown for slack
func FormatMarkdownForSlack(str string) string {
	str = sanitizeEmbedImage(str)
	str = sanitizeEmptyLink(str)
	str = toSlackLink(str)
	return str
}

// ![text](url) -> text
func sanitizeEmbedImage(str string) string {
	return regexp.MustCompile(`\!\[(.*?)\]\(.*?\)`).ReplaceAllString(str, "$1")
}

// [text](url) -> <url|text>
func toSlackLink(str string) string {
	return regexp.MustCompile(`\[(\S+?)\]\((\S+?)\)`).ReplaceAllString(str, "<$2|$1>")
}

func sanitizeEmptyLink(str string) string {
	// [](url) -> url
	str = regexp.MustCompile(`\[\s*?\]\((\S+?)\)`).ReplaceAllString(str, "$1")

	// [text]() -> text
	str = regexp.MustCompile(`\[(\S+?)\]\(\s*?\)`).ReplaceAllString(str, "$1")

	return str
}

// WithDebugLogging executes blocks and outputs debug logs if necessary
func WithDebugLogging[T any](label string, isDebugLogging bool, fn func() (*T, error)) (*T, error) {
	start := time.Now()

	ret, err := fn()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if isDebugLogging {
		duration := time.Since(start)
		fmt.Printf("[DEBUG] %s : duration=%s, ret=%+v\n", label, duration, ret)
	}

	return ret, nil
}
