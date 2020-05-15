package item

const (
	KindClothing Kind = "clothing"
)

type Clothing struct {
	Item `bson:",inline"`

	Type      string    `json:"type"`
	Blocking  []string  `json:"blocking"`
	Penalties Penalties `json:"penalties"`
	Slots     Slots     `json:"slots"`
}

type ClothingResult struct {
	*Result
	Items []Clothing `json:"items"`
}

func (r *ClothingResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var clothingFilter = Filter{
	"type": {
		"eyewear",
		"faceCover",
		"headwear",
		"unknown",
	},
}
