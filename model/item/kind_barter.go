package item

const (
	KindBarter Kind = "barter"
)

type Barter struct {
	Item `bson:",inline"`
}

type BarterResult struct {
	Count int64    `json:"total"`
	Items []Barter `json:"items"`
}

func (r *BarterResult) GetCount() int64 {
	return r.Count
}

func (r *BarterResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
