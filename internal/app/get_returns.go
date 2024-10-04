package app

import (
	"context"
	"fmt"
	"homework/internal/storage/postgres"
	"strconv"
)

type GetReturnsRequest struct {
	Page int `json:"Page"` // Поле для номера страницы
}

type GetReturnsResponse struct {
	Returns []*postgres.Order `json:"Returns"`
}

const (
	returnsPerPage               = 5
	countOfArgumentsToGetReturns = 2
)

func (s *Service) GetReturns(ctx context.Context, req *GetReturnsRequest, parts []string) (*GetReturnsResponse, error) {
	if len(parts) < countOfArgumentsToGetReturns-1 || len(parts) > countOfArgumentsToGetReturns {
		return nil, fmt.Errorf("should be 1 or 2 arguments: page number (optional)")
	}

	page := 1 // По умолчанию первая страница
	if len(parts) == 2 {
		var err error
		page, err = strconv.Atoi(parts[1])
		if err != nil || page <= 0 {
			return nil, fmt.Errorf("invalid page number, must be greater than 0")
		}
	}

	returns, err := s.storage.GetReturns(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get returns: %v", err)
	}

	if len(returns) == 0 {
		return &GetReturnsResponse{Returns: nil}, nil // Возвращаем пустой список, если возвратов нет
	}

	// Пагинация
	start := (page - 1) * returnsPerPage
	if start >= len(returns) {
		return nil, fmt.Errorf("page number exceeds the available range")
	}

	end := start + returnsPerPage
	if end > len(returns) {
		end = len(returns)
	}

	fmt.Printf("Showing returns, page %d:\n", page)
	for i := start; i < end; i++ {
		order := returns[i]
		fmt.Printf("Order ID: %d,\t Client ID: %d,\t Date of return: %s\t\n", order.OrderID, order.ClientID, order.ReturnedAt)
	}

	return &GetReturnsResponse{Returns: returns[start:end]}, nil
}
