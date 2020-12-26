package item

type Effects struct {
	Energy            Effect   `json:"energy,omitempty"`
	EnergyRate        Effect   `json:"energyRate,omitempty"`
	Hydration         Effect   `json:"hydration,omitempty"`
	HydrationRate     Effect   `json:"hydrationRate,omitempty"`
	Stamina           Effect   `json:"stamina,omitempty"`
	StaminaRate       Effect   `json:"staminaRate,omitempty"`
	Health            Effect   `json:"health,omitempty"`
	HealthRate        Effect   `json:"healthRate,omitempty"`
	LightBleeding     Effect   `json:"lightBleeding,omitempty"`
	HeavyBleeding     Effect   `json:"heavyBleeding,omitempty"`
	Fracture          Effect   `json:"fracture,omitempty"`
	Contusion         Effect   `json:"contusion,omitempty"`
	Pain              Effect   `json:"pain,omitempty"`
	TunnelVision      Effect   `json:"tunnelVision,omitempty"`
	Tremor            Effect   `json:"tremor,omitempty"`
	Toxication        Effect   `json:"toxication,omitempty"`
	Antidote          Effect   `json:"antidote,omitempty"`
	RadiationExposure Effect   `json:"radExposure,omitempty"`
	BodyTemperature   Effect   `json:"bodyTemperature,omitempty"`
	Mobility          Effect   `json:"mobility,omitempty"`
	Recoil            Effect   `json:"recoil,omitempty"`
	ReloadSpeed       Effect   `json:"reloadSpeed,omitempty"`
	LootSpeed         Effect   `json:"lootSpeed,omitempty"`
	UnlockSpeed       Effect   `json:"unlockSpeed,omitempty"`
	DestroyedPart     Effect   `json:"destroyedPart,omitempty"`
	WeightLimit       Effect   `json:"weightLimit,omitempty"`
	DamageModifier    Effect   `json:"damageModifier,omitempty"`
	Skill             []Effect `json:"skill,omitempty"`
}

type Effect struct {
	Name          string          `json:"name,omitempty"`
	ResourceCosts int64           `json:"resourceCosts"`
	FadeIn        float64         `json:"fadeIn"`
	FadeOut       float64         `json:"fadeOut"`
	Chance        float64         `json:"chance"`
	Delay         float64         `json:"delay"`
	Duration      float64         `json:"duration"`
	Value         float64         `json:"value"`
	IsPercent     bool            `json:"isPercent"`
	Removes       bool            `json:"removes"`
	Penalties     EffectPenalties `json:"penalties"`
}

type EffectPenalties struct {
	HealthMin float64 `json:"healthMin,omitempty"`
	HealthMax float64 `json:"healthMax,omitempty"`
}
