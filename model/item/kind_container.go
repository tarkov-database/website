package item

const (
	KindContainer Kind = "container"
)

type Container struct {
	Item

	Grids []Grid `json:"grids"`
}

type ContainerResult struct {
	*Result
	Items []Container `json:"items"`
}

func (r *ContainerResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
