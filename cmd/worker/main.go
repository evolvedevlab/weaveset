package main

import (
	"context"
	"fmt"
	"log"

	"github.com/evolvedevlab/weaveset/scraper"
)

func main() {
	sc, err := scraper.NewGRScraper("https://www.goodreads.com/list/show/399714.Goodreads_Editors_Picks_for_2026?ref=ls_fl_0_seeall")
	if err != nil {
		log.Fatal(err)
	}

	result, err := sc.Scrape(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
