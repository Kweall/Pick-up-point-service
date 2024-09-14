package packaging

type Film struct{}

func (f *Film) NewPrice(weight float64, oldPrice int64) (int64, error) {
	return oldPrice + 1, nil
}
