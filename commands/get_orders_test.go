package commands

import (
	"fmt"
	"homework/commands/mocks"
	"homework/storage/json_file"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetOrders(t *testing.T) {
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
			name: "success with no limit",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(map[int64]*json_file.Order{
					1: {ID: 1, ClientID: 1, CreatedAt: time.Now(), ExpiredAt: time.Now().Add(24 * time.Hour)},
					2: {ID: 2, ClientID: 1, CreatedAt: time.Now().Add(-1 * time.Hour), ExpiredAt: time.Now().Add(24 * time.Hour)},
				}, nil)
			},
			parts:   []string{"get_orders", "1"},
			wantErr: false,
		},
		{
			name: "success with limit",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(map[int64]*json_file.Order{
					1: {ID: 1, ClientID: 1, CreatedAt: time.Now(), ExpiredAt: time.Now().Add(24 * time.Hour)},
					2: {ID: 2, ClientID: 1, CreatedAt: time.Now().Add(-1 * time.Hour), ExpiredAt: time.Now().Add(24 * time.Hour)},
					3: {ID: 3, ClientID: 2, CreatedAt: time.Now().Add(-2 * time.Hour), ExpiredAt: time.Now().Add(24 * time.Hour)},
				}, nil)
			},
			parts:   []string{"get_orders", "1", "1"},
			wantErr: false,
		},
		{
			name: "no orders found for clientID",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(map[int64]*json_file.Order{}, nil)
			},
			parts:       []string{"get_orders", "1"},
			wantErr:     true,
			expectedErr: fmt.Errorf("no orders found for clientID: %v", 1),
		},
		{
			name: "failed to read from file",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(nil, fmt.Errorf("read error"))
			},
			parts:       []string{"get_orders", "1"},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to read from file: %v", "read error"),
		},
		{
			name: "clientID is incorrect",
			mockSetup: func() {
				// Парсинг числа упадет до вызова метода GetAll
			},
			parts:       []string{"get_orders", "invalid_client_id"},
			wantErr:     true,
			expectedErr: fmt.Errorf("clientID is incorrect"),
		},
		{
			name: "limit is incorrect",
			mockSetup: func() {
				// Парсинг числа упадет до вызова метода GetAll
			},
			parts:       []string{"get_orders", "1", "invalid_limit"},
			wantErr:     true,
			expectedErr: fmt.Errorf("limit is incorrect"),
		},
		{
			name: "wrong count of arguments",
			mockSetup: func() {
				// Не требуется мок, так как проверка аргументов произойдет до вызова GetAll
			},
			parts:       []string{"get_orders", "1", "2", "3"},
			wantErr:     true,
			expectedErr: fmt.Errorf("should be maximum 2 arguments: clientID (int) and count of orders you want to get"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := GetOrders(mockStorage, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrders() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			}
		})
	}
}
