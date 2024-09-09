package main

import (
	"bufio"
	"fmt"
	"homework/commands"
	"os"
	"strings"
	//"gitlab.ozon.dev/go/classroom-15/students/workshop-1/storage/json_file"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Console Utility")
	fmt.Println("---------------------")
	fmt.Println("Type 'help' to see available commands or 'exit' to quit.")

	// Основной цикл программы
	for {
		fmt.Print("> ")                       // Выводим приглашение для ввода команды
		input, err := reader.ReadString('\n') // Считываем строку до символа новой строки
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Удаляем символы перевода строки и пробелы
		command := strings.TrimSpace(input)

		parts := strings.Fields(command)

		// Обрабатываем команды
		switch parts[0] {
		case "help":
			fmt.Println("List of commands:")
			fmt.Println(" - help: Show this help message")
			fmt.Println(" - exit: Quit the application")
			fmt.Printf(" - CREATE: Receive the order from the courier and record it in the database\n\tlayout: CREATE clientID orderID dd.mm.yyyy\n\texample: CREATE 10 20 13.09.2024\n\v")
			fmt.Printf(" - DELETE: Return the order to the courier using the orderID and delete the entry from the file\n\tlayout: DELETE orderID\n\texample: DELETE 10\n\v")
			fmt.Printf(" - GIVE: Issue orders to ONE! client\n\tlayout: GIVE []orderID\n\texample: GIVE 10 20 30\n\v")
			fmt.Printf(" - GET_ORDERS: Get a list of customer orders\n\tlayout: GET_ORDERS clientID (+limit optionally)\n\texample: GET_ORDERS 10 (5)\n\v")
			fmt.Printf(" - ACCEPT_RETURN: Accept return from customer\n\tlayout: ACCEPT_RETURN clientID orderID\n\texample: ACCEPT_RETURN 10 20\n\v")
			fmt.Printf(" - GET_RETURNS: Get a list of returns (max 5 per page)\n\tlayout: GET_RETURNS page\n\texample: GET_RETURNS 1\n")
		case "exit":
			fmt.Println("Exiting...")
			return // Завершаем программу
		case "CREATE": // Принять заказ от курьера
			commands.Create(parts)
		case "DELETE": // Вернуть заказ курьеру
			commands.Delete(parts)
		case "GIVE": // Выдать заказы клиента
			commands.Give(parts)
		case "GET_ORDERS": // Получить список заказов определенного клиента
			commands.GetOrders(parts)
		case "ACCEPT_RETURN": // Принять возврат от клиента
			commands.AcceptReturn(parts)
		case "GET_RETURNS": // Получить список возвратов (номер страницы, количество записей на одной странице - 5)
			commands.GetReturns(parts)
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}
