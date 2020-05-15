package item

const (
	KindMedical Kind = "medical"
)

type Medical struct {
	Item `bson:",inline"`

	Type         string  `json:"type"`
	Resources    int64   `json:"resources"`
	ResourceRate int64   `json:"resourceRate"`
	UseTime      float64 `json:"useTime"`
	Effects      Effects `json:"effects"`
}

type MedicalResult struct {
	*Result
	Items []Medical `json:"items"`
}

func (r *MedicalResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var medicalFilter = Filter{
	"type": {
		"accessory",
		"drug",
		"medkit",
		"stimulator",
	},
}
