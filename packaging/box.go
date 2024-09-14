package packaging

import "fmt"

type Box struct{}

func (b *Box) NewPrice(weight float64, oldPrice int64) (int64, error) {
	if weight >= 30 {
		return 0, fmt.Errorf("weight exceeds 10 kg for bag")
	}
	return oldPrice + 20, nil
}
