package json_file

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type OutputData struct {
	Users []User `json:"users"`
}

type User struct {
	UserID int64  `json:"userID"`
	Items  []Item `json:"items"`
}

type Item struct {
	ID     int64  `json:"sku"`
	Date   string `json:"date"`
	Valid  bool   `json:"at the pick-up point"`
	Return bool   `json:"has returned"`
}

type Storage struct {
	Carts map[int64]map[int64]string
	Path  string
}

func NewStorage(path string) (*Storage, error) {
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return &Storage{Carts: make(map[int64]map[int64]string), Path: path}, nil
}

func (s *Storage) AddItemToCart(userID, skuID int64, date string) (err error) {
	_, err = s.ReadDataFromFile()
	if err != nil {
		return err
	}

	if _, ok := s.Carts[userID]; !ok {
		s.Carts[userID] = make(map[int64]string)
	}

	s.Carts[userID][skuID] = date

	// Создание структуры для записи данных
	var data OutputData
	for userID, items := range s.Carts {
		user := User{UserID: userID}
		for itemID, date := range items {
			user.Items = append(user.Items, Item{ID: itemID, Date: date, Valid: true})
		}
		data.Users = append(data.Users, user)
	}

	// Запись в файл
	err = s.WriteDataFromFile(data)
	if err != nil {
		return err
	}
	return nil
}

// Проверка существования элемента в корзине
func (s *Storage) CheckIfItemExists(userID, skuID int64, date string) (bool, error) {
	_, err := s.ReadDataFromFile()
	if err != nil {
		return false, fmt.Errorf("failed to read from file: %v", err)
	}
	if cart, ok := s.Carts[userID]; ok {
		if existingDate, exists := cart[skuID]; exists && existingDate == date {
			return true, nil
		}
	}
	return false, nil
}

// Удаление заказа с указанным skuID из файла
func (s *Storage) DeleteItemBySkuID(skuID int64) error {
	_, err := s.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	// Флаг для проверки, был ли элемент найден и удален
	itemDeleted := false

	// Итерация по всем пользователям и удаление элементов с указанным skuID
	for userID, cart := range s.Carts {
		if _, exists := cart[skuID]; exists {
			delete(cart, skuID)
			itemDeleted = true

			// Если корзина пользователя пуста после удаления, удаляем запись о пользователе
			if len(cart) == 0 {
				delete(s.Carts, userID)
			}
		}
	}

	if !itemDeleted {
		fmt.Println("Item with the given skuID not found in any user's cart")
		return nil
	}

	// Создание структуры для записи данных
	var data OutputData
	for userID, items := range s.Carts {
		user := User{UserID: userID}
		for itemID, date := range items {
			user.Items = append(user.Items, Item{ID: itemID, Date: date, Valid: true})
		}
		data.Users = append(data.Users, user)
	}
	// Запись обновленных данных в файл
	err = s.WriteDataFromFile(data)
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Item deleted successfully")
	return nil
}

func (s *Storage) GiveOrdersToClient(skuIDs []int64) error {

	// Чтение данных из файла
	file, err := os.OpenFile(s.Path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var data OutputData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	// Проверяем заказы
	var foundUserID int64
	orderExists := make(map[int64]bool) // Карта для проверки существования заказов

	// Собираем все существующие заказы и проверяем их
	//fmt.Println(data)
	for _, user := range data.Users {
		for _, item := range user.Items {
			orderExists[item.ID] = true
		}
	}

	for _, skuID := range skuIDs {
		if !orderExists[skuID] {
			fmt.Printf("Order with skuID %d not found.\n", skuID)
			return err
		}

		// Проверяем срок хранения и принадлежность одного пользователя
		found := false
		for _, user := range data.Users {
			for _, item := range user.Items {
				if item.ID == skuID {
					if !item.Valid {
						fmt.Println("This sku was given before")
						return nil
					}
					parsedDate, err := time.Parse("02.01.2006", item.Date)
					if err != nil {
						return fmt.Errorf("invalid date format for skuID %d: %v", skuID, err)
					}

					current_time, _ := time.Parse("02.01.2006", time.Now().Format("02.01.2006"))
					if parsedDate.Before(current_time) {
						fmt.Printf("Cannot give order with skuID %d, order expired.\n", skuID)
						return nil
					}

					if foundUserID == 0 {
						foundUserID = user.UserID
					} else if foundUserID != user.UserID {
						fmt.Println("All orders should belong to the same user.")
						return nil
					}

					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			fmt.Printf("Order with skuID %d not found.\n", skuID)
			return nil
		}
	}

	// Обновляем статус заказов
	for userIndex := range data.Users {
		user := &data.Users[userIndex]
		for itemIndex := range user.Items {
			item := &user.Items[itemIndex]
			for _, skuID := range skuIDs {
				if item.ID == skuID {
					if !item.Return {
						last_day, _ := time.Parse("02.01.2006", time.Now().Format("02.01.2006"))
						last_day = last_day.Add(48 * time.Hour) // Последний день для возврата
						item.Date = last_day.Format("02.01.2006")
						item.Valid = false
					} else {
						fmt.Println("This item has returned")
						return nil
					}
				}
			}
		}
	}

	// Записываем обновленные данные в файл
	err = s.WriteDataFromFile(data)
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Orders have been successfully given to the client.")
	return nil
}

func (s *Storage) ReadDataFromFile() (*OutputData, error) {
	file, err := os.OpenFile(s.Path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data OutputData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	// Преобразование данных в карту
	s.Carts = make(map[int64]map[int64]string)
	for _, user := range data.Users {
		if _, ok := s.Carts[user.UserID]; !ok {
			s.Carts[user.UserID] = make(map[int64]string)
		}
		for _, item := range user.Items {
			s.Carts[user.UserID][item.ID] = item.Date
		}
	}
	return &data, nil
}

func (s *Storage) WriteDataFromFile(data OutputData) (err error) {
	file, err := os.OpenFile(s.Path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}
