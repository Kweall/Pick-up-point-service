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
	Orders       map[int64]*Order
	OrderHistory []int64
	Path         string
}

// Конструктор для инициализации хранилища
func NewStorage(path string) (*Storage, error) {
	storage := &Storage{Orders: make(map[int64]*Order), Path: path}
	err := storage.readDataFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}
	return storage, nil
}

func NewHistoryStorage(path string) ([]int64, error) {
	var orderIDs []int64
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open file, err: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&orderIDs)
	if err != nil && err.Error() != "unexpected end of JSON input" {
		return nil, fmt.Errorf("failed to decode order IDs: %v", err)
	}

	return orderIDs, nil
}

// Добавление заказа в data.json
func (s *Storage) AddOrder(order *Order) error {
	s.Orders[order.ID] = order

	// Запись данных в файл
	err := s.writeDataToFile(s.Orders)
	if err != nil {
		return fmt.Errorf("can't write data to file, err: %v", err)
	}

	return nil
}

// Добавление в общую историю заказов
func (s *Storage) AddOrderToStory(orderID int64, path string) error {
	// Проверяем, существует ли уже этот заказ в истории
	for _, id := range s.OrderHistory {
		if id == orderID {
			return fmt.Errorf("orderID %d is already exists", orderID)
		}
	}

	s.OrderHistory = append(s.OrderHistory, orderID)

	err := s.writeDataToHistory(path)
	if err != nil {
		return fmt.Errorf("can't write data to file, err: %v", err)
	}

	fmt.Println("Order added to story successfully.")
	return nil
}

func (s *Storage) writeDataToHistory(path string) error {
	// Открываем файл для записи истории заказов
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("can't open file, err: %v", err)
	}
	defer file.Close()

	// Запись обновленной истории
	encoder := json.NewEncoder(file)
	err = encoder.Encode(s.OrderHistory)
	if err != nil {
		return fmt.Errorf("failed to encode order IDs: %v", err)
	}
	return nil
}

// Чтение данных из файла
func (s *Storage) readDataFromFile() error {
	file, err := os.Open(s.Path)
	if err != nil {
		return fmt.Errorf("can't open file, err: %v", err)
	}
	defer file.Close()

	if stat, _ := file.Stat(); stat.Size() == 0 {
		return nil
	}

	err = json.NewDecoder(file).Decode(&s.Orders)
	if err != nil {
		return fmt.Errorf("can't decode, err: %v", err)
	}

	return nil
}

func (s *Storage) GetAll() (map[int64]*Order, error) {
	return s.Orders, nil
}

// Запись данных в файл
func (s *Storage) writeDataToFile(orders map[int64]*Order) error {
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
	if _, exists := s.Orders[orderID]; !exists {
		return fmt.Errorf("order not found")
	}

	// Удаление заказа из памяти
	delete(s.Orders, orderID)

	// Запись обновленных данных в файл
	err := s.writeDataToFile(s.Orders)
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Order deleted successfully")
	return nil
}

// Выдача заказов клиенту
func (s *Storage) GiveOrdersToClient(orderIDs []int64) error {

	// Проверка существования всех указанных orderID
	orderExists := make(map[int64]bool)
	for _, orderID := range orderIDs {
		if _, exists := s.Orders[orderID]; !exists {
			return fmt.Errorf("OrderID %d not found", orderID)
		}
		orderExists[orderID] = true
	}

	// Проверка срока хранения и принадлежности всех заказов одному клиенту
	var clientID int64
	for _, orderID := range orderIDs {
		order := s.Orders[orderID]

		if clientID == 0 {
			clientID = order.ClientID
		} else if order.ClientID != clientID {
			return fmt.Errorf("all orders should belong to the same client")
		}

		// Проверка срока хранения
		if !order.ExpiredAt.After(time.Now()) {
			return fmt.Errorf("can't give order with ID %d, order expired", order.ID)
		}

		// Проверка, был ли заказ уже выдан
		if !order.RecievedAt.IsZero() {
			return fmt.Errorf("order with ID %d was already given", order.ID)
		}
	}

	// Обновление данных заказов
	for _, orderID := range orderIDs {
		order := s.Orders[orderID]
		order.RecievedAt = time.Now() // Устанавливаем текущую дату как дату получения
	}

	// Запись обновленных данных в файл
	err := s.writeDataToFile(s.Orders)
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Orders have been successfully given to the client.")
	return nil
}

// Принятие возврата заказа
func (s *Storage) AcceptReturn(clientID, orderID int64) error {
	// Проверяем, существует ли заказ в памяти
	order, exists := s.Orders[orderID]
	if !exists || order.ClientID != clientID {
		return fmt.Errorf("order not found or does not belong to the given client")
	}

	// Проверяем, был ли заказ выдан
	if order.RecievedAt.IsZero() {
		return fmt.Errorf("order has not been given")
	}

	// Проверяем, был ли заказ возвращен ранее
	if !order.ReturnedAt.IsZero() {
		return fmt.Errorf("order has already been returned")
	}

	// Проверяем срок возврата
	returnDate := order.RecievedAt.Add(48 * time.Hour)
	if time.Now().After(returnDate) {
		return fmt.Errorf("can't return, return period has expired")
	}

	// Помечаем заказ как возвращённый
	order.ReturnedAt = time.Now()

	// Запись обновленных данных в файл
	err := s.writeDataToFile(s.Orders)
	if err != nil {
		return fmt.Errorf("failed to write updated data to file: %v", err)
	}

	fmt.Println("Order has been successfully returned.")
	return nil
}

func (s *Storage) CheckIfExists(orderID int64) error {
	// Проверяем, принимался ли этот заказ ранее
	for _, id := range s.OrderHistory {
		if id == orderID {
			return fmt.Errorf("this OrderID already used")
		}
	}
	return nil
}
