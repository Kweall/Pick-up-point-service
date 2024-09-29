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

func TestGetOrders(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	mockStorage := mocks.NewFacadeMock(ctrl)

	tests := []struct {
		name        string
		mockSetup   func()
		parts       []string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success with no limit",
			mockSetup: func() {
				mockStorage.GetOrdersMock.Expect(context.Background(), int64(1)).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, CreatedAt: &time.Time{}},
					{OrderID: 2, ClientID: 1, CreatedAt: &time.Time{}},
				}, nil)
			},
			parts:   []string{"get_orders", "1"},
			wantErr: false,
		},
		{
			name: "success with limit",
			mockSetup: func() {
				mockStorage.GetOrdersMock.Expect(context.Background(), int64(1)).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, CreatedAt: &time.Time{}},
					{OrderID: 2, ClientID: 1, CreatedAt: &time.Time{}},
					{OrderID: 3, ClientID: 2, CreatedAt: &time.Time{}},
				}, nil)
			},
			parts:   []string{"get_orders", "1", "2"},
			wantErr: false,
		},
		{
			name: "no orders found for clientID",
			mockSetup: func() {
				mockStorage.GetOrdersMock.Expect(context.Background(), int64(1)).Return([]*postgres.Order{}, nil)
			},
			parts:       []string{"get_orders", "1"},
			wantErr:     true,
			expectedErr: fmt.Errorf("no orders found for clientID: %v", 1),
		},
		{
			name: "failed to get orders",
			mockSetup: func() {
				mockStorage.GetOrdersMock.Expect(context.Background(), int64(1)).Return(nil, fmt.Errorf("database error"))
			},
			parts:       []string{"get_orders", "1"},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to get orders: %v", "database error"),
		},
		{
			name: "clientID is incorrect",
			mockSetup: func() {
			},
			parts:       []string{"get_orders", "invalid_client_id"},
			wantErr:     true,
			expectedErr: fmt.Errorf("clientID is incorrect"),
		},
		{
			name: "limit is incorrect",
			mockSetup: func() {
			},
			parts:       []string{"get_orders", "1", "invalid_limit"},
			wantErr:     true,
			expectedErr: fmt.Errorf("limit is incorrect"),
		},
		{
			name: "wrong count of arguments",
			mockSetup: func() {
			},
			parts:       []string{"get_orders", "1", "2", "3"},
			wantErr:     true,
			expectedErr: fmt.Errorf("should be maximum 2 arguments: clientID (int) and count of orders you want to get"),
		},
	}

	service := &Service{storage: mockStorage}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := &GetOrdersRequest{}
			resp, err := service.GetOrders(context.Background(), req, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrders() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			} else if !tt.wantErr {
				assert.NotNil(t, resp, "expected a response but got nil")
			}
		})
	}
}
