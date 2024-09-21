package commands

import (
	"fmt"
	"strconv"
)

const countOfArgumentstoReturn = 3

func AcceptReturn(storage Storage, parts []string) error {
	if len(parts) != countOfArgumentstoReturn {
		return fmt.Errorf("should be 2 arguments: clientID (int), orderID (int)")
	}

	clientID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("clientID is incorrect")
	}

	orderID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("orderID is incorrect")
	}
	err = storage.AcceptReturn(clientID, orderID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Printf("Return accepted for Order ID: %d, ClientID: %d\n", orderID, clientID)
	return nil
}
