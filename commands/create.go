package commands

import (
	"bufio"
	"fmt"
	"homework/storage/json_file"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	countOfArgumentsToCreate = 4
	timeLayout               = "02.01.2006"
)

func Create(storage Storage, parts []string) error {

	if len(parts) != countOfArgumentsToCreate {
		return fmt.Errorf("should be 3 arguments: clientID (int), OrderID (int), Expired_date (dd.mm.yyyy)")
	}

	// Получение и преобразование аргументов
	clientID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("clientID is incorrect")
	}

	orderID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("orderID is incorrect")
	}

	// Проверяем, принимался ли этот заказ ранее
	err = storage.Validation(orderID)
	if err != nil {
		return fmt.Errorf("validation have a problem")
	}

	date := parts[3]
	parsedDate, err := time.Parse(timeLayout, date)
	if err != nil {
		return fmt.Errorf("date is incorrect")
	}

	currentTime := time.Now()
	// Сравниваем даты
	if parsedDate.Before(currentTime) {
		return fmt.Errorf("date is incorrect. Must be today or later")
	}

	// Добавляем вес, цену и выбор упаковки
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Specify the weight in kg (up to 3 decimal places)")
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}
	weight, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		return fmt.Errorf("invalid format, err: %v", err)
	}
	weight = math.Round(weight*1000) / 1000

	fmt.Println("Specify the price (int)")
	input, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}
	price, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid format, err: %v", err)
	}

	fmt.Printf("Specify the type of packaging \nbag - 5 units\n box - 20 units\n film - 1 unit\n")
	packaging, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}
	packaging = strings.ToLower(strings.TrimSpace(packaging))
	switch packaging {
	case "bag":
		if weight >= 10 {
			return fmt.Errorf("weight exceeds 10 kg")
		} else {
			price += 5
		}
	case "box":
		if weight >= 30 {
			return fmt.Errorf("weight exceeds 30 kg")
		} else {
			price += 20
		}
	case "film":
		price += 1
	default:
		return fmt.Errorf("there is no such type of packaging")
	}
	var additional_film string
	if packaging != "film" {
		fmt.Println("Would you like to add additional film for 1 units? (yes or no)")
		additional_film, err = reader.ReadString('\n') // Считываем строку до символа новой строки
		if err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}
		additional_film = strings.TrimSpace(additional_film)
		switch additional_film {
		case "yes":
			price += 1
		case "no":
		default:
			return fmt.Errorf("invalid format")
		}
	}

	// Создание нового заказа
	newOrder := &json_file.Order{
		ID:              orderID,
		ClientID:        clientID,
		CreatedAt:       currentTime,
		ExpiredAt:       parsedDate,
		Weight:          weight,
		Price:           price,
		Packaging:       packaging,
		Additional_film: additional_film,
	}

	// Добавление заказа в хранилище
	err = storage.AddOrder(newOrder)
	if err != nil {
		return fmt.Errorf("failed to add order: %v", err)
	}
	err = storage.AddOrderToStory(newOrder.ID, "storage/json_file/story_of_orders.json")
	if err != nil {
		return fmt.Errorf("failed to add ortder to story: %v", err)
	}
	fmt.Printf("Creating order for client %d with orderID %d on %s\n", clientID, orderID, date)
	return nil
}
