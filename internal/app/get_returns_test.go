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

func TestGetReturns(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)

	mockStorage := mocks.NewFacadeMock(ctrl)

	tests := []struct {
		name          string
		mockSetup     func()
		parts         []string
		wantErr       bool
		expectedErr   error
		expectedCount int
	}{
		{
			name: "success (page 1)",
			mockSetup: func() {
				mockStorage.GetReturnsMock.Expect(context.Background()).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, ReturnedAt: &time.Time{}},
					{OrderID: 2, ClientID: 1, ReturnedAt: &time.Time{}},
					{OrderID: 3, ClientID: 2, ReturnedAt: &time.Time{}},
					{OrderID: 4, ClientID: 2, ReturnedAt: &time.Time{}},
					{OrderID: 5, ClientID: 3, ReturnedAt: &time.Time{}},
				}, nil)
			},
			parts:         []string{"get_returns", "1"},
			wantErr:       false,
			expectedCount: 5,
		},
		{
			name: "success (page 2)",
			mockSetup: func() {
				mockStorage.GetReturnsMock.Expect(context.Background()).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, ReturnedAt: &time.Time{}},
					{OrderID: 2, ClientID: 1, ReturnedAt: &time.Time{}},
					{OrderID: 3, ClientID: 2, ReturnedAt: &time.Time{}},
					{OrderID: 4, ClientID: 2, ReturnedAt: &time.Time{}},
					{OrderID: 5, ClientID: 3, ReturnedAt: &time.Time{}},
					{OrderID: 6, ClientID: 3, ReturnedAt: &time.Time{}},
				}, nil)
			},
			parts:         []string{"get_returns", "2"},
			wantErr:       false,
			expectedCount: 1,
		},
		{
			name: "no returns found",
			mockSetup: func() {
				mockStorage.GetReturnsMock.Expect(context.Background()).Return([]*postgres.Order{}, nil)
			},
			parts:         []string{"get_returns", "1"},
			wantErr:       false,
			expectedCount: 0,
		},
		{
			name: "failed to get returns",
			mockSetup: func() {
				mockStorage.GetReturnsMock.Expect(context.Background()).Return(nil, fmt.Errorf("database error"))
			},
			parts:       []string{"get_returns", "1"},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to get returns: %v", "database error"),
		},
		{
			name: "invalid page number",
			mockSetup: func() {
			},
			parts:       []string{"get_returns", "invalid_page"},
			wantErr:     true,
			expectedErr: fmt.Errorf("invalid page number, must be greater than 0"),
		},
		{
			name: "page number exceeds available range",
			mockSetup: func() {
				mockStorage.GetReturnsMock.Expect(context.Background()).Return([]*postgres.Order{
					{OrderID: 1, ClientID: 1, ReturnedAt: &time.Time{}},
					{OrderID: 2, ClientID: 1, ReturnedAt: &time.Time{}},
					{OrderID: 3, ClientID: 2, ReturnedAt: &time.Time{}},
					{OrderID: 4, ClientID: 2, ReturnedAt: &time.Time{}},
					{OrderID: 5, ClientID: 3, ReturnedAt: &time.Time{}},
				}, nil)
			},
			parts:       []string{"get_returns", "2"},
			wantErr:     true,
			expectedErr: fmt.Errorf("page number exceeds the available range"),
		},
	}

	service := &Service{storage: mockStorage}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := &GetReturnsRequest{}
			resp, err := service.GetReturns(context.Background(), req, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetReturns() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			} else if !tt.wantErr {
				assert.NotNil(t, resp, "expected a response but got nil")
				assert.Equal(t, len(resp.Returns), tt.expectedCount, "expected number of returns to match")
			}
		})
	}
}
