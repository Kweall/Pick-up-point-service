package commands

import (
	"fmt"
	"homework/storage/json_file"
	"strconv"
	"time"
)

const countOfArgumentstoReturn = 2

func AcceptReturn(parts []string) error {
	// Проверяем количество аргументов
	if len(parts)-1 != countOfArgumentstoReturn {
		fmt.Println("Should be 2 arguments: userID (int), skuID (int)")
		return nil
	}

	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("userID isn't correct")
		return err
	}

	skuID, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Println("skuID isn't correct")
		return err
	}

	storage, err := json_file.NewStorage("storage/json_file/data.json")
	if err != nil {
		return fmt.Errorf("can't init storage: %v", err)
	}

	data, err := storage.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	// Ищем пользователя и заказ
	var orderFound bool
	for i, user := range data.Users {
		if user.UserID == int64(userID) {
			for j, item := range user.Items {
				if item.ID == int64(skuID) {
					time_for_return, err := time.Parse("02.01.2006", data.Users[i].Items[j].Date)
					if err != nil {
						fmt.Println("Date incorrect")
						return nil
					}
					today, _ := time.Parse("02.01.2006", time.Now().Format("02.01.2006"))
					if time_for_return.Before(today) {
						fmt.Println("Can't return, date is expired")
						return nil
					} else if data.Users[i].Items[j].Return {
						fmt.Println("Already hes returned")
						return nil
					} else {
						// Меняем значение Return на true
						data.Users[i].Items[j].Return = true
						data.Users[i].Items[j].Valid = true
						orderFound = true
						break
					}
				}
			}
		}
		if orderFound {
			break
		}
	}

	if !orderFound {
		fmt.Printf("Order with ID %d for user %d not found.\n", skuID, userID)
		return nil
	}

	err = storage.WriteDataFromFile(*data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	fmt.Printf("Return accepted for Order ID: %d, User ID: %d\n", skuID, userID)
	return nil
}
