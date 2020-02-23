package item

const (
	KindKey Kind = "key"
)

type Key struct {
	Item `bson:",inline"`

	Location string `json:"location"`
}

type KeyResult struct {
	Count int64 `json:"total"`
	Items []Key `json:"items"`
}

func (r *KeyResult) GetCount() int64 {
	return r.Count
}

func (r *KeyResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
