package item

const (
	KindTacticalrig Kind = "tacticalrig"
)

type TacticalRig struct {
	Item

	Capacity  int64     `json:"capacity"`
	Grids     []Grid    `json:"grids"`
	Penalties Penalties `json:"penalties"`
	Armor     ArmorProp `json:"armor,omitempty"`
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

var tacticalFilter = Filter{
	"armored": {
		"true",
		"false",
	},
	"class":    armorClasses[:],
	"material": armorMaterial[:],
}
