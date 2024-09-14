package packaging

import "fmt"

type Bag struct{}

func (b *Bag) NewPrice(weight float64, oldPrice int64) (int64, error) {
	if weight >= 10 {
		return 0, fmt.Errorf("weight exceeds 10 kg for bag")
	}
	return oldPrice + 5, nil
}
