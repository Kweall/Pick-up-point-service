package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	point_service "homework/pkg/point-service/v1"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
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

func sendAcceptReturnEvent(ctx context.Context, pointServiceClient point_service.PointServiceClient, prod KafkaProducerInterface, topic, eventType string, req *point_service.AcceptReturnRequest) (resp proto.Message, respErr error) {
	event := map[string]interface{}{
		"eventType": eventType,
		"clientId":  req.ClientId,
		"orderId":   req.OrderId,
		"timestamp": time.Now().Truncate(time.Minute).UnixNano(),
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %v", err)
	}
	err = prod.SendMessage(topic, eventType, string(eventData))
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %v", err)
	}
	resp, respErr = pointServiceClient.AcceptReturn(ctx, req)
	return resp, respErr
}

func sendGiveOrderEvent(ctx context.Context, pointServiceClient point_service.PointServiceClient, prod KafkaProducerInterface, topic, eventType string, req *point_service.GiveOrderRequest) (resp proto.Message, respErr error) {
	event := map[string]interface{}{
		"eventType": eventType,
		"orderIds":  req.OrderIds,
		"timestamp": time.Now().Truncate(time.Minute).UnixNano(),
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %v", err)
	}
	err = prod.SendMessage(topic, eventType, string(eventData))
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %v", err)
	}
	resp, respErr = pointServiceClient.GiveOrder(ctx, req)
	return resp, respErr
}

func sendAddOrderEvent(ctx context.Context, pointServiceClient point_service.PointServiceClient, prod KafkaProducerInterface, topic, eventType string, req *point_service.AddOrderRequest) (resp proto.Message, respErr error) {
	event := map[string]interface{}{
		"eventType":       eventType,
		"clientId":        req.ClientId,
		"orderId":         req.OrderId,
		"expiredAt":       req.ExpiredAt.AsTime().UnixNano(),
		"weight":          req.Weight,
		"price":           req.Price,
		"packaging":       req.Packaging,
		"additional_film": req.AdditionalFilm,
		"timestamp":       time.Now().Truncate(time.Minute).UnixNano(),
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %v", err)
	}
	err = prod.SendMessage(topic, eventType, string(eventData))
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %v", err)
	}
	resp, respErr = pointServiceClient.AddOrder(ctx, req)
	return resp, respErr
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
