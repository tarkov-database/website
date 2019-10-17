package item

const (
	KindHeadphone Kind = "headphone"
)

type Headphone struct {
	Item `bson:",inline"`

	AmbientVolume  float64    `json:"ambientVol" bson:"ambientVol"`
	DryVolume      float64    `json:"dryVol" bson:"dryVol"`
	Distortion     float64    `json:"distortion" bson:"distortion"`
	HighPassFilter HighPass   `json:"hpf" bson:"hpf"`
	Compressor     Compressor `json:"compressor" bson:"compressor"`
}

type HighPass struct {
	CutoffFrequency float64 `json:"cutoffFreq" bson:"cutoffFreq"`
	Resonance       float64 `json:"resonance" bson:"resonance"`
}

type Compressor struct {
	Attack    float64 `json:"attack" bson:"attack"`
	Gain      float64 `json:"gain" bson:"gain"`
	Release   float64 `json:"release" bson:"release"`
	Treshhold float64 `json:"treshhold" bson:"treshhold"`
	Volume    float64 `json:"volume" bson:"volume"`
}

type HeadphoneResult struct {
	Count int64       `json:"total"`
	Items []Headphone `json:"items"`
}

func (r *HeadphoneResult) GetCount() int64 {
	return r.Count
}

func (r *HeadphoneResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}
