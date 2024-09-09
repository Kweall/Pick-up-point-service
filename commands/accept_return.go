package commands

import (
	"fmt"
	"homework/storage/json_file"
	"strconv"
	"time"
)

const countOfArgumentstoReturn = 2

func AcceptReturn(parts []string) error {
	if len(parts)-1 != countOfArgumentstoReturn {
		fmt.Println("Should be 2 arguments: clientID (int), orderID (int)")
		return nil
	}

	clientID, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("clientID is incorrect")
		return err
	}

	orderID, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Println("orderID is incorrect")
		return err
	}

	storage, err := json_file.NewStorage("storage/json_file/data.json")
	if err != nil {
		return fmt.Errorf("can't init storage: %v", err)
	}

	err = storage.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	// Ищем пользователя и заказ
	var orderFound bool
	for _, order := range storage.Orders {
		if order.ClientID == int64(clientID) && order.ID == int64(orderID) {
			// Проверяем, был ли заказ выдан
			if order.RecievedAt == "" {
				fmt.Println("Order has not been given")
				return nil
			}
			// Проверяем, был ли заказ возвращен раннее
			if order.ReturnedAt != "" {
				fmt.Println("Order has already been returned")
				return nil
			}
			// Проверяем срок возврата
			returnDate, err := time.Parse("02.01.2006", order.RecievedAt)
			if err != nil {
				fmt.Println("Invalid return date format")
				return err
			}
			returnDate = returnDate.Add(time.Hour * 48)

			currentDate := time.Now()
			if currentDate.After(returnDate) {
				fmt.Println("Can't return, return period has expired")
				return nil
			}

			// Помечаем заказ как возвращенный
			order.ReturnedAt = time.Now().Format("02.01.2006")
			orderFound = true
			break
		}
	}

	if !orderFound {
		fmt.Printf("Order with ID %d for client %d not found.\n", orderID, clientID)
		return nil
	}

	err = storage.WriteDataToFile()
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	fmt.Printf("Return accepted for Order ID: %d, ClientID: %d\n", orderID, clientID)
	return nil
}
