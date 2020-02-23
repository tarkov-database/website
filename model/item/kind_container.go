package item

const (
	KindContainer Kind = "container"
)

type Container struct {
	Item `bson:",inline"`

	Grids []Grid `json:"grids"`
}

type ContainerResult struct {
	Count int64       `json:"total"`
	Items []Container `json:"items"`
}

func (r *ContainerResult) GetCount() int64 {
	return r.Count
}

func (r *ContainerResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
