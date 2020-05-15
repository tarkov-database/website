package item

const (
	KindFirearm Kind = "firearm"
)

type Firearm struct {
	Item `bson:",inline"`

	Type              string   `json:"type"`
	Class             string   `json:"class"`
	Caliber           string   `json:"caliber"`
	RateOfFire        int64    `json:"rof"`
	Action            string   `json:"action"`
	Modes             []string `json:"modes"`
	Velocity          float64  `json:"velocity"`
	EffectiveDistance int64    `json:"effectiveDist"`
	Ergonomics        float64  `json:"ergonomicsFP"`
	FoldRectractable  bool     `json:"foldRectractable"`
	RecoilVertical    int64    `json:"recoilVertical"`
	RecoilHorizontal  int64    `json:"recoilHorizontal"`
	Slots             Slots    `json:"slots"`
}

type FirearmResult struct {
	*Result
	Items []Firearm `json:"items"`
}

func (r *FirearmResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var firearmFilter = Filter{
	"type": {
		"primary",
		"secondary",
	},
	"class": {
		"assaultCarbine",
		"assaultRifle",
		"grenadeLauncher",
		"machinegun",
		"marksmanRifle",
		"pistol",
		"shotgun",
		"smg",
		"sniperRifle",
	},
	"caliber": calibers[:],
}
