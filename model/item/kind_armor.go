package item

const (
	KindArmor Kind = "armor"
)

type Armor struct {
	Item `bson:",inline"`

	Type          string    `json:"type" bson:"type"`
	Armor         ArmorProp `json:"armor" bson:"armor"`
	Penalties     Penalties `json:"penalties" bson:"penalties"`
	Blocking      []string  `json:"blocking" bson:"blocking"`
	Slots         Slots     `json:"slots" bson:"slots"`
	Compatibility ItemList  `json:"compatibility" bson:"compatibility"`
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
	Class           int64         `json:"class" bson:"class"`
	Durability      float64       `json:"durability" bson:"durability"`
	Material        ArmorMaterial `json:"material" bson:"material"`
	BluntThroughput float64       `json:"bluntThroughput" bson:"bluntThroughput"`
	Zones           []string      `json:"zones" bson:"zones"`
}

type ArmorMaterial struct {
	Name            string  `json:"name" bson:"name"`
	Destructibility float64 `json:"destructibility" bson:"destructibility"`
}
