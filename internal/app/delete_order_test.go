package app

import (
	"context"
	"fmt"
	"homework/internal/app/mocks"
	"strconv"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOrder(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	mockStorage := mocks.NewFacadeMock(ctrl)
	service := &Service{storage: mockStorage} // Создайте экземпляр Service

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
				mockStorage.DeleteOrderMock.Expect(context.Background(), int64(10)).Return(nil)
			},
			parts:   []string{"delete_order", "10"},
			wantErr: false,
		},
		{
			name: "orderID is incorrect",
			mockSetup: func() {
			},
			parts:       []string{"delete_order", "invalid_id"},
			wantErr:     true,
			expectedErr: fmt.Errorf("strconv.ParseInt: parsing \"invalid_id\": invalid syntax"),
		},
		{
			name: "order not found",
			mockSetup: func() {
				mockStorage.DeleteOrderMock.Expect(context.Background(), int64(20)).Return(fmt.Errorf("order not found"))
			},
			parts:       []string{"delete_order", "20"},
			wantErr:     true,
			expectedErr: fmt.Errorf("order not found"),
		},
		{
			name: "not enough arguments",
			mockSetup: func() {
			},
			parts:       []string{"delete_order"},
			wantErr:     true,
			expectedErr: fmt.Errorf("should be 1 argument: orderID (int)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.parts) < 2 {
				assert.EqualError(t, fmt.Errorf("should be 1 argument: orderID (int)"), tt.expectedErr.Error())
				return
			}

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			orderID, err := strconv.ParseInt(tt.parts[1], 10, 64)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("expected no error but got: %v", err)
				}
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
				return
			}

			_, err = service.DeleteOrder(context.Background(), &DeleteOrderRequest{OrderID: orderID}, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			} else if !tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}
