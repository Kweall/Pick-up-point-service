package main

import (
	"time"

	flag "github.com/spf13/pflag"
)

type flags struct {
	repeatCnt int
	startID   int
	count     int
	topic     string
	interval  time.Duration
}

var cliFlags = flags{}

func init() {
	flag.IntVar(&cliFlags.repeatCnt, "repeat-count", 3, "count times all messages sent")
	flag.IntVar(&cliFlags.startID, "start-id", 1, "start order-id of all messages")
	flag.IntVar(&cliFlags.count, "count", 10, "count of orders to emit events")
	flag.StringVar(&cliFlags.topic, "topic", "pvz.events-log", "topic to produce")
	flag.DurationVar(&cliFlags.interval, "interval", 100*time.Millisecond, "duration between messages")

	flag.Parse()
}
