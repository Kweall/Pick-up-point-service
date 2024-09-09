package storage

import "homework/storage/json_file"

type Storage interface {
	AddOrder(*json_file.Order) (err error)
	DeleteOrderByID(int64) error
	GiveOrdersToClient(skuIDs []int64) error
}
