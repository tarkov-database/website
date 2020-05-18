package item

const (
	KindBackpack Kind = "backpack"
)

type Backpack struct {
	Item

	Grids     []Grid    `json:"grids"`
	Penalties Penalties `json:"penalties"`
}

type BackpackResult struct {
	*Result
	Items []Backpack `json:"items"`
}

func (r *BackpackResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
