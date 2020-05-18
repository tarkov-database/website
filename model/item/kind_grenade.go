package item

const (
	KindGrenade Kind = "grenade"
)

type Grenade struct {
	Item

	Type              string  `json:"type"`
	Delay             float64 `json:"delay"`
	FragmentCount     float64 `json:"fragCount"`
	MinDistance       float64 `json:"minDistance"`
	MaxDistance       float64 `json:"maxDistance"`
	ContusionDistance float64 `json:"contusionDistance"`
	Strength          float64 `json:"strength"`
	EmitTime          float64 `json:"emitTime"`
}

type GrenadeResult struct {
	*Result
	Items []Grenade `json:"items"`
}

func (r *GrenadeResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var grenadeFilter = Filter{
	"type": {
		"flash",
		"frag",
		"smoke",
	},
}
