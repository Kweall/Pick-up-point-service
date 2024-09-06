package storage

type Storage interface {
	AddItemToCart(userID, skuID int64, date string) (err error)
	CheckIfItemExists(userID, skuID int64, date string) (bool, error)
	DeleteItemBySkuID(skuID int64) error
	GiveOrdersToClient(skuIDs []int64) error
}
