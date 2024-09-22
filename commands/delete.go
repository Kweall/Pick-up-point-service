package commands

import (
	"fmt"
	"strconv"
)

const countOfArgumentsToDelete = 1

func Delete(storage Storage, parts []string) (err error) {
	if len(parts) != countOfArgumentsToDelete {
		return fmt.Errorf("should be 1 argument: orderID (int)")
	}

	orderID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("orderID is incorrect")
	}

	// Удаляем заказ из файла data.json
	err = storage.DeleteOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
