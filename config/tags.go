package config

import (
	"bufio"
	"embed"
	"io"
	"log"
	"os"
	"strings"

	"github.com/evolvedevlab/weavedeck/util"
)

//go:embed *
var embedFS embed.FS

var (
	Tags      = make(map[string]struct{})
	Stopwords = make(map[string]struct{})
)

func init() {
	if !strings.HasSuffix(os.Args[0], ".test") {
		sw, err := embedFS.Open("stopwords.txt")
		if err != nil {
			log.Fatal(err)
		}
		t, err := embedFS.Open("tags.txt")
		if err != nil {
			log.Fatal(err)
		}

		rd := bufio.NewReader(sw)
		rd.ReadString('\n')
		for {
			s, err := rd.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			Stopwords[strings.TrimSpace(s)] = struct{}{}
		}

		rd = bufio.NewReader(t)
		rd.ReadString('\n')
		for {
			s, err := rd.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			Tags[util.Normalize(s)] = struct{}{}
		}
	}
}
