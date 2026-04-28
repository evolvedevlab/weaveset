package queue

import "github.com/prometheus/client_golang/prometheus"

var (
	jobsProcessedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jobs_processed_total",
			Help: "Total number of jobs processed",
		},
		[]string{"status"}, // job success or fail
	)
	jobDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "job_duration_seconds",
			Help:    "Time taken to process a job",
			Buckets: prometheus.DefBuckets,
		},
	)
	jobRetriesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_retries_total",
			Help: "Total number of job retries",
		},
	)
	jobsDroppedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_dropped_total",
			Help: "Total number of jobs dropped after exceeding retry limit",
		},
	)
)

func init() {
	prometheus.MustRegister(
		jobsProcessedTotal,
		jobDuration,
		jobRetriesTotal,
		jobsDroppedTotal,
	)
}
