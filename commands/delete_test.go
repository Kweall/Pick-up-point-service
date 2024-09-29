package commands

import (
	"fmt"
	"homework/commands/mocks"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
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
				mockStorage.DeleteOrderByIDMock.Expect(int64(10)).Return(nil)
			},
			parts:   []string{"DELETE", "10"},
			wantErr: false,
		},
		{
			name: "orderID is incorrect",
			mockSetup: func() {
			},
			parts:       []string{"DELETE", "invalid_id"},
			wantErr:     true,
			expectedErr: fmt.Errorf("orderID is incorrect"),
		},
		{
			name: "order not found",
			mockSetup: func() {
				mockStorage.DeleteOrderByIDMock.Expect(int64(20)).Return(fmt.Errorf("order not found"))
			},
			parts:       []string{"DELETE", "20"},
			wantErr:     true,
			expectedErr: fmt.Errorf("order not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := Delete(mockStorage, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			}
		})
	}
}
