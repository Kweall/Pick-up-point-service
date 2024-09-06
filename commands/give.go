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
		fmt.Println("Should be at least 1 argument: list of skuID (int) separated by space")
		return nil
	}

	// Преобразуем список skuID в массив int64
	var skuIDs []int64
	for _, part := range parts[1:] {
		skuID, err := strconv.Atoi(part)
		if err != nil {
			fmt.Printf("skuID isn't correct")
			return err
		}
		skuIDs = append(skuIDs, int64(skuID))
	}

	// Выдача заказов клиенту
	err = storage.GiveOrdersToClient(skuIDs)
	if err != nil {
		return err
	}

	return nil
}
