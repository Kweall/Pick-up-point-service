package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"homework/internal/app"
	"homework/internal/config"
	"homework/internal/storage/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	dataFlag = flag.String("data", "{}", "data in JSON format")
)

func main() {
	// const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, cfg.PsqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// err = app.ClearTables(pool)
	// err = app.GenerateFakeOrders(pool, 100)
	// if err != nil {
	// 	log.Fatalf("Error generating fake orders: %v", err)
	// }

	fmt.Println("Successfully generated test data!")
	storageFacade := newStorageFacade(pool)
	service := app.NewService(storageFacade)

	var (
		resp    any
		respErr error
	)
	fmt.Println("Console Utility")
	fmt.Println("---------------------")
	fmt.Println("Type 'help' to see available commands or 'exit' to quit.")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		parts := strings.Fields(command)
		switch parts[0] {
		case "help":
			fmt.Println("List of commands:")
			fmt.Println(" - help: Show this help message")
			fmt.Println(" - exit: Quit the application")
			fmt.Printf(" - AddOrder: Receive the order from the courier and record it in the database\n\tlayout: AddOrder clientID orderID dd.mm.yyyy\n\texample: AddOrder 10 20 20.10.2024\n\v")
			fmt.Printf(" - DeleteOrder: Return the order to the courier using the orderID and delete the entry from the file\n\tlayout: DeleteOrder orderID\n\texample: DeleteOrder 10\n\v")
			fmt.Printf(" - GiveOrder: Issue orders to ONE! client\n\tlayout: GiveOrder []orderID\n\texample: GiveOrder 10 20 30\n\v")
			fmt.Printf(" - GetOrders: Get a list of customer orders\n\tlayout: GetOrders clientID (+limit optionally)\n\texample: GetOrders 10 (5)\n\v")
			fmt.Printf(" - AcceptReturn: Accept return from customer\n\tlayout: AcceptReturn clientID orderID\n\texample: AcceptReturn 10 20\n\v")
			fmt.Printf(" - GetReturns: Get a list of returns (max 5 per page)\n\tlayout: GetReturns page\n\texample: GetReturns 1\n")
		case "exit":
			fmt.Println("Exiting...")
			return
		case "AddOrder":
			var req app.AddOrderRequest
			if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
				log.Fatal(err)
			}
			resp, respErr = service.AddOrder(ctx, &req, parts)
		case "DeleteOrder":
			var req app.DeleteOrderRequest
			if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
				log.Fatal(err)
			}
			resp, respErr = service.DeleteOrder(ctx, &req, parts)
		case "GetOrders":
			var req app.GetOrdersRequest
			if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
				log.Fatal(err)
			}
			resp, respErr = service.GetOrders(ctx, &req, parts)
		case "GiveOrders":
			var req app.GiveOrderRequest
			if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
				log.Fatal(err)
			}
			resp, respErr = service.GiveOrders(ctx, &req, parts)
		case "AcceptReturn":
			var req app.AcceptReturnRequest
			if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
				log.Fatal(err)
			}
			resp, respErr = service.AcceptReturn(ctx, &req, parts)
		case "GetReturns":
			var req app.GetReturnsRequest
			if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
				log.Fatal(err)
			}
			resp, respErr = service.GetReturns(ctx, &req, parts)
		default:
			respErr = fmt.Errorf("unknown command: %s", parts[0])
		}
		if err != nil {
			fmt.Println(err)
		}

		data, _ := json.Marshal(resp)
		resp = nil
		if respErr != nil {
			log.Printf("resp: %s, err: %v\n", data, respErr)
			respErr = nil
		}
	}
}

func newStorageFacade(pool *pgxpool.Pool) app.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgRepository(txManager)

	return app.NewStorageFacade(txManager, pgRepository)
}
