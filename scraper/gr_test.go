package scraper

import (
	"fmt"
	"testing"
)

func TestGRScrape(t *testing.T) {
	sc, err := NewGRScraper("https://www.goodreads.com/list/show/43")
	if err != nil {
		t.Error(err)
	}

	result, err := sc.Scrape(t.Context())
	if err != nil {
		t.Error(err)
	}

	_ = result
	for i, item := range result.Items {
		fmt.Printf("%+v\n", item)
		if i == 10 {
			break
		}
	}
}
