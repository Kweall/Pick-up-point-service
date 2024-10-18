package main

import (
	"time"
)

type appConfig struct {
	repeatCnt int
	startID   int
	count     int
	interval  time.Duration
}

type producerConfig struct {
	topic string
}

type config struct {
	app      appConfig
	kafka    kafkaConfig
	producer producerConfig
}

type kafkaConfig struct {
	Brokers []string
}

func newConfig(f flags) config {
	return config{
		app: appConfig{
			repeatCnt: f.repeatCnt,
			startID:   f.startID,
			count:     f.count,
			interval:  f.interval,
		},
		kafka: kafkaConfig{
			Brokers: []string{
				"localhost:9092",
			},
		},
		producer: producerConfig{
			topic: f.topic,
		},
	}
}
