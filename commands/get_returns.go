package commands

import (
	"fmt"
	"homework/storage/json_file"
	"log"

	"github.com/nsf/termbox-go"
)

func GetReturns(storage Storage, parts []string) error {

	orders, err := storage.GetAll()
	if err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	// Поиск всех возвратов
	var returns []json_file.Order
	for _, order := range orders {
		if !order.ReturnedAt.IsZero() {
			returns = append(returns, *order)
		}
	}

	if len(returns) == 0 {
		return fmt.Errorf("no returns found")
	}

	// Инициализация termbox
	err = termbox.Init()
	if err != nil {
		log.Fatalf("Failed to initialize termbox: %v", err)
	}
	defer termbox.Close()

	// Переменные для управления прокруткой
	index := 0

	// Функция для отображения текущего возврата
	displayCurrentReturn := func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		order := returns[index]
		output := fmt.Sprintf("Order ID: %d,\t Client ID: %d,\t Date of return: %s\t\n", order.ID, order.ClientID, order.ReturnedAt)
		for i, ch := range output {
			termbox.SetCell(i, 0, ch, termbox.ColorWhite, termbox.ColorDefault)
		}
		termbox.Flush()
	}

	displayCurrentReturn()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyArrowDown {
				if index < len(returns)-1 {
					index++
					displayCurrentReturn()
				}
			} else if ev.Key == termbox.KeyArrowUp {
				if index > 0 {
					index--
					displayCurrentReturn()
				}
			} else if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				return nil // Выход из функции
			}
		case termbox.EventError:
			log.Printf("Termbox event error: %v", ev.Err)
			return ev.Err
		}
	}
}
