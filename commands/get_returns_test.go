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

func TestGetReturns(t *testing.T) {
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
			name: "success (page 1)",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(map[int64]*json_file.Order{
					1: {ID: 1, ClientID: 1, ReturnedAt: time.Now().Add(-24 * time.Hour)},
					2: {ID: 2, ClientID: 1, ReturnedAt: time.Now().Add(-48 * time.Hour)},
				}, nil)
			},
			parts:   []string{"get_returns", "1"},
			wantErr: false,
		},
		{
			name: "success (page 2)",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(map[int64]*json_file.Order{
					1: {ID: 1, ClientID: 1, ReturnedAt: time.Now().Add(-24 * time.Hour)},
					2: {ID: 2, ClientID: 1, ReturnedAt: time.Now().Add(-48 * time.Hour)},
					3: {ID: 3, ClientID: 1, ReturnedAt: time.Now().Add(-72 * time.Hour)},
					4: {ID: 4, ClientID: 2, ReturnedAt: time.Now().Add(-96 * time.Hour)},
					5: {ID: 5, ClientID: 2, ReturnedAt: time.Now().Add(-120 * time.Hour)},
					6: {ID: 6, ClientID: 2, ReturnedAt: time.Now().Add(-144 * time.Hour)},
				}, nil)
			},
			parts:   []string{"get_returns", "2"},
			wantErr: false,
		},
		{
			name: "no returns found",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(map[int64]*json_file.Order{
					1: {ID: 1, ClientID: 1, ReturnedAt: time.Time{}}, // No returned orders
				}, nil)
			},
			parts:       []string{"get_returns", "1"},
			wantErr:     true,
			expectedErr: fmt.Errorf("no returns found"),
		},
		{
			name: "failed to read from file",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(nil, fmt.Errorf("read error"))
			},
			parts:       []string{"get_returns", "1"},
			wantErr:     true,
			expectedErr: fmt.Errorf("failed to read from file: %v", "read error"),
		},
		{
			name: "invalid page number",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(map[int64]*json_file.Order{
					1: {ID: 1, ClientID: 1, ReturnedAt: time.Now()},
					2: {ID: 2, ClientID: 1, ReturnedAt: time.Now()},
					3: {ID: 3, ClientID: 2, ReturnedAt: time.Now()},
				}, nil)
			},
			parts:   []string{"get_returns", "invalid_page"},
			wantErr: false,
		},

		{
			name: "page number exceeds available range but shows page 1",
			mockSetup: func() {
				mockStorage.GetAllMock.Expect().Return(map[int64]*json_file.Order{
					1: {ID: 1, ClientID: 1, ReturnedAt: time.Now()},
					2: {ID: 2, ClientID: 1, ReturnedAt: time.Now()},
					3: {ID: 3, ClientID: 2, ReturnedAt: time.Now()},
				}, nil)
			},
			parts:       []string{"get_returns", "10"},
			wantErr:     true,
			expectedErr: fmt.Errorf("page number exceeds the available range"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := GetReturns(mockStorage, tt.parts)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetReturns() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error does not match")
			}
		})
	}
}
