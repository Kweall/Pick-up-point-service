package app

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type AddOrderRequest struct {
	OrderID        int64     `json:"OrderID"`
	ClientID       int64     `json:"ClientID"`
	CreatedAt      time.Time `json:"CreatedAt"`
	ExpiredAt      time.Time `json:"ExpiredAt"`
	Weight         float64   `json:"Weight"`
	Price          int64     `json:"Price"`
	Packaging      string    `json:"Packaging"`
	AdditionalFilm string    `json:"AdditionalFilm"`
}

type AddOrderResponse struct {
}

const (
	countOfArgumentsToCreate = 3
	timeLayout               = "02.01.2006"
)

func (s *Service) AddOrder(ctx context.Context, req *AddOrderRequest, parts []string) (*AddOrderResponse, error) {

	if len(parts)-1 != countOfArgumentsToCreate {
		return nil, fmt.Errorf("should be 3 arguments: clientID (int), OrderID (int), Expired_date (dd.mm.yyyy)")
	}
	req.AdditionalFilm = "no"
	req.CreatedAt = time.Now().Truncate(time.Minute)
	var err error

	req.ClientID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("clientID is incorrect")
	}
	req.OrderID, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("orderID is incorrect")
	}
	req.ExpiredAt, err = time.Parse(timeLayout, parts[3])
	if err != nil {
		return nil, fmt.Errorf("date is incorrect")
	}
	if req.ExpiredAt.Before(req.CreatedAt) {
		return nil, fmt.Errorf("date is incorrect. Must be today or later")
	}

	// Добавляем вес, цену и выбор упаковки
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Specify the weight in kg (up to 3 decimal places)\n> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}
	req.Weight, err = strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid format, err: %v", err)
	}
	req.Weight = math.Round(req.Weight*1000) / 1000

	fmt.Printf("Specify the price (int)\n> ")
	input, err = reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}
	req.Price, err = strconv.ParseInt(strings.TrimSpace(input), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid format, err: %v", err)
	}
	fmt.Printf("Specify the type of packaging \nbag - 5 units\nbox - 20 units\nfilm - 1 unit\n> ")
	req.Packaging, err = reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}
	req.Packaging = strings.ToLower(strings.TrimSpace(req.Packaging))

	packaging, err := GetPackaging(req.Packaging)
	if err != nil {
		return nil, fmt.Errorf("failed to get packaging: %v", err)
	}

	if !packaging.CheckWeight(req.Weight) {
		return nil, fmt.Errorf("weight is too big for the selected packaging")
	}

	req.Price += packaging.GetPrice()
	if req.Packaging != "film" {
		fmt.Printf("Would you like to add additional film for 1 units? (yes or no)\n> ")
		req.AdditionalFilm, err = reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading input: %v", err)
		}
		req.AdditionalFilm = strings.ToLower(strings.TrimSpace(req.AdditionalFilm))
		switch req.AdditionalFilm {
		case "yes":
			req.Price += 1
		case "no":
		default:
			return nil, fmt.Errorf("invalid format (you can choose between yes or no)")
		}
	}

	err = s.storage.AddOrder(ctx, req.OrderID, req.ClientID, req.CreatedAt, req.ExpiredAt, req.Weight, req.Price, req.Packaging, req.AdditionalFilm)
	if err != nil {
		return nil, fmt.Errorf("failed to add order: %v", err)
	}

	fmt.Printf("Final price for order %d: %d\n", req.OrderID, req.Price)
	return &AddOrderResponse{}, nil
}
