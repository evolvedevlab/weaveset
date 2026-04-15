package scraper

import (
	"context"
	"net/http"

	"github.com/evolvedevlab/weaveset/data"
)

type Scraper interface {
	Scrape(context.Context) (*data.List, error)
}

func addRequestHeaders(r *http.Request) {
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/122.0.0.0 Safari/537.36")
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	r.Header.Set("Accept-Language", "en-US,en;q=0.9")
	r.Header.Set("Connection", "keep-alive")
}
