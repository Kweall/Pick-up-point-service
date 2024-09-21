package commands

import (
	"fmt"
	"homework/commands/mocks"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestGive(t *testing.T) {
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
				mockStorage.GiveOrdersToClientMock.Expect([]int64{1, 2}).Return(nil)
			},
			parts:   []string{"give", "1", "2"},
			wantErr: false,
		},
		{
			name: "orderID is incorrect",
			mockSetup: func() {
			},
			parts:       []string{"give", "invalid_orderID"},
			wantErr:     true,
			expectedErr: fmt.Errorf("orderID is incorrect"),
		},
		{
			name: "wrong count of arguments",
			mockSetup: func() {
			},
			parts:       []string{"give"},
			wantErr:     true,
			expectedErr: fmt.Errorf("should be at least 1 argument: list of orderID's (int) separated by space"),
		},
		{
			name: "order not found in storage",
			mockSetup: func() {
				mockStorage.GiveOrdersToClientMock.Expect([]int64{1, 2}).Return(fmt.Errorf("OrderID 2 not found"))
			},
			parts:       []string{"give", "1", "2"},
			wantErr:     true,
			expectedErr: fmt.Errorf("OrderID 2 not found"),
		},
		{
			name: "failed to write updated data to file",
			mockSetup: func() {
				mockStorage.GiveOrdersToClientMock.Expect([]int64{1, 2}).Return(fmt.Errorf("failed to write updated data to file"))
			},
			parts:       []string{"give", "1", "2"},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to write updated data to file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := Give(mockStorage, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Give() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			}
		})
	}
}
