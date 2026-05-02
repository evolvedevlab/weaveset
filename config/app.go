package config

const (
	// redis stream
	Stream = "jobs"
	Group  = "workers"

	TagsKey = "tags"

	TriggerModifyFilename = ".changed"
	MaxJobRetryLimit      = 10
)
