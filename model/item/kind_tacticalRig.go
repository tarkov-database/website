package item

const (
	KindTacticalrig Kind = "tacticalrig"
)

type TacticalRig struct {
	Item `bson:",inline"`

	Grids     []Grid    `json:"grids" bson:"grids"`
	Penalties Penalties `json:"penalties" bson:"penalties"`
	Armor     ArmorProp `json:"armor,omitempty" bson:"armor,omitempty"`
}

type TacticalRigResult struct {
	Count int64         `json:"total"`
	Items []TacticalRig `json:"items"`
}

func (r *TacticalRigResult) GetCount() int64 {
	return r.Count
}

func (r *TacticalRigResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
