package commands

import (
	"fmt"
	"homework/storage/json_file"
	"sort"
	"strconv"
)

const countOfArgumentsToGetOrders = 2

func GetOrders(storage Storage, parts []string) error {
	orders, err := storage.GetAll()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	if len(parts) < countOfArgumentsToGetOrders-1 {
		return fmt.Errorf("should be at least 1 argument: clientID (int)")
	} else if len(parts) > countOfArgumentsToGetOrders {
		return fmt.Errorf("should be maximum 2 arguments: clientID (int) and count of orders you want to get")
	}

	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("clientID is incorrect")
	}

	// Определяем значение limit, если оно указано
	limit := -1 // -1 будет означать, что ограничение не задано
	if len(parts) == countOfArgumentsToGetOrders {
		limit, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("limit is incorrect")
		}
	}

	// Поиск заказов пользователя
	var ordersSlice []*json_file.Order
	for _, order := range orders {
		if order.ClientID == int64(clientID) {
			ordersSlice = append(ordersSlice, order)
		}
	}

	if len(ordersSlice) == 0 {
		return fmt.Errorf("no orders found for clientID: %v", clientID)
	}

	// Сортируем заказы по дате (начиная с нового)
	sort.Slice(ordersSlice, func(i, j int) bool {
		date1 := ordersSlice[i].CreatedAt
		date2 := ordersSlice[j].CreatedAt
		return date1.After(date2)
	})

	// Если установлен limit и он больше 0, ограничиваем вывод
	if limit > 0 && limit < len(ordersSlice) {
		ordersSlice = ordersSlice[:limit]
	}

	// Вывод заказов
	for _, order := range ordersSlice {
		fmt.Printf("Order ID: %d,\t Client ID: %d,\t Created at: %s,\t Expired at: %s\n", order.ID, order.ClientID, order.CreatedAt, order.ExpiredAt)
	}

	return nil
}
