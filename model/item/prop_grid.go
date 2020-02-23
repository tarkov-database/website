package item

type Grid struct {
	ID        string   `json:"id"`
	Height    int64    `json:"height"`
	Width     int64    `json:"width"`
	MaxWeight float64  `json:"maxWeight"`
	Filter    ItemList `json:"filter"`
}
