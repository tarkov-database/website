package item

const (
	KindFood Kind = "food"
)

type Food struct {
	Item `bson:",inline"`

	Type      string  `json:"type"`
	Resources int64   `json:"resources"`
	UseTime   float64 `json:"useTime"`
	Effects   Effects `json:"effects"`
}

type FoodResult struct {
	Count int64  `json:"total"`
	Items []Food `json:"items"`
}

func (r *FoodResult) GetCount() int64 {
	return r.Count
}

func (r *FoodResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
