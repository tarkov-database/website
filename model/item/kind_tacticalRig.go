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

func (t TacticalRig) ClassRange() ClassRange {
	if len(t.ArmorComponents) == 0 {
		return ClassRange{}
	}

	min, max := t.ArmorComponents[0].Class, t.ArmorComponents[0].Class
	for _, component := range t.ArmorComponents {
		if component.Class < min {
			min = component.Class
		}
		if component.Class > max {
			max = component.Class
		}
	}

	return ClassRange{Min: min, Max: max}
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
