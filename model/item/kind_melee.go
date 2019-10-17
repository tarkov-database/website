package item

const (
	KindMelee Kind = "melee"
)

type Melee struct {
	Item `bson:",inline"`

	Slash MeleeAttack `json:"slash" bson:"slash"`
	Stab  MeleeAttack `json:"stab" bson:"stab"`
}

type MeleeAttack struct {
	Damage      float64 `json:"damage" bson:"damage"`
	Rate        float64 `json:"rate" bson:"rate"`
	Range       float64 `json:"range" bson:"range"`
	Consumption float64 `json:"consumption" bson:"consumption"`
}

type MeleeResult struct {
	Count int64   `json:"total"`
	Items []Melee `json:"items"`
}

func (r *MeleeResult) GetCount() int64 {
	return r.Count
}

func (r *MeleeResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
