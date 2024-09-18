package packaging

type Film struct{}

func (f *Film) CheckWeight(weight float64) bool {
	return true
}

func (f *Film) GetPrice() int64 {
	return 1
}
