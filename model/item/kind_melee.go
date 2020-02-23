package item

const (
	KindMelee Kind = "melee"
)

type Melee struct {
	Item `bson:",inline"`

	Slash MeleeAttack `json:"slash"`
	Stab  MeleeAttack `json:"stab"`
}

type MeleeAttack struct {
	Damage      float64 `json:"damage"`
	Rate        float64 `json:"rate"`
	Range       float64 `json:"range"`
	Consumption float64 `json:"consumption"`
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
