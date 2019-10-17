package item

const (
	KindMoney Kind = "money"
)

type Money struct {
	Item `bson:",inline"`
}

type MoneyResult struct {
	Count int64   `json:"total"`
	Items []Money `json:"items"`
}

func (r *MoneyResult) GetCount() int64 {
	return r.Count
}

func (r *MoneyResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
