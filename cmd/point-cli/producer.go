package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KafkaProducerInterface interface {
	SendMessage(topic, key, value string) error
	Close() error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "failed to connect to Kafka broker")
	}

	return &KafkaProducer{producer: producer}, nil
}

func (p *KafkaProducer) SendMessage(topic, key, value string) error {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(value),
		Timestamp: time.Now(),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return status.Error(codes.Internal, "failed to send message to Kafka")
	}
	log.Printf("Message is stored in partition %d, offset %d\n", partition, offset)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
