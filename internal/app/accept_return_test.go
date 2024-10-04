package app

import (
	"context"
	"fmt"
	"homework/internal/app/mocks"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAcceptReturn(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	mockStorage := mocks.NewFacadeMock(ctrl)

	service := &Service{
		storage: mockStorage,
	}

	tests := []struct {
		name        string
		mockSetup   func()
		parts       []string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success",
			mockSetup: func() {
				mockStorage.CheckOrderStatusMock.Expect(context.Background(), int64(100)).Return(true, false, nil)
				mockStorage.AcceptReturnMock.Expect(context.Background(), int64(1), int64(100)).Return(nil)
			},
			parts:   []string{"ACCEPT_RETURN", "1", "100"},
			wantErr: false,
		},
		{
			name: "order has not been received yet",
			mockSetup: func() {
				mockStorage.CheckOrderStatusMock.Expect(context.Background(), int64(100)).Return(false, false, nil)
			},
			parts:       []string{"ACCEPT_RETURN", "1", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("order 100 has not been received yet"),
		},
		{
			name: "order has already been returned",
			mockSetup: func() {
				mockStorage.CheckOrderStatusMock.Expect(context.Background(), int64(100)).Return(true, true, nil)
			},
			parts:       []string{"ACCEPT_RETURN", "1", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("order 100 has already been returned"),
		},
		{
			name: "failed to check order status",
			mockSetup: func() {
				mockStorage.CheckOrderStatusMock.Expect(context.Background(), int64(100)).Return(false, false, fmt.Errorf("database error"))
			},
			parts:       []string{"ACCEPT_RETURN", "1", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to check order status: database error"),
		},
		{
			name: "error during return acceptance",
			mockSetup: func() {
				mockStorage.CheckOrderStatusMock.Expect(context.Background(), int64(100)).Return(true, false, nil)
				mockStorage.AcceptReturnMock.Expect(context.Background(), int64(1), int64(100)).Return(fmt.Errorf("return acceptance failed"))
			},
			parts:       []string{"ACCEPT_RETURN", "1", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("return acceptance failed"),
		},
		{
			name:        "clientID is incorrect",
			mockSetup:   func() {},
			parts:       []string{"ACCEPT_RETURN", "invalid_client_id", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("clientID is incorrect"),
		},
		{
			name:        "orderID is incorrect",
			mockSetup:   func() {},
			parts:       []string{"ACCEPT_RETURN", "1", "invalid_order_id"},
			wantErr:     true,
			expectedErr: fmt.Errorf("orderID is incorrect"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := &AcceptReturnRequest{}
			_, err := service.AcceptReturn(context.Background(), req, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("AcceptReturn() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			}
		})
	}
}
