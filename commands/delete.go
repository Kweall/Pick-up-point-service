package commands

import (
	"fmt"
	"strconv"
)

const countOfArgumentsToDelete = 2

func Delete(storage Storage, parts []string) (err error) {
	if len(parts) != countOfArgumentsToDelete {
		return fmt.Errorf("should be 1 argument: orderID (int)")
	}

	orderID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("orderID is incorrect")
	}

	// Удаляем заказ из файла data.json
	err = storage.DeleteOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("deleteOrderByID failed with error: %v", err)
	}

	return nil
}
