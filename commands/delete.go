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
		fmt.Println("Should be 1 argument: skuID (int)")
		return nil
	}

	skuID, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("skuID isn't correct")
		return err
	}

	// Удаляем элемент из всех корзин пользователей по skuID
	err = storage.DeleteItemBySkuID(int64(skuID))
	if err != nil {
		fmt.Printf("DeleteItemBySkuID failed with error: %v\n", err)
		return err
	}

	return nil
}
