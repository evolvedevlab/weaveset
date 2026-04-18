package util

import (
	"os"
	"regexp"
	"strings"
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
