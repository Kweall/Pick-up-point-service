package commands

import (
	"fmt"
	"homework/commands/mocks"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAcceptReturn(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	mockStorage := mocks.NewStorageMock(ctrl)

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
				mockStorage.AcceptReturnMock.Expect(int64(1), int64(100)).Return(nil)
			},
			parts:   []string{"ACCEPT_RETURN", "1", "100"},
			wantErr: false,
		},
		{
			name: "order not found or does not belong to the given client",
			mockSetup: func() {
				mockStorage.AcceptReturnMock.Expect(int64(1), int64(100)).Return(fmt.Errorf("order not found or does not belong to the given client"))
			},
			parts:       []string{"ACCEPT_RETURN", "1", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("order not found or does not belong to the given client"),
		},
		{
			name: "order has not been given",
			mockSetup: func() {
				mockStorage.AcceptReturnMock.Expect(int64(1), int64(100)).Return(fmt.Errorf("order has not been given"))
			},
			parts:       []string{"ACCEPT_RETURN", "1", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("order has not been given"),
		},
		{
			name: "order has already been returned",
			mockSetup: func() {
				mockStorage.AcceptReturnMock.Expect(int64(1), int64(100)).Return(fmt.Errorf("order has already been returned"))
			},
			parts:       []string{"ACCEPT_RETURN", "1", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("order has already been returned"),
		},
		{
			name: "can't return, return period has expired",
			mockSetup: func() {
				mockStorage.AcceptReturnMock.Expect(int64(1), int64(100)).Return(fmt.Errorf("can't return, return period has expired"))
			},
			parts:       []string{"ACCEPT_RETURN", "1", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("can't return, return period has expired"),
		},
		{
			name: "clientID is incorrect",
			mockSetup: func() {
			},
			parts:       []string{"ACCEPT_RETURN", "invalid_client_id", "100"},
			wantErr:     true,
			expectedErr: fmt.Errorf("clientID is incorrect"),
		},
		{
			name: "orderID is incorrect",
			mockSetup: func() {
			},
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

			err := AcceptReturn(mockStorage, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("AcceptReturn() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			}
		})
	}
}
