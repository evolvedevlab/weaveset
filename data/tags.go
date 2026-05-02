package data

import (
	"strings"
	"unicode"

	"github.com/evolvedevlab/weavedeck/util"
)

func GenerateTags(s string, sw map[string]struct{}, presets map[string]struct{}) (tags []string, patterns []string) {
	var (
		tokens = tokenize(util.Normalize(s), sw)

		unigrams      = ngrams(tokens, 1)
		bigrams       = ngrams(tokens, 2)
		patternTokens = generatePatterns(tokens)
	)

	// combine all
	all := append(patternTokens, bigrams...)
	all = append(all, unigrams...)

	return match(all, presets), patternTokens
}

// matching against pre-defined tags
func match(tokens []string, presets map[string]struct{}) []string {
	var result []string

	var (
		n    = len(tokens)
		seen = make(map[string]struct{})
		used = make([]bool, n) // track token positions that's already included in one of the phrases
	)

	// Phrase matching
	for i := 0; i < n; i++ {
		for size := 3; size >= 2; size-- {
			if i+size > n {
				continue
			}

			candidate := strings.Join(tokens[i:i+size], " ")
			if _, ok := presets[candidate]; ok {
				if _, exists := seen[candidate]; !exists {
					seen[candidate] = struct{}{}
					result = append(result, candidate)
				}

				// mark tokens as used
				for j := i; j < i+size; j++ {
					used[j] = true
				}

				i += size - 1
				break
			}
		}
	}

	// Single token matching (only if not part of a phrase)
	for i, t := range tokens {
		if used[i] {
			continue
		}
		if _, ok := presets[t]; !ok {
			continue
		}
		if _, exists := seen[t]; exists {
			continue
		}

		seen[t] = struct{}{}
		result = append(result, t)
	}

	return result
}

func generatePatterns(tokens []string) []string {
	var result []string

	for i := 0; i < len(tokens); i++ {
		// eg. "20th century"
		if i < len(tokens)-1 {
			if isOrdinal(tokens[i]) && tokens[i+1] == "century" {
				result = append(result, tokens[i]+" "+tokens[i+1])
			}
		}
		// eg. "world war 2"
		if i < len(tokens)-2 {
			if tokens[i] == "world" && tokens[i+1] == "war" && isNumberLike(tokens[i+2]) {
				result = append(result, tokens[i]+" "+tokens[i+1]+" "+tokens[i+2])
			}
		}
	}

	return result
}

func ngrams(tokens []string, n int) []string {
	if n <= 0 || len(tokens) < n {
		return nil
	}

	var result []string

	for i := 0; i <= len(tokens)-n; i++ {
		gram := strings.Join(tokens[i:i+n], " ")
		result = append(result, gram)
	}

	return result
}

func tokenize(s string, sw map[string]struct{}) []string {
	raw := strings.Fields(s)

	var tokens []string
	for _, t := range raw {
		// remove stopwords
		if _, ok := sw[t]; ok {
			continue
		}
		// remove short tokens
		if len(t) < 3 {
			continue
		}
		// skip pure numbers
		if isPureNumber(t) {
			continue
		}

		tokens = append(tokens, t)
	}

	return tokens
}

// matches 1st, 2nd, 3rd, 4th...
func isOrdinal(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return strings.HasSuffix(s, "st") ||
		strings.HasSuffix(s, "nd") ||
		strings.HasSuffix(s, "rd") ||
		strings.HasSuffix(s, "th")
}

func isPureNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isNumberLike(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}
