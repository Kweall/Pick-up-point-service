package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	events "homework/internal/events"
	point_service "homework/pkg/point-service/v1"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) SendMessage(topic, key, value string) error {
	args := m.Called(topic, key, value)
	return args.Error(0)
}

func (m *MockKafkaProducer) Close() error {
	return nil
}

type MockPointServiceClient struct {
	mock.Mock
}

func (m *MockPointServiceClient) AddOrder(ctx context.Context, req *point_service.AddOrderRequest, opts ...grpc.CallOption) (*point_service.AddOrderResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*point_service.AddOrderResponse), args.Error(1)
}

func (m *MockPointServiceClient) DeleteOrder(ctx context.Context, req *point_service.DeleteOrderRequest, opts ...grpc.CallOption) (*point_service.DeleteOrderResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*point_service.DeleteOrderResponse), args.Error(1)
}

func (m *MockPointServiceClient) GetOrders(ctx context.Context, req *point_service.GetOrdersRequest, opts ...grpc.CallOption) (*point_service.GetOrdersResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*point_service.GetOrdersResponse), args.Error(1)
}

func (m *MockPointServiceClient) GetReturns(ctx context.Context, req *point_service.GetReturnsRequest, opts ...grpc.CallOption) (*point_service.GetReturnsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*point_service.GetReturnsResponse), args.Error(1)
}

func (m *MockPointServiceClient) AcceptReturn(ctx context.Context, req *point_service.AcceptReturnRequest, opts ...grpc.CallOption) (*point_service.AcceptReturnResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*point_service.AcceptReturnResponse), args.Error(1)
}

func (m *MockPointServiceClient) GiveOrder(ctx context.Context, req *point_service.GiveOrderRequest, opts ...grpc.CallOption) (*point_service.GiveOrderResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*point_service.GiveOrderResponse), args.Error(1)
}

func TestLogAddOrderEvent(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	mockClient := new(MockPointServiceClient)
	ctx := context.Background()
	topic := "pvz.events-log"
	eventLogger := events.NewEventLogger(mockProducer, topic)
	additionalFilm := true
	req := &point_service.AddOrderRequest{
		ClientId:       1,
		OrderId:        1001,
		ExpiredAt:      timestamppb.New(time.Now().Add(24 * time.Hour)),
		Weight:         10,
		Price:          200,
		Packaging:      "box",
		AdditionalFilm: &additionalFilm,
	}

	event := events.AddOrderEvent{
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
	require.NoError(t, err)

	mockProducer.On("SendMessage", topic, "AddOrder", string(eventData)).Return(nil)
	mockClient.On("AddOrder", ctx, req).Return(&point_service.AddOrderResponse{}, nil)

	resp, err := eventLogger.LogAddOrderEvent(ctx, mockClient, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	mockProducer.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestLogGiveOrderEvent(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	mockClient := new(MockPointServiceClient)
	ctx := context.Background()
	topic := "pvz.events-log"
	eventLogger := events.NewEventLogger(mockProducer, topic)
	req := &point_service.GiveOrderRequest{
		OrderIds: []int64{1001},
	}

	event := events.GiveOrderEvent{
		OrderIds:  req.OrderIds,
		Timestamp: time.Now().Truncate(time.Minute).UnixNano(),
	}

	eventData, err := json.Marshal(event)
	require.NoError(t, err)

	mockProducer.On("SendMessage", topic, "GiveOrder", string(eventData)).Return(nil)
	mockClient.On("GiveOrder", ctx, req).Return(&point_service.GiveOrderResponse{}, nil)

	resp, err := eventLogger.LogGiveOrderEvent(ctx, mockClient, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	mockProducer.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestLogAcceptReturnEvent(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	mockClient := new(MockPointServiceClient)
	ctx := context.Background()
	topic := "pvz.events-log"
	eventLogger := events.NewEventLogger(mockProducer, topic)
	req := &point_service.AcceptReturnRequest{
		OrderId:  2001,
		ClientId: 1,
	}

	event := events.AcceptReturnEvent{
		ClientId:  req.ClientId,
		OrderId:   req.OrderId,
		Timestamp: time.Now().Truncate(time.Minute).UnixNano(),
	}

	eventData, err := json.Marshal(event)
	require.NoError(t, err)

	mockProducer.On("SendMessage", topic, "AcceptReturn", string(eventData)).Return(nil)
	mockClient.On("AcceptReturn", ctx, req).Return(&point_service.AcceptReturnResponse{}, nil)

	resp, err := eventLogger.LogAcceptReturnEvent(ctx, mockClient, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	mockProducer.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestNewKafkaProducer(t *testing.T) {
	brokers := []string{"localhost:9092"}
	producer, err := NewKafkaProducer(brokers)
	require.NoError(t, err)
	require.NotNil(t, producer)

	err = producer.Close()
	require.NoError(t, err)
}

func TestNewKafkaProducerFailure(t *testing.T) {
	brokers := []string{"invalid-broker"}
	producer, err := NewKafkaProducer(brokers)
	require.Error(t, err)
	require.Nil(t, producer)
}

func TestKafkaProducerSendMessage(t *testing.T) {
	brokers := []string{"localhost:9092"}
	producer, err := NewKafkaProducer(brokers)
	require.NoError(t, err)

	err = producer.SendMessage("test-topic", "test-key", "test-message")
	require.NoError(t, err)

	err = producer.Close()
	require.NoError(t, err)
}
