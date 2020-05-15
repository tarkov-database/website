package item

const (
	KindMap Kind = "map"
)

type Map struct {
	Item `bson:",inline"`
}

type MapResult struct {
	*Result
	Items []Map `json:"items"`
}

func (r *MapResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
