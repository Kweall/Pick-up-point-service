package commands

import (
	"fmt"
	"homework/commands/mocks"
	"homework/storage/json_file"
	"os"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	const timeLayout = "02.01.2006"

	mc := minimock.NewController(t)

	tests := []struct {
		name string
		args struct {
			storage   Storage
			parts     []string
			userInput *os.File
		}
		wantErr     bool
		expectedErr error
		mock        func(m *mocks.StorageMock)
	}{
		{
			name: "success",
			args: struct {
				storage   Storage
				parts     []string
				userInput *os.File
			}{
				parts:   []string{"create", "1", "100", "15.12.2024"},
				storage: mocks.NewStorageMock(mc),
				userInput: func() *os.File {
					f, _ := os.CreateTemp("", "mock_input")
					f.WriteString("1.5\n500\nbox\nyes\n")
					f.Seek(0, 0)
					return f
				}(),
			},
			wantErr: false,
			mock: func(m *mocks.StorageMock) {
				parsedDate, _ := time.Parse(timeLayout, "15.12.2024")
				fixedCreatedAt := time.Now().Truncate(time.Minute) // Фиксированное время для теста
				order := &json_file.Order{
					ID:             100,
					ClientID:       1,
					CreatedAt:      fixedCreatedAt,
					ExpiredAt:      parsedDate,
					Weight:         1.5,
					Price:          500 + 20 + 1,
					Packaging:      "box",
					AdditionalFilm: "yes",
				}

				m.CheckIfExistsMock.Expect(int64(100)).Return(nil)
				m.AddOrderMock.Expect(order).Return(nil)
				m.AddOrderToStoryMock.Expect(order.ID, "storage/json_file/story_of_orders.json").Return(nil)
			},
		},
		{
			name: "error when not enough arguments",
			args: struct {
				storage   Storage
				parts     []string
				userInput *os.File
			}{
				parts:     []string{"create", "1", "100"},
				storage:   mocks.NewStorageMock(mc),
				userInput: nil,
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("should be 3 arguments: clientID (int), OrderID (int), Expired_date (dd.mm.yyyy)"),
			mock:        func(m *mocks.StorageMock) {},
		},
		{
			name: "error when order already exists",
			args: struct {
				storage   Storage
				parts     []string
				userInput *os.File
			}{
				parts:   []string{"create", "1", "100", "15.12.2024"},
				storage: mocks.NewStorageMock(mc),
				userInput: func() *os.File {
					f, _ := os.CreateTemp("", "mock_input")
					f.WriteString("1.5\n500\nbox\nyes\n")
					f.Seek(0, 0)
					return f
				}(),
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("this orderID already exists"),
			mock: func(m *mocks.StorageMock) {
				m.CheckIfExistsMock.Expect(int64(100)).Return(fmt.Errorf("this orderID already exists"))
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

			tt.mock(tt.args.storage.(*mocks.StorageMock))

			err := Create(tt.args.storage, tt.args.parts)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error(), "expected error message does not match")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
