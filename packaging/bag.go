package packaging

type Bag struct{}

func (b *Bag) CheckWeight(weight float64) bool {
	return weight < 10
}

func (b *Bag) GetPrice() int64 {
	return 5
}
