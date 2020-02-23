package item

const (
	KindFirearm Kind = "firearm"
)

type Firearm struct {
	Item `bson:",inline"`

	Type              string   `json:"type" bson:"type"`
	Class             string   `json:"class" bson:"class"`
	Caliber           string   `json:"caliber" bson:"caliber"`
	RateOfFire        int64    `json:"rof" bson:"rof"`
	Action            string   `json:"action" bson:"action"`
	Modes             []string `json:"modes" bson:"modes"`
	Velocity          float64  `json:"velocity" bson:"velocity"`
	EffectiveDistance int64    `json:"effectiveDist" bson:"effectiveDist"`
	Ergonomics        float64  `json:"ergonomicsFP" bson:"ergonomicsFP"`
	FoldRectractable  bool     `json:"foldRectractable" bson:"foldRectractable"`
	RecoilVertical    int64    `json:"recoilVertical" bson:"recoilVertical"`
	RecoilHorizontal  int64    `json:"recoilHorizontal" bson:"recoilHorizontal"`
	Slots             Slots    `json:"slots" bson:"slots"`
}

type FirearmResult struct {
	Count int64     `json:"total"`
	Items []Firearm `json:"items"`
}

func (r *FirearmResult) GetCount() int64 {
	return r.Count
}

func (r *FirearmResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
