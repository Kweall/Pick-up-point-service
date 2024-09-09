package commands

import (
	"fmt"
	"homework/storage"
	"homework/storage/json_file"
	"strconv"
)

const countOfArgumentsToDelete = 1

func Delete(parts []string) (err error) {
	var storage storage.Storage
	storage, err = json_file.NewStorage("storage/json_file/data.json")
	if err != nil {
		fmt.Printf("can't init storage: %v\n", err)
		return err
	}

	if len(parts)-1 != countOfArgumentsToDelete {
		fmt.Println("Should be 1 argument: orderID (int)")
		return nil
	}

	orderID, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("orderID is incorrect")
		return err
	}

	// Удаляем заказ из файла data.json
	err = storage.DeleteOrderByID(int64(orderID))
	if err != nil {
		fmt.Printf("DeleteOrderByID failed with error: %v\n", err)
		return err
	}

	return nil
}
