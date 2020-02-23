package item

type Effects struct {
	Energy            Effect `json:"energy,omitempty"`
	Hydration         Effect `json:"hydration,omitempty"`
	Bloodloss         Effect `json:"bloodloss,omitempty"`
	Fracture          Effect `json:"fracture,omitempty"`
	Contusion         Effect `json:"contusion,omitempty"`
	Pain              Effect `json:"pain,omitempty"`
	Toxication        Effect `json:"toxication,omitempty"`
	RadiationExposure Effect `json:"radExposure,omitempty"`
	Mobility          Effect `json:"mobility,omitempty"`
	Recoil            Effect `json:"recoil,omitempty"`
	ReloadSpeed       Effect `json:"reloadSpeed,omitempty"`
	LootSpeed         Effect `json:"lootSpeed,omitempty"`
	UnlockSpeed       Effect `json:"unlockSpeed,omitempty"`
}

type Effect struct {
	ResourceCosts int64   `json:"resourceCosts"`
	FadeIn        float64 `json:"fadeIn"`
	FadeOut       float64 `json:"fadeOut"`
	Duration      float64 `json:"duration"`
	Value         float64 `json:"value"`
	IsPercent     bool    `json:"isPercent"`
	Removes       bool    `json:"removes"`
}
