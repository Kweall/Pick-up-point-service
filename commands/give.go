package commands

import (
	"fmt"
	"strconv"
)

func Give(storage Storage, parts []string) (err error) {

	if len(parts) < 1 {
		return fmt.Errorf("should be at least 1 argument: list of orderID's (int) separated by space")
	}

	// Преобразуем список orderID в массив int64
	var orderIDs []int64
	for _, part := range parts[0:] {
		orderID, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return fmt.Errorf("orderID is incorrect")
		}
		orderIDs = append(orderIDs, orderID)
	}

	// Выдача заказов клиенту
	err = storage.GiveOrdersToClient(orderIDs)
	if err != nil {
		return err
	}

	return nil
}
