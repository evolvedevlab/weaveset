package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/evolvedevlab/weavedeck/config"
	"github.com/evolvedevlab/weavedeck/data"
	"github.com/evolvedevlab/weavedeck/util"
)

func main() {
	tags, _ := data.GenerateTags("best books of all time", config.Stopwords, config.Tags)
	fmt.Println(tags, len(tags))
}

func cleanTags() {
	file, err := os.Open("config/tags.txt")
	if err != nil {
		log.Fatal(err)
	}

	new, err := os.Create("new.txt")
	if err != nil {
		log.Fatal(err)
	}
	_ = new

	rd := bufio.NewReader(file)
	rd.ReadBytes('\n')

loop:
	for {
		data, err := rd.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		data = strings.Split(data, ",")[0]
		data = util.Normalize(data)

		if len(util.GenerateSlug(data, "")) > 0 {
			tokens := strings.Fields(data)

			for _, t := range tokens {
				if _, ok := config.Stopwords[t]; ok {
					continue loop
				}
			}
			if len(tokens) == 0 || len(tokens) > 3 {
				continue
			}
			if len(tokens) == 1 && len(tokens[0]) <= 3 {
				continue
			}
			if !hasNoDuplicates(tokens) {
				continue
			}
			if !hasValidContent(tokens) {
				continue
			}
			if isPureNumber(tokens) {
				continue
			}

			new.WriteString(data + "\n")
		}
	}
}

func hasNoDuplicates(tokens []string) bool {
	seen := map[string]struct{}{}
	for _, t := range tokens {
		if _, ok := seen[t]; ok {
			return false
		}
		seen[t] = struct{}{}
	}
	return true
}

func hasValidContent(tokens []string) bool {
	for _, t := range tokens {
		if len(t) >= 3 { // simple heuristic
			return true
		}
	}
	return false
}

func isPureNumber(tokens []string) bool {
	for _, t := range tokens {
		for _, r := range t {
			if !unicode.IsDigit(r) {
				return false
			}
		}
	}
	return true
}
