package item

const (
	KindMoney Kind = "money"
)

type Money struct {
	Item
}

type MoneyResult struct {
	*Result
	Items []Money `json:"items"`
}

func (r *MoneyResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
