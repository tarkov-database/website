package item

const (
	KindFood Kind = "food"
)

type Food struct {
	Item `bson:",inline"`

	Type      string  `json:"type" bson:"type"`
	Resources int64   `json:"resources" bson:"resources"`
	UseTime   float64 `json:"useTime" bson:"useTime"`
	Effects   Effects `json:"effects" bson:"effects"`
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
