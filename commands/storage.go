package commands

import "homework/storage/json_file"

type Storage interface {
	AddOrder(*json_file.Order) (err error)
	AddOrderToStory(orderID int64, path string) error
	DeleteOrderByID(int64) error
	GiveOrdersToClient(skuIDs []int64) error
	GetAll() (map[int64]*json_file.Order, error)
	AcceptReturn(clientID, orderID int64) error
	CheckIfExists(orderID int64) error
}
