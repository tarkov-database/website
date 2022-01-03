package item

const (
	KindFirearm Kind = "firearm"
)

type Firearm struct {
	Item

	Type               string   `json:"type"`
	Class              string   `json:"class"`
	Caliber            string   `json:"caliber"`
	Manufacturer       string   `json:"manufacturer"`
	RateOfFire         int64    `json:"rof"`
	BurstRounds        int64    `json:"burstRounds,omitempty"`
	Action             string   `json:"action"`
	Modes              []string `json:"modes"`
	Velocity           float64  `json:"velocity"`
	EffectiveDistance  int64    `json:"effectiveDist"`
	Ergonomics         float64  `json:"ergonomicsFP"`
	FoldRectractable   bool     `json:"foldRectractable"`
	RecoilVertical     int64    `json:"recoilVertical"`
	RecoilHorizontal   int64    `json:"recoilHorizontal"`
	OperatingResources float64  `json:"operatingResources"`
	MalfunctionChance  float64  `json:"malfunctionChance"`
	DurabilityRatio    float64  `json:"durabilityRatio"`
	HeatFactor         float64  `json:"heatFactor"`
	HeatFactorByShot   float64  `json:"heatFactorByShot"`
	CoolFactor         float64  `json:"coolFactor"`
	CoolFactorMods     float64  `json:"coolFactorMods"`
	CenterOfImpact     float64  `json:"centerOfImpact"`
	Slots              Slots    `json:"slots"`
}

func (f *Firearm) AccuracyMoa() float64 {
	return 100 * f.CenterOfImpact / 2.9089
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
		"specialWeapon",
	},
	"caliber": calibers[:],
}
