package item

const (
	KindTacticalrig Kind = "tacticalrig"
)

type TacticalRig struct {
	Item

	Capacity        int64            `json:"capacity"`
	Grids           []Grid           `json:"grids"`
	Penalties       Penalties        `json:"penalties"`
	ArmorComponents []ArmorComponent `json:"armorComponents,omitempty" bson:"armorComponents,omitempty"`
	IsPlateCarrier  bool             `json:"isPlateCarrier" bson:"isPlateCarrier"`
	Slots           Slots            `json:"slots" bson:"slots"`
}

func (t TacticalRig) TotalDurability() float64 {
	var total float64
	for _, component := range t.ArmorComponents {
		total += component.Durability
	}

	return total
}

type TacticalRigResult struct {
	*Result
	Items []TacticalRig `json:"items"`
}

func (r *TacticalRigResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var tacticalRigFilter = Filter{
	"armored": {
		"true",
		"false",
	},
	"plateCarrier": {
		"true",
		"false",
	},
	"class":    armorClasses[:],
	"material": armorMaterial[:],
}
