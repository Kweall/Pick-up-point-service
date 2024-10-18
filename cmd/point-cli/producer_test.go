package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestSendOrderEvent(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	topic := "pvz.events-log"
	eventType := "AddOrder"
	id := 12345

	event := map[string]interface{}{
		"eventType": eventType,
		"id":        id,
	}
	_, err := json.Marshal(event)
	require.NoError(t, err)

	mockProducer.On("SendMessage", topic, eventType, mock.Anything).Return(nil)

	sendOrderEvent(mockProducer, topic, eventType, id)

	mockProducer.AssertExpectations(t)
}
