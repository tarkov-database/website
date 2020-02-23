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
	Count int64        `json:"total"`
	Items []Ammunition `json:"items"`
}

func (r *AmmunitionResult) GetCount() int64 {
	return r.Count
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
