package item

const (
	KindMap Kind = "map"
)

type Map struct {
	Item `bson:",inline"`
}

type MapResult struct {
	Count int64 `json:"total"`
	Items []Map `json:"items"`
}

func (r *MapResult) GetCount() int64 {
	return r.Count
}

func (r *MapResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
