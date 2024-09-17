package commands

import (
	"fmt"
	"homework/storage/json_file"
	"strconv"
	"time"
)

const (
	countOfArgumentsToCreate = 3
	timeLayout               = "02.01.2006"
)

func Create(storage Storage, parts []string) error {

	if len(parts) != countOfArgumentsToCreate {
		return fmt.Errorf("should be 3 arguments: clientID (int), OrderID (int), Expired_date (dd.mm.yyyy)")
	}

	// Получение и преобразование аргументов
	clientID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("clientID is incorrect")
	}

	orderID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("orderID is incorrect")
	}

	// Проверяем, принимался ли этот заказ ранее
	err = storage.Validation(orderID)
	if err != nil {
		return fmt.Errorf("validation have a problem")
	}

	date := parts[2]
	parsedDate, err := time.Parse(timeLayout, date)
	if err != nil {
		return fmt.Errorf("date is incorrect")
	}

	currentTime := time.Now()
	// Сравниваем даты
	if parsedDate.Before(currentTime) {
		return fmt.Errorf("date is incorrect. Must be today or later")
	}

	// Создание нового заказа
	newOrder := &json_file.Order{
		ID:        orderID,
		ClientID:  clientID,
		CreatedAt: currentTime,
		ExpiredAt: parsedDate,
	}

	// Добавление заказа в хранилище
	err = storage.AddOrder(newOrder)
	if err != nil {
		return fmt.Errorf("failed to add order: %v", err)
	}
	err = storage.AddOrderToStory(newOrder.ID, "storage/json_file/story_of_orders.json")
	if err != nil {
		return fmt.Errorf("failed to add ortder to story: %v", err)
	}
	fmt.Printf("Creating order for client %d with orderID %d on %s\n", clientID, orderID, date)
	return nil
}
