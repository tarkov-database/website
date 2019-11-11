package item

const (
	KindAmmunition Kind = "ammunition"
)

type Ammunition struct {
	Item `bson:",inline"`

	Caliber             string         `json:"caliber" bson:"caliber"`
	Type                string         `json:"type" bson:"type"`
	Tracer              bool           `json:"tracer" bson:"tracer"`
	TracerColor         string         `json:"tracerColor" bson:"tracerColor"`
	Subsonic            bool           `json:"subsonic" bson:"subsonic"`
	Velocity            float64        `json:"velocity" bson:"velocity"`
	BallisticCoeficient float64        `json:"ballisticCoef" bson:"ballisticCoef"`
	Damage              float64        `json:"damage" bson:"damage"`
	Penetration         float64        `json:"penetration" bson:"penetration"`
	ArmorDamage         float64        `json:"armorDamage" bson:"armorDamage"`
	Fragmentation       AmmoFrag       `json:"fragmentation" bson:"fragmentation"`
	Projectiles         int64          `json:"projectiles" bson:"projectiles"`
	WeaponModifier      WeaponModifier `json:"weaponModifier" bson:"weaponModifier"`
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
	Chance float64 `json:"chance" bson:"chance"`
	Min    int64   `json:"min" bson:"min"`
	Max    int64   `json:"max" bson:"max"`
}

type WeaponModifier struct {
	Accuracy float64 `json:"accuracy" bson:"accuracy"`
	Recoil   float64 `json:"recoil" bson:"recoil"`
}
