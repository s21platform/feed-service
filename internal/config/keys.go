package config

type key string

const (
	KeyUUID        = key("uuid")
	KeyMetrics key = key("metrics")
	KeyLogger      = key("logger")
)
