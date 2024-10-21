package main

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestKafkaIntegration(t *testing.T) {
	brokers := []string{"kafka0:29092"}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	assert.NoError(t, err)
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: "pvz.events-log",
		Value: sarama.StringEncoder(`{"eventType":"AddOrder","data":""additional_film":true,"clientId":1,"eventType":"AddOrder","expiredAt":1734220800000000000,"orderId":100000,"packaging":"box","price":100,"timestamp":1729523760000000000,"weight":1"}`),
	}

	_, _, err = producer.SendMessage(message)
	assert.NoError(t, err)

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumer(brokers, consumerConfig)
	assert.NoError(t, err)
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("pvz.events-log", 0, sarama.OffsetOldest)
	assert.NoError(t, err)
	defer partitionConsumer.Close()

	select {
	case msg := <-partitionConsumer.Messages():
		assert.Equal(t, message.Value, msg.Value)
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for message")
	}
}
