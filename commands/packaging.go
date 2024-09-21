package commands

import (
	"fmt"
	"homework/packaging"
)

//go:generate mkdir -p ./mocks
//go:generate minimock -g -i * -o ./mocks/ -s ".mock.go"
type Packaging interface {
	CheckWeight(weight float64) bool
	GetPrice() int64
}

func GetPackaging(packagingType string) (Packaging, error) {
	switch packagingType {
	case "bag":
		return &packaging.Bag{}, nil
	case "box":
		return &packaging.Box{}, nil
	case "film":
		return &packaging.Film{}, nil
	default:
		return nil, fmt.Errorf("unknown packaging type")
	}
}
