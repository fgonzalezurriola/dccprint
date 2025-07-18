package scripts

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// EscapeFilename cleans the filename by removing diacritics, spaces, and disallowed characters.
func EscapeFilename(name string) string {
	name = strings.ToLower(name)

	// Normalize to NFD and remove diacritics
	t := norm.NFD.String(name)
	normRunes := make([]rune, 0, len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		normRunes = append(normRunes, r)
	}
	name = string(normRunes)

	name = strings.ReplaceAll(name, " ", "")

	// Keep only lowercase letters, numbers, dot, underscore, and hyphen
	re := regexp.MustCompile(`[^a-z0-9._-]`)
	name = re.ReplaceAllString(name, "")

	if name == "" {
		return "file.pdf"
	}
	return name
}
