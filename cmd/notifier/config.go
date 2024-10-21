package main

type kafkaConfig struct {
	Brokers   []string
	GroupID   string
	AutoReset string
}

type config struct {
	kafka kafkaConfig
}

func newConfig() config {
	return config{
		kafka: kafkaConfig{
			Brokers: []string{
				"kafka0:29092",
			},
			GroupID:   "notifier-group",
			AutoReset: "earliest",
		},
	}
}
