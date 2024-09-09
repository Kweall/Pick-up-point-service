package commands

import (
	"encoding/json"
	"fmt"
	"homework/storage/json_file"
	"os"
	"strconv"
	"time"
)

const countOfArgumentsToCreate = 3

func Create(parts []string) error {
	// Инициализация хранилища
	storage, err := json_file.NewStorage("storage/json_file/data.json")
	if err != nil {
		fmt.Printf("can't init storage: %v\n", err)
		return err
	}

	if len(parts)-1 != countOfArgumentsToCreate {
		fmt.Println("Should be 3 arguments: clientID (int), OrderID (int), Expired_date (dd.mm.yyyy)")
		return nil
	}

	// Получение и преобразование аргументов
	clientID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		fmt.Println("ClientID isn't correct")
		return err
	}

	orderID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		fmt.Println("OrderID isn't correct")
		return err
	}

	//
	// Проверяем, принимался ли этот заказ ранее
	story_file, err := os.OpenFile("storage/json_file/story_of_orders.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer story_file.Close()

	var orderIDs []int64
	decoder := json.NewDecoder(story_file)
	err = decoder.Decode(&orderIDs)
	if err != nil {
		return fmt.Errorf("failed to decode order IDs: %v", err)
	}
	for i := 0; i < len(orderIDs); i++ {
		if orderIDs[i] == orderID {
			fmt.Printf("This OrderID already used\n")
			return nil
		}
	}
	//

	date := parts[3]
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		fmt.Println("Date isn't correct")
		return err
	}

	// Получаем текущую дату
	currentTime, _ := time.Parse("02.01.2006", time.Now().Format("02.01.2006"))

	// Сравниваем даты
	if parsedDate.Before(currentTime) {
		fmt.Println("Date is incorrect. Must be today or later.")
		return nil
	}

	// Проверка существования заказа
	for _, order := range storage.Orders {
		if order.ClientID == clientID && order.ID == orderID && order.CreatedAt == date {
			fmt.Println("This order already exists.")
			return nil
		}
	}

	// Создание нового заказа
	newOrder := &json_file.Order{
		ID:        orderID,
		ClientID:  clientID,
		CreatedAt: time.Now().Format("02.01.2006"),
		ExpiredAt: date,
	}

	// Добавление заказа в хранилище
	err = storage.AddOrder(newOrder)
	if err != nil {
		fmt.Printf("Failed to add order: %v\n", err)
		return err
	}
	err = storage.AddOrderToStory(newOrder.ID, "storage/json_file/story_of_orders.json")
	if err != nil {
		fmt.Printf("Failed to add order to story: %v\n", err)
		return err
	}
	fmt.Printf("Creating order for client %d with orderID %d on %s\n", clientID, orderID, date)
	return nil
}
