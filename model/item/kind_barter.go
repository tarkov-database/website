package item

const (
	KindBarter Kind = "barter"
)

type Barter struct {
	Item
}

type BarterResult struct {
	*Result
	Items []Barter `json:"items"`
}

func (r *BarterResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
