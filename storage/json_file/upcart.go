package json_file

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Order struct {
	ID         int64  `json:"Order_ID"`
	ClientID   int64  `json:"Client_ID"`
	CreatedAt  string `json:"Created_at"`
	ExpiredAt  string `json:"Expired_at"`
	RecievedAt string `json:"Received_at"`
	ReturnedAt string `json:"Returned_at"`
}

type Storage struct {
	Orders map[int64]*Order
	Path   string
}

// Конструктор для инициализации хранилища
func NewStorage(path string) (*Storage, error) {
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return &Storage{Orders: make(map[int64]*Order), Path: path}, nil
}

// Добавление заказа в data.json
func (s *Storage) AddOrder(order *Order) error {
	err := s.ReadDataFromFile()
	if err != nil {
		return err
	}

	// Добавление заказа в память
	s.Orders[order.ID] = order

	// Запись данных в файл
	err = s.WriteDataToFile()
	if err != nil {
		return err
	}

	return nil
}

// Добавление в общую историю заказов
func (s *Storage) AddOrderToStory(orderID int64, path string) error {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var orderIDs []int64
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&orderIDs)
	if err != nil {
		return fmt.Errorf("failed to decode order IDs: %v", err)
	}

	// Добавляем новый ID заказа
	orderIDs = append(orderIDs, orderID)

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	// Запись обновленной истории
	encoder := json.NewEncoder(file)
	err = encoder.Encode(orderIDs)
	if err != nil {
		return fmt.Errorf("failed to encode order IDs: %v", err)
	}

	return nil
}

// Чтение данных из файла
func (s *Storage) ReadDataFromFile() error {
	file, err := os.OpenFile(s.Path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var orders []*Order
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&orders)
	if err != nil {
		return err
	}

	// Перезаписываем существующие заказы в память
	s.Orders = make(map[int64]*Order)
	for _, order := range orders {
		s.Orders[order.ID] = order
	}
	return nil
}

// Запись данных в файл
func (s *Storage) WriteDataToFile() error {
	file, err := os.OpenFile(s.Path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var orders []*Order
	for _, order := range s.Orders {
		orders = append(orders, order)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(orders)
	if err != nil {
		return err
	}
	return nil
}

// Удаление заказа по ID
func (s *Storage) DeleteOrderByID(orderID int64) error {
	err := s.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	if _, exists := s.Orders[orderID]; !exists {
		fmt.Println("Order not found")
		return nil
	}

	// Удаление заказа из памяти
	delete(s.Orders, orderID)

	// Запись обновленных данных в файл
	err = s.WriteDataToFile()
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Order deleted successfully")
	return nil
}

// Выдача заказов клиенту
func (s *Storage) GiveOrdersToClient(orderIDs []int64) error {

	err := s.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	orderExists := make(map[int64]bool)

	// Проверка существования всех указанных orderID
	for _, orderID := range orderIDs {
		found := false
		for _, order := range s.Orders {
			if order.ID == orderID {
				found = true
				orderExists[orderID] = true
				break
			}
		}
		if !found {
			fmt.Printf("OrderID %d not found.\n", orderID)
			return fmt.Errorf("OrderID %d not found", orderID)
		}
	}

	// Проверка срока хранения и принадлежности всех заказов одному клиенту
	var clientID int64
	for _, order := range s.Orders {
		for _, orderID := range orderIDs {
			if order.ID == orderID {
				if clientID == 0 {
					clientID = order.ClientID
				} else if order.ClientID != clientID {
					//fmt.Println(order.ClientID, clientID)
					fmt.Println("All orders should belong to the same client.")
					return nil
				}

				// Проверка срока хранения
				parsedDate, err := time.Parse("02.01.2006", order.ExpiredAt)
				if err != nil {
					fmt.Printf("invalid date format for orderID %d: %v", order.ID, err)
					return nil
				}

				currentTime := time.Now()
				if parsedDate.Before(currentTime) {
					fmt.Printf("Cannot give order with ID %d, order expired.\n", order.ID)
					return nil
				}

				// Проверка, был ли заказ уже выдан
				if order.RecievedAt != "" {
					fmt.Printf("Order with ID %d was already given.\n", order.ID)
					return fmt.Errorf("order already given")
				}
			}
		}
	}

	// Обновление данных заказов
	for _, orderID := range orderIDs {
		if order, exists := s.Orders[orderID]; exists {
			order.RecievedAt = time.Now().Format("02.01.2006") // Устанавливаем текущую дату как дату получения
		}
	}

	// Запись обновленных данных в файл
	err = s.WriteDataToFile()
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Orders have been successfully given to the client.")
	return nil
}

// Принятие возврата заказа
func (s *Storage) AcceptReturn(clientID, orderID int64) error {
	err := s.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	order, exists := s.Orders[orderID]
	if !exists || order.ClientID != clientID {
		fmt.Println("Order not found or does not belong to the given client")
		return nil
	}

	// Запись обновленных данных в файл
	err = s.WriteDataToFile()
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Order return accepted successfully")
	return nil
}

// Получение списка возвратов с пагинацией
/*func (s *Storage) GetReturns(page, perPage int) ([]*Order, error) {
	err := s.ReadDataFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to read from file: %v", err)
	}

	var returns []*Order

	// Пагинация
	start := (page - 1) * perPage
	end := start + perPage
	if start >= len(returns) {
		return nil, nil // Нет данных для данной страницы
	}
	if end > len(returns) {
		end = len(returns)
	}

	return returns[start:end], nil
}*/
