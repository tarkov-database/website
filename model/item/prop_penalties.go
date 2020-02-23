package item

type Penalties struct {
	Mouse      float64 `json:"mouse,omitempty"`
	Speed      float64 `json:"speed,omitempty"`
	Ergonomics float64 `json:"ergonomicsFP,omitempty"`
	Deafness   string  `json:"deafness,omitempty"`
}
