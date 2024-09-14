package json_file

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Order struct {
	ID             int64     `json:"order_id"`
	ClientID       int64     `json:"client_id"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiredAt      time.Time `json:"expired_at"`
	RecievedAt     time.Time `json:"received_at"`
	ReturnedAt     time.Time `json:"returned_at"`
	Weight         float64   `json:"weight"`
	Price          int64     `json:"price"`
	Packaging      string    `json:"packaging"`
	AdditionalFilm string    `json:"additional_film"`
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
	orders, err := s.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("can't read file, err: %v", err)
	}

	// Добавление заказа в память
	orders[order.ID] = order

	// Запись данных в файл
	err = s.WriteDataToFile(orders)
	if err != nil {
		return fmt.Errorf("can't write data to file, err: %v", err)
	}

	return nil
}

// Добавление в общую историю заказов
func (s *Storage) AddOrderToStory(orderID int64, path string) error {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("can't open file, err: %v", err)
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
		return fmt.Errorf("can't truncate, err: %v", err)
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
func (s *Storage) ReadDataFromFile() (map[int64]*Order, error) {
	file, err := os.OpenFile(s.Path, os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("can't open file, err: %v", err)
	}
	defer file.Close()

	var orders map[int64]*Order
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&orders)
	if err != nil {
		return nil, fmt.Errorf("can't decode, err: %v", err)
	}

	return orders, nil
}

func (s *Storage) GetAll() (map[int64]*Order, error) {
	return s.ReadDataFromFile()
}

// Запись данных в файл
func (s *Storage) WriteDataToFile(orders map[int64]*Order) error {
	file, err := os.OpenFile(s.Path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("can't open file, err: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(orders)
	if err != nil {
		return fmt.Errorf("can't encode, err: %v", err)
	}
	return nil
}

// Удаление заказа по ID
func (s *Storage) DeleteOrderByID(orderID int64) error {
	orders, err := s.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	if _, exists := orders[orderID]; !exists {
		return fmt.Errorf("Order not found")
	}

	// Удаление заказа из памяти
	delete(orders, orderID)

	// Запись обновленных данных в файл
	err = s.WriteDataToFile(orders)
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Order deleted successfully")
	return nil
}

// Выдача заказов клиенту
func (s *Storage) GiveOrdersToClient(orderIDs []int64) error {

	orders, err := s.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	orderExists := make(map[int64]bool)

	// Проверка существования всех указанных orderID
	for _, orderID := range orderIDs {
		found := false
		for _, order := range orders {
			if order.ID == orderID {
				found = true
				orderExists[orderID] = true
				break
			}
		}
		if !found {
			return fmt.Errorf("OrderID %d not found", orderID)
		}
	}

	// Проверка срока хранения и принадлежности всех заказов одному клиенту
	var clientID int64
	for _, order := range orders {
		for _, orderID := range orderIDs {
			if order.ID == orderID {
				if clientID == 0 {
					clientID = order.ClientID
				} else if order.ClientID != clientID {
					return fmt.Errorf("all orders should belong to the same client")
				}

				// Проверка срока хранения
				parsedDate := order.ExpiredAt

				if !parsedDate.After(time.Now()) {
					return fmt.Errorf("can't give order with ID %d, order expired", order.ID)
				}

				// Проверка, был ли заказ уже выдан
				if !order.RecievedAt.IsZero() {
					return fmt.Errorf("order with ID %d was already given", order.ID)
				}
			}
		}
	}

	// Обновление данных заказов
	for _, orderID := range orderIDs {
		if order, exists := orders[orderID]; exists {
			order.RecievedAt = time.Now() // Устанавливаем текущую дату как дату получения
		}
	}

	// Запись обновленных данных в файл
	err = s.WriteDataToFile(orders)
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Orders have been successfully given to the client.")
	return nil
}

// Принятие возврата заказа
func (s *Storage) AcceptReturn(clientID, orderID int64) error {
	orders, err := s.ReadDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	order, exists := orders[orderID]
	if !exists || order.ClientID != clientID {
		return fmt.Errorf("Order not found or does not belong to the given client")
	}
	// Ищем пользователя и заказ
	var orderFound bool
	if order.ClientID == clientID && order.ID == orderID {
		// Проверяем, был ли заказ выдан
		if order.RecievedAt.IsZero() {
			return fmt.Errorf("order has not been given")
		}
		// Проверяем, был ли заказ возвращен раннее
		if !order.ReturnedAt.IsZero() {
			return fmt.Errorf("order has already been returned")
		}
		// Проверяем срок возврата
		returnDate := order.RecievedAt
		returnDate = returnDate.Add(time.Hour * 48)

		currentDate := time.Now()
		if currentDate.After(returnDate) {
			return fmt.Errorf("can't return, return period has expired")
		}

		// Помечаем заказ как возвращенный
		orders[orderID].ReturnedAt = time.Now()
		orderFound = true
	}

	if !orderFound {
		return fmt.Errorf("order with ID %d for client %d not found", orderID, clientID)
	}
	// Запись обновленных данных в файл
	err = s.WriteDataToFile(orders)
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	return nil
}

func (s *Storage) Validation(orderID int64) error {
	// Проверяем, принимался ли этот заказ ранее
	story_file, err := os.OpenFile("storage/json_file/story_of_orders.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("can't open file, err: %v", err)
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
			return fmt.Errorf("this OrderID already used")
		}
	}
	return nil
}
