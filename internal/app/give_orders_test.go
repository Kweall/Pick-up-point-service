package app

import (
	"context"
	"fmt"
	"homework/internal/app/mocks"
	"homework/internal/storage/postgres"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestGiveOrders(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	mockStorage := mocks.NewFacadeMock(ctrl)

	tests := []struct {
		name          string
		mockSetup     func()
		parts         []string
		wantErr       bool
		expectedErr   error
		expectedCount int64
	}{
		{
			name: "success",
			mockSetup: func() {
				mockStorage.GetOrdersByIDsMock.Expect(context.Background(), []int64{1, 2}).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, ReceivedAt: nil},
					{OrderID: 2, ClientID: 1, ReceivedAt: nil},
				}, nil)
				mockStorage.GiveOrdersMock.Expect(context.Background(), []int64{1, 2}).Return(nil)
			},
			parts:         []string{"give_orders", "1", "2"},
			wantErr:       false,
			expectedCount: 2,
		},
		{
			name: "orderID is incorrect",
			mockSetup: func() {
			},
			parts:       []string{"give_orders", "invalid_orderID"},
			wantErr:     true,
			expectedErr: fmt.Errorf("orderID is incorrect"),
		},
		{
			name: "wrong count of arguments",
			mockSetup: func() {
			},
			parts:       []string{"give_orders"},
			wantErr:     true,
			expectedErr: fmt.Errorf("should be at least 1 argument: list of orderID's (int) separated by space"),
		},
		{
			name: "no orders found",
			mockSetup: func() {
				mockStorage.GetOrdersByIDsMock.Expect(context.Background(), []int64{1, 2}).Return(nil, fmt.Errorf("no orders found"))
			},
			parts:       []string{"give_orders", "1", "2"},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to get orders: %v", "no orders found"),
		},
		{
			name: "all orders must belong to the same client",
			mockSetup: func() {
				mockStorage.GetOrdersByIDsMock.Expect(context.Background(), []int64{1, 2}).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, ReceivedAt: nil},
					{OrderID: 2, ClientID: 2, ReceivedAt: nil},
				}, nil)
			},
			parts:       []string{"give_orders", "1", "2"},
			wantErr:     true,
			expectedErr: fmt.Errorf("all orders must belong to the same client"),
		},
		{
			name: "order has already been received",
			mockSetup: func() {
				mockStorage.GetOrdersByIDsMock.Expect(context.Background(), []int64{1, 2}).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, ReceivedAt: nil},
					{OrderID: 2, ClientID: 1, ReceivedAt: &time.Time{}},
				}, nil)
			},
			parts:       []string{"give_orders", "1", "2"},
			wantErr:     true,
			expectedErr: fmt.Errorf("order 2 has already been received"),
		},
		{
			name: "failed to write updated data to storage",
			mockSetup: func() {
				mockStorage.GetOrdersByIDsMock.Expect(context.Background(), []int64{1, 2}).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, ReceivedAt: nil},
					{OrderID: 2, ClientID: 1, ReceivedAt: nil},
				}, nil)
				mockStorage.GiveOrdersMock.Expect(context.Background(), []int64{1, 2}).Return(fmt.Errorf("failed to write updated data to storage"))
			},
			parts:       []string{"give_orders", "1", "2"},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to write updated data to storage"),
		},
	}

	service := &Service{storage: mockStorage}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := &GiveOrderRequest{}
			resp, err := service.GiveOrders(context.Background(), req, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("GiveOrders() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			} else if !tt.wantErr {
				assert.NotNil(t, resp, "expected a response but got nil")
				assert.Equal(t, resp.UpdatedCount, tt.expectedCount, "expected updated count to match")
			}
		})
	}
}
