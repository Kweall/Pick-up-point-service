package commands

import (
	"fmt"
	"homework/storage"
	"homework/storage/json_file"
	"strconv"
)

func Give(parts []string) (err error) {
	var storage storage.Storage
	storage, err = json_file.NewStorage("storage/json_file/data.json")
	if err != nil {
		fmt.Printf("can't init storage: %v\n", err)
		return err
	}

	if len(parts) < 2 {
		fmt.Println("Should be at least 1 argument: list of orderID's (int) separated by space")
		return nil
	}

	// Преобразуем список orderID в массив int64
	var orderIDs []int64
	for _, part := range parts[1:] {
		orderID, err := strconv.Atoi(part)
		if err != nil {
			fmt.Printf("orderID is incorrect")
			return err
		}
		orderIDs = append(orderIDs, int64(orderID))
	}

	// Выдача заказов клиенту
	err = storage.GiveOrdersToClient(orderIDs)
	if err != nil {
		return err
	}

	return nil
}
