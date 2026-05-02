package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/evolvedevlab/weavedeck/util"
)

var (
	Tags      = make(map[string]struct{})
	Stopwords = make(map[string]struct{})
)

func init() {
	var data struct {
		Tags      []string `json:"tags"`
		Stopwords []string `json:"stopwords"`
	}

	file, err := os.Open("config/data.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatal(err)
	}

	for _, t := range data.Tags {
		Tags[util.Normalize(t)] = struct{}{}
	}
	for _, s := range data.Stopwords {
		Stopwords[s] = struct{}{}
	}
}
