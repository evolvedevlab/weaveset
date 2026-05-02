package util

import (
	"os"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func GetEnv(key string, fallback ...string) string {
	value := os.Getenv(key)
	if len(value) == 0 && len(fallback) > 0 {
		value = fallback[0]
	}
	return value
}

func GenerateSlug(s string, sep string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`[^a-z0-9 ]+`)
	s = re.ReplaceAllString(s, "")
	reSpaces := regexp.MustCompile(`\s+`)
	s = reSpaces.ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, " ", sep)
	return s
}

func Normalize(s string) string {
	// lowercase
	s = strings.ToLower(s)
	// remove possessive 's first
	s = strings.ReplaceAll(s, "'s", "")
	s = strings.ReplaceAll(s, "’s", "")
	// remove remaining apostrophes
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "’", "")
	// unicode normalize
	s = norm.NFD.String(s)
	// remove diacritics (accents)
	s = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, s)
	// normalize characters
	s = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			return r
		case unicode.IsSpace(r):
			return ' '
		// convert seperators to space
		case r == '-', r == '_', r == '/', r == '.', r == ',', r == ':', r == ';', r == '|', r == '+', r == '&':
			return ' '
		default:
			return -1
		}
	}, s)
	// collapse spaces
	s = strings.Join(strings.Fields(s), " ")

	return s
}
