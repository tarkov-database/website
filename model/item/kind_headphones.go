package item

const (
	KindHeadphone Kind = "headphone"
)

type Headphone struct {
	Item

	AmbientVolume  float64    `json:"ambientVol"`
	DryVolume      float64    `json:"dryVol"`
	Distortion     float64    `json:"distortion"`
	HighPassFilter HighPass   `json:"hpf"`
	Compressor     Compressor `json:"compressor"`
}

type HighPass struct {
	CutoffFrequency float64 `json:"cutoffFreq"`
	Resonance       float64 `json:"resonance"`
}

type Compressor struct {
	Attack    float64 `json:"attack"`
	Gain      float64 `json:"gain"`
	Release   float64 `json:"release"`
	Treshhold float64 `json:"treshhold"`
	Volume    float64 `json:"volume"`
}

type HeadphoneResult struct {
	*Result
	Items []Headphone `json:"items"`
}

func (r *HeadphoneResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
