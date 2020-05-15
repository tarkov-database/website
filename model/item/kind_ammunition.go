package item

const (
	KindAmmunition Kind = "ammunition"
)

type Ammunition struct {
	Item `bson:",inline"`

	Caliber             string         `json:"caliber"`
	Type                string         `json:"type"`
	Tracer              bool           `json:"tracer"`
	TracerColor         string         `json:"tracerColor"`
	Subsonic            bool           `json:"subsonic"`
	Velocity            float64        `json:"velocity"`
	BallisticCoeficient float64        `json:"ballisticCoef"`
	Damage              float64        `json:"damage"`
	Penetration         float64        `json:"penetration"`
	ArmorDamage         float64        `json:"armorDamage"`
	Fragmentation       AmmoFrag       `json:"fragmentation"`
	Projectiles         int64          `json:"projectiles"`
	WeaponModifier      WeaponModifier `json:"weaponModifier"`
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

type WeaponModifier struct {
	Accuracy float64 `json:"accuracy"`
	Recoil   float64 `json:"recoil"`
}

var calibers = [...]string{
	".366 TKM",
	"11.43x23mm ACP",
	"12.7x108mm",
	"12.7x55mm STs-130",
	"12ga",
	"20ga",
	"30x29mm",
	"5.45x39mm",
	"5.56x45mm NATO",
	"5.7x28mm",
	"7.62x25mm Tokarev",
	"7.62x39mm",
	"7.62x51mm NATO",
	"7.62x54mmR",
	"9x18mm Makarov",
	"9x19mm Parabellum",
	"9x21mm Gyurza",
	"9x39mm",
	"HK 4.6x30mm",
}

var ammunitionFilter = Filter{
	"caliber": calibers[:],
}
