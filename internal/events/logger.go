package events

import (
	"context"
	"encoding/json"
	"fmt"
	point_service "homework/pkg/point-service/v1"
	"time"

	"google.golang.org/protobuf/proto"
)

type AcceptReturnEvent struct {
	ClientId  int64 `json:"clientId"`
	OrderId   int64 `json:"orderId"`
	Timestamp int64 `json:"timestamp"`
}

type GiveOrderEvent struct {
	OrderIds  []int64 `json:"orderIds"`
	Timestamp int64   `json:"timestamp"`
}

type AddOrderEvent struct {
	ClientId       int64   `json:"clientId"`
	OrderId        int64   `json:"orderId"`
	ExpiredAt      int64   `json:"expiredAt"`
	Weight         float32 `json:"weight"`
	Price          int64   `json:"price"`
	Packaging      string  `json:"packaging"`
	AdditionalFilm *bool   `json:"additional_film"`
	Timestamp      int64   `json:"timestamp"`
}

type KafkaProducer interface {
	SendMessage(topic string, key string, value string) error
	Close() error
}

type EventLogger struct {
	producer KafkaProducer
	topic    string
}

func NewEventLogger(producer KafkaProducer, topic string) *EventLogger {
	return &EventLogger{
		producer: producer,
		topic:    topic,
	}
}

func (e *EventLogger) LogAddOrderEvent(ctx context.Context, client point_service.PointServiceClient, req *point_service.AddOrderRequest) (resp proto.Message, respErr error) {
	resp, respErr = client.AddOrder(ctx, req)
	if respErr != nil {
		return nil, respErr
	}

	event := AddOrderEvent{
		ClientId:       req.ClientId,
		OrderId:        req.OrderId,
		ExpiredAt:      req.ExpiredAt.AsTime().UnixNano(),
		Weight:         req.Weight,
		Price:          req.Price,
		Packaging:      req.Packaging,
		AdditionalFilm: req.AdditionalFilm,
		Timestamp:      time.Now().Truncate(time.Minute).UnixNano(),
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	err = e.producer.SendMessage(e.topic, "AddOrder", string(eventData))
	if err != nil {
		return nil, fmt.Errorf("failed to send Kafka message: %w", err)
	}

	return resp, respErr
}

func (e *EventLogger) LogGiveOrderEvent(ctx context.Context, client point_service.PointServiceClient, req *point_service.GiveOrderRequest) (resp proto.Message, respErr error) {
	resp, respErr = client.GiveOrder(ctx, req)
	if respErr != nil {
		return nil, respErr
	}

	event := GiveOrderEvent{
		OrderIds:  req.OrderIds,
		Timestamp: time.Now().Truncate(time.Minute).UnixNano(),
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	err = e.producer.SendMessage(e.topic, "GiveOrder", string(eventData))
	if err != nil {
		return nil, fmt.Errorf("failed to send Kafka message: %w", err)
	}

	return resp, respErr
}

func (e *EventLogger) LogAcceptReturnEvent(ctx context.Context, client point_service.PointServiceClient, req *point_service.AcceptReturnRequest) (resp proto.Message, respErr error) {
	resp, respErr = client.AcceptReturn(ctx, req)
	if respErr != nil {
		return nil, respErr
	}

	event := AcceptReturnEvent{
		ClientId:  req.ClientId,
		OrderId:   req.OrderId,
		Timestamp: time.Now().Truncate(time.Minute).UnixNano(),
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	err = e.producer.SendMessage(e.topic, "AcceptReturn", string(eventData))
	if err != nil {
		return nil, fmt.Errorf("failed to send Kafka message: %w", err)
	}

	return resp, respErr
}
