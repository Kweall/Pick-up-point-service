package packaging

type Box struct{}

func (b *Box) CheckWeight(weight float64) bool {
	return weight < 30
}

func (b *Box) GetPrice() int64 {
	return 20
}
