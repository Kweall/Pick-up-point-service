package commands

import (
	"fmt"
	"homework/storage/json_file"
	"sort"
	"strconv"
	"time"
)

func GetOrders(parts []string) error {
	storage, err := json_file.NewStorage("storage/json_file/data.json")
	if err != nil {
		return fmt.Errorf("can't init storage: %v", err)
	}

	data, err := storage.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	if len(parts) < 2 {
		fmt.Println("Should be at least 1 argument: userID (int)")
		return nil
	} else if len(parts) > 3 {
		fmt.Println("Should be maximum 2 arguments: userID (int) and count of orders you want to get")
		return nil
	}

	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("userID isn't correct")
		return err
	}

	// Определяем значение limit, если оно указано
	limit := -1 // -1 будет означать, что ограничение не задано
	if len(parts) == 3 {
		limit, err = strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("limit isn't correct")
			return err
		}
	}

	// Поиск заказов пользователя
	var orders []json_file.Item
	for _, user := range data.Users {
		if user.UserID == int64(userID) {
			orders = user.Items
			break
		}
	}

	if len(orders) == 0 {
		fmt.Printf("No orders found for userID: %v\n", userID)
		return nil
	}

	// Сортируем заказы по дате (начиная с нового)
	sort.Slice(orders, func(i, j int) bool {
		date1, _ := time.Parse("02.01.2006", orders[i].Date)
		date2, _ := time.Parse("02.01.2006", orders[j].Date)
		return date1.After(date2)
	})

	// Если установлен limit и он больше 0, ограничиваем вывод
	if limit > 0 && limit < len(orders) {
		orders = orders[:limit]
	}

	// Вывод заказов
	for _, order := range orders {
		fmt.Printf("Order ID: %d, Date: %s, Valid: %t\n", order.ID, order.Date, order.Valid)
	}

	return nil
}
