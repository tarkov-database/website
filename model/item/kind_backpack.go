package item

const (
	KindBackpack Kind = "backpack"
)

type Backpack struct {
	Item `bson:",inline"`

	Grids     []Grid    `json:"grids" bson:"grids"`
	Penalties Penalties `json:"penalties" bson:"penalties"`
}

type BackpackResult struct {
	Count int64      `json:"total"`
	Items []Backpack `json:"items"`
}

func (r *BackpackResult) GetCount() int64 {
	return r.Count
}

func (r *BackpackResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
