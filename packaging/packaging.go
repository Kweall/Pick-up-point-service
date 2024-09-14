package packaging

import "fmt"

type Packaging interface {
	NewPrice(weight float64, oldPrice int64) (int64, error)
}

func GetPackaging(packagingType string) (Packaging, error) {
	switch packagingType {
	case "bag":
		return &Bag{}, nil
	case "box":
		return &Box{}, nil
	case "film":
		return &Film{}, nil
	default:
		return nil, fmt.Errorf("unknown packaging type")
	}
}
