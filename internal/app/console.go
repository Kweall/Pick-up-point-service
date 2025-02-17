package app

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type Task struct {
	Command  string
	Parts    []string
	DataFlag *string
}

func worker(ctx context.Context, id int, taskChan chan Task, service *Service, wg *sync.WaitGroup) error {
	defer wg.Done()
	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return fmt.Errorf("worker %d have a problem", id)
			}
			if err := runCommand(ctx, service, task.Parts, task.DataFlag); err != nil {
				return fmt.Errorf("worker %d: error executing command: %v", id, err)
			}
		case <-ctx.Done():
			return fmt.Errorf("worker %d: received shutdown signal", id)
		}
	}
}

func RunCLI(ctx context.Context, service *Service, dataFlag *string) error {
	fmt.Println("Console Utility")
	fmt.Println("---------------------")
	fmt.Println("Type 'help' to see available commands or 'exit' to quit.")
	taskChan := make(chan Task)
	var wg sync.WaitGroup
	numWorkers := 2
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, taskChan, service, &wg)
	}

	shutdownChan := make(chan struct{})

	go func() {
		<-ctx.Done()
		fmt.Println("\nReceived shutdown signal, terminating gracefully...")
		close(taskChan)
		wg.Wait()
		close(shutdownChan)
	}()

	go userInput(ctx, service, taskChan, &wg, shutdownChan, cancel, &numWorkers, dataFlag)

	<-shutdownChan
	fmt.Println("Gracefully shutdown completed.")
	return nil
}

func userInput(ctx context.Context, service *Service, taskChan chan Task, wg *sync.WaitGroup, shutdownChan chan struct{}, cancel context.CancelFunc, numWorkers *int, dataFlag *string) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-shutdownChan:
			return fmt.Errorf("closing")
		default:
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading input: %v", err)
			}

			command := strings.TrimSpace(input)
			parts := strings.Fields(command)
			if len(parts) == 0 {
				continue
			}

			switch parts[0] {
			case "help":
				printHelp()
			case "exit":
				fmt.Println("Exiting...")
				cancel()
				return nil
			case "AddWorkers":
				if err = addWorkers(parts, ctx, service, taskChan, wg, numWorkers); err != nil {
					return err
				}
			default:
				t := Task{
					Command:  parts[0],
					Parts:    parts,
					DataFlag: dataFlag,
				}
				if t.Command == "AddOrder" {
					var mu sync.Mutex
					mu.Lock()
					defer mu.Unlock()
					if err := runCommand(ctx, service, t.Parts, t.DataFlag); err != nil {
						fmt.Println(err)
					}
					continue
				}
				select {
				case taskChan <- t:
				case <-ctx.Done():
					return nil
				}
			}
		}
	}
}

func addWorkers(parts []string, ctx context.Context, service *Service, taskChan chan Task, wg *sync.WaitGroup, numWorkers *int) error {
	if len(parts) != 2 {
		return fmt.Errorf("usage: AddWorkers count")
	}
	count, err := strconv.Atoi(parts[1])
	if err != nil || count <= 0 {
		return fmt.Errorf("invalid count for workers")
	}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go worker(ctx, *numWorkers+1, taskChan, service, wg)
		*numWorkers++
	}
	fmt.Printf("%d workers added. Total workers: %d\n", count, *numWorkers)
	return nil
}

func runCommand(ctx context.Context, service *Service, parts []string, dataFlag *string) error {
	var (
		resp    any
		respErr error
		err     error
	)
	if len(parts) == 0 {
		return fmt.Errorf("input message is empty")
	}
	switch parts[0] {
	case "AddOrder":
		var req AddOrderRequest
		if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
			log.Fatal(err)
		}
		resp, respErr = service.AddOrder(ctx, &req, parts)
	case "DeleteOrder":
		var req DeleteOrderRequest
		if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
			log.Fatal(err)
		}
		resp, respErr = service.DeleteOrder(ctx, &req, parts)
	case "GetOrders":
		var req GetOrdersRequest
		if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
			log.Fatal(err)
		}
		resp, respErr = service.GetOrders(ctx, &req, parts)
	case "GiveOrders":
		var req GiveOrderRequest
		if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
			log.Fatal(err)
		}
		resp, respErr = service.GiveOrders(ctx, &req, parts)
	case "AcceptReturn":
		var req AcceptReturnRequest
		if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
			log.Fatal(err)
		}
		resp, respErr = service.AcceptReturn(ctx, &req, parts)
	case "GetReturns":
		var req GetReturnsRequest
		if err = json.Unmarshal([]byte(*dataFlag), &req); err != nil {
			log.Fatal(err)
		}
		resp, respErr = service.GetReturns(ctx, &req, parts)
	default:
		respErr = fmt.Errorf("unknown command: %s", parts[0])
	}
	if err != nil {
		return err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("error with marshaling: %v", err)
	}
	if respErr != nil {
		return fmt.Errorf("resp: %s, err: %v", data, respErr)
	}
	return nil
}

func printHelp() {
	fmt.Println("List of commands:")
	fmt.Println(" - help: Show this help message")
	fmt.Println(" - exit: Quit the application")
	fmt.Println(" - AddWorkers: Add more workers for execution tasks!\n\tlayout: AddWorkers count\n\v")
	fmt.Printf(" - AddOrder: Receive the order from the courier and record it in the database\n\tlayout: AddOrder clientID orderID dd.mm.yyyy\n\texample: AddOrder 10 20 20.10.2024\n\v")
	fmt.Printf(" - DeleteOrder: Return the order to the courier using the orderID and delete the entry from the file\n\tlayout: DeleteOrder orderID\n\texample: DeleteOrder 10\n\v")
	fmt.Printf(" - GiveOrder: Issue orders to ONE! client\n\tlayout: GiveOrder []orderID\n\texample: GiveOrder 10 20 30\n\v")
	fmt.Printf(" - GetOrders: Get a list of customer orders\n\tlayout: GetOrders clientID (+limit optionally)\n\texample: GetOrders 10 (5)\n\v")
	fmt.Printf(" - AcceptReturn: Accept return from customer\n\tlayout: AcceptReturn clientID orderID\n\texample: AcceptReturn 10 20\n\v")
	fmt.Printf(" - GetReturns: Get a list of returns (max 5 per page)\n\tlayout: GetReturns page\n\texample: GetReturns 1\n")
}
