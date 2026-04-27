package scraper

import "github.com/prometheus/client_golang/prometheus"

var (
	scrapeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "scrape_total",
			Help: "Total scrape attempts",
		},
		[]string{"status"}, // scrape success or fail
	)
	scrapeDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "scrape_duration_seconds",
			Help:    "Scrape duration",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(scrapeTotal, scrapeDuration)
}
