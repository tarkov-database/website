package item

const (
	KindArmor Kind = "armor"
)

type Armor struct {
	Item `bson:",inline"`

	Type          string    `json:"type"`
	Armor         ArmorProp `json:"armor"`
	Penalties     Penalties `json:"penalties"`
	Blocking      []string  `json:"blocking"`
	Slots         Slots     `json:"slots"`
	Compatibility ItemList  `json:"compatibility"`
}

type ArmorResult struct {
	Count int64   `json:"total"`
	Items []Armor `json:"items"`
}

func (r *ArmorResult) GetCount() int64 {
	return r.Count
}

func (r *ArmorResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type ArmorProp struct {
	Class           int64         `json:"class"`
	Durability      float64       `json:"durability"`
	Material        ArmorMaterial `json:"material"`
	BluntThroughput float64       `json:"bluntThroughput"`
	Zones           []string      `json:"zones"`
}

type ArmorMaterial struct {
	Name            string  `json:"name"`
	Destructibility float64 `json:"destructibility"`
}
