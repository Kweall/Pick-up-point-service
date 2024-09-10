package commands

import (
	"fmt"
	"homework/storage/json_file"
	"strconv"
)

const returnsPerPage = 5

func GetReturns(storage Storage, parts []string) error {

	orders, err := storage.GetAll()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	// Поиск всех возвратов
	var returns []json_file.Order
	for _, order := range orders {
		if !order.ReturnedAt.IsZero() {
			returns = append(returns, *order)
		}
	}

	if len(returns) == 0 {
		return fmt.Errorf("no returns found")
	}

	// Определение страницы
	page := 1
	if len(parts) > 1 {
		page, err = strconv.Atoi(parts[1])
		if err != nil || page <= 0 {
			fmt.Println("Invalid page number. Showing page 1.")
			page = 1
		}
	}

	// Пагинация: расчет начального и конечного индекса для вывода записей
	start := (page - 1) * returnsPerPage
	if start >= len(returns) {
		return fmt.Errorf("page number exceeds the available range")
	}
	end := start + returnsPerPage
	if end > len(returns) {
		end = len(returns)
	}

	// Вывод возвратов на текущей странице
	fmt.Printf("Showing returns, page %d:\n", page)
	for i := start; i < end; i++ {
		order := returns[i]
		fmt.Printf("Order ID: %d,\t Client ID: %d,\t Date of return: %s\t\n", order.ID, order.ClientID, order.ReturnedAt)
	}

	return nil
}
