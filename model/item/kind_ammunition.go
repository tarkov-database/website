package item

const (
	KindAmmunition Kind = "ammunition"
)

type Ammunition struct {
	Item

	Caliber             string                `json:"caliber"`
	Type                string                `json:"type"`
	Tracer              bool                  `json:"tracer"`
	TracerColor         string                `json:"tracerColor"`
	Subsonic            bool                  `json:"subsonic"`
	CasingMass          float64               `json:"casingMass"`
	BulletMass          float64               `json:"bulletMass"`
	BulletDiameter      float64               `json:"bulletDiameter"`
	Velocity            float64               `json:"velocity"`
	BallisticCoeficient float64               `json:"ballisticCoef"`
	Damage              float64               `json:"damage"`
	Penetration         float64               `json:"penetration"`
	ArmorDamage         float64               `json:"armorDamage"`
	Fragmentation       AmmoFrag              `json:"fragmentation"`
	Effects             AmmoEffects           `json:"effects"`
	Projectiles         int64                 `json:"projectiles"`
	MisfireChance       float64               `json:"misfireChance"`
	FailureToFeedChance float64               `json:"failureToFeedChance"`
	WeaponModifier      WeaponModifier        `json:"weaponModifier"`
	GrenadeProperties   AmmoGrenadeProperties `json:"grenadeProps,omitempty"`
}

func (b *Ammunition) Heat() float64 {
	return (b.WeaponModifier.HeatFactor - 1) * 100
}

func (b *Ammunition) DurabilityBurnPercent() float64 {
	return (b.WeaponModifier.DurabilityBurn - 1) * 100
}

type AmmunitionResult struct {
	*Result
	Items []Ammunition `json:"items"`
}

func (r *AmmunitionResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type AmmoFrag struct {
	Chance float64 `json:"chance"`
	Min    int64   `json:"min"`
	Max    int64   `json:"max"`
}

type AmmoEffects struct {
	LightBleedingChance float64 `json:"lightBleedingChance,omitempty"`
	HeavyBleedingChance float64 `json:"heavyBleedingChance,omitempty"`
}

type WeaponModifier struct {
	Accuracy       float64 `json:"accuracy"`
	Recoil         float64 `json:"recoil"`
	DurabilityBurn float64 `json:"durabilityBurn"`
	HeatFactor     float64 `json:"heatFactor"`
}

type AmmoGrenadeProperties struct {
	Delay         float64 `json:"delay"`
	FragmentCount float64 `json:"fragCount"`
	MinRadius     float64 `json:"minRadius"`
	MaxRadius     float64 `json:"maxRadius"`
}

var calibers = [...]string{
	".366 TKM",
	"11.43x23mm ACP",
	"12.7x108mm",
	"12.7x55mm STs-130",
	"12ga",
	"20ga",
	"23x75mmR",
	"26.5x75mm",
	"30x29mm",
	"40mmRU",
	"40x46mm",
	"5.45x39mm",
	"5.56x45mm NATO",
	"5.7x28mm",
	"6.8x51mm",
	"7.62x25mm Tokarev",
	"7.62x35mm",
	"7.62x39mm",
	"7.62x51mm NATO",
	"7.62x54mmR",
	"8.6x70mm",
	"9x18mm Makarov",
	"9x19mm Parabellum",
	"9x21mm Gyurza",
	"9x33mmR",
	"9x39mm",
	"HK 4.6x30mm",
}

var ammunitionFilter = Filter{
	"caliber": calibers[:],
}
