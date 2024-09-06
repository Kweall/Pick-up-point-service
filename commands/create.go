package commands

import (
	"fmt"
	"homework/storage"
	"homework/storage/json_file"
	"strconv"
	"time"
)

const countOfArgumentsToCreate = 3

func Create(parts []string) (err error) {
	var storage storage.Storage
	storage, err = json_file.NewStorage("storage/json_file/data.json")
	if err != nil {
		fmt.Printf("can't init storage %v", err)
		//return fmt.Errorf("can't init storage %v", err)
		return err
	}
	if len(parts)-1 == countOfArgumentsToCreate {
		userID, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("userID isn't correct")
			//return fmt.Errorf("userID isn't correct")
			return err
		}
		skuID, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("skuID isn't correct")
			//return fmt.Errorf("skuID isn't correct")
			return err
		}
		date := parts[3]
		// Парсим дату из строки. Указываем формат, соответствующий строке.
		parsedDate, err := time.Parse("02.01.2006", date)
		if err != nil {
			fmt.Println("Date isn't correct")
			//return fmt.Errorf("Date isn't correct")
			return err
		}

		// Получаем текущую дату
		currentTime, _ := time.Parse("02.01.2006", time.Now().Format("02.01.2006"))

		// Сравниваем даты
		if parsedDate.After(currentTime) || parsedDate.Equal(currentTime) {
			// Проверяем, существует ли уже такая запись
			exists, err := storage.CheckIfItemExists(int64(userID), int64(skuID), date)
			if err != nil {
				fmt.Printf("Error checking if item exists: %v\n", err)
				return err
			}
			if exists {
				fmt.Println("This item already exists in the cart.")
				return fmt.Errorf("item already exists")
			}
			err = storage.AddItemToCart(int64(userID), int64(skuID), date)
			if err != nil {
				fmt.Printf("AddItemToCart failed with error %v", err)
				return err
			}
			fmt.Println("Creating", parts[1], "for", parts[2], "at", parts[3])
		} else {
			fmt.Println("Date incorrect")
		}

	} else {
		fmt.Println("Should be 3 arguments: userID (int), skuID (int), date (dd.mm.yyyy)")
	}
	return nil
}
