package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func runConsumer(conf config) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(conf.kafka.Brokers, conf.kafka.GroupID, config)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	handler := &consumerGroupHandler{}

	log.Println("Consumer started. Waiting for messages...")

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{"pvz.events-log"}, handler); err != nil {
				log.Printf("Error consuming: %v\n", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-sigchan
	log.Println("Caught signal: terminating...")

	return nil
}

type consumerGroupHandler struct{}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message received: key=%s, value=%s", string(msg.Key), string(msg.Value))

		err := processMessage(msg)
		if err != nil {
			log.Printf("Failed to process message: %v\n", err)
			continue
		}

		session.MarkMessage(msg, "")
		log.Printf("Offset %d committed for message with key: %s", msg.Offset, string(msg.Key))
	}

	return nil
}

func processMessage(msg *sarama.ConsumerMessage) error {
	log.Printf("Processing message with key: %s, value: %s", string(msg.Key), string(msg.Value))
	return nil
}
