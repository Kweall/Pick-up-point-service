package main

import (
	"log"
)

func main() {
	conf := newConfig()

	err := runConsumer(conf)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
}
