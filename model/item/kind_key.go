package item

const (
	KindKey Kind = "key"
)

type Key struct {
	Item

	Location string `json:"location"`
	Usages   int64  `json:"usages,omitempty"`
}

type KeyResult struct {
	*Result
	Items []Key `json:"items"`
}

func (r *KeyResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
