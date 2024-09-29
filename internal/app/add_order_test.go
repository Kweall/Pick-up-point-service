package app

import (
	"context"
	"homework/internal/app/mocks"
	"os"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	t.Parallel()
	const timeLayout = "02.01.2006"

	mc := minimock.NewController(t)

	tests := []struct {
		name string
		args struct {
			storage   Facade
			parts     []string
			userInput *os.File
		}
		wantErr     bool
		expectedErr error
		mock        func(m *mocks.FacadeMock)
	}{
		{
			name: "success",
			args: struct {
				storage   Facade
				parts     []string
				userInput *os.File
			}{
				parts:   []string{"AddOrder", "1", "100", "15.12.2024"},
				storage: mocks.NewFacadeMock(mc),
				userInput: func() *os.File {
					f, _ := os.CreateTemp("", "mock_input")
					f.WriteString("1.5\n500\nbox\nyes\n")
					f.Seek(0, 0)
					return f
				}(),
			},
			wantErr: false,
			mock: func(m *mocks.FacadeMock) {
				parsedDate, _ := time.Parse(timeLayout, "15.12.2024")
				fixedCreatedAt := time.Now().Truncate(time.Minute)
				orderRequest := &AddOrderRequest{
					ClientID:       1,
					OrderID:        100,
					CreatedAt:      fixedCreatedAt,
					ExpiredAt:      parsedDate,
					Weight:         1.5,
					Price:          500,
					Packaging:      "box",
					AdditionalFilm: "yes",
				}

				m.AddOrderMock.Expect(context.Background(), orderRequest.OrderID, orderRequest.ClientID, orderRequest.CreatedAt, orderRequest.ExpiredAt, orderRequest.Weight, orderRequest.Price, orderRequest.Packaging, orderRequest.AdditionalFilm).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.userInput != nil {
				defer tt.args.userInput.Close()
				oldStdin := os.Stdin
				defer func() { os.Stdin = oldStdin }()
				os.Stdin = tt.args.userInput
			}

			tt.mock(tt.args.storage.(*mocks.FacadeMock))

			if !tt.wantErr {
				parsedDate, _ := time.Parse(timeLayout, "15.12.2024")
				req := &AddOrderRequest{
					ClientID:       1,
					OrderID:        100,
					CreatedAt:      time.Now().Truncate(time.Minute),
					ExpiredAt:      parsedDate,
					Weight:         1.5,
					Price:          500,
					Packaging:      "box",
					AdditionalFilm: "yes",
				}

				// Вызов метода AddOrder
				err := tt.args.storage.AddOrder(context.Background(), req.OrderID, req.ClientID, req.CreatedAt, req.ExpiredAt, req.Weight, req.Price, req.Packaging, req.AdditionalFilm)
				assert.NoError(t, err)
			} else {
				// Не ожидаем вызов метода AddOrder
				assert.Panics(t, func() {
					_ = tt.args.storage.AddOrder(context.Background(), 0, 0, time.Time{}, time.Time{}, 0, 0, "", "")
				}, "expected to panic due to nil request")
			}
		})
	}
}
