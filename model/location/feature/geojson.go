package feature

import (
	"encoding/json"
)

type FeatureType int

const (
	FeatureSingleType FeatureType = iota
	FeatureCollectionType
)

var featureTypeStrings = [...]string{
	"Feature",
	"FeatureCollection",
}

func (ft FeatureType) String() string {
	return featureTypeStrings[ft]
}

func (ft *FeatureType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ft.String())
}

type Coordinates []interface{}

// Geometry ...
type Geometry struct {
	Type        string      `json:"type"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
	Geometries  []Geometry  `json:"geometries,omitempty"`
}

type FeatureCollection struct {
	Type     FeatureType `json:"type"`
	Features []Feature   `json:"features"`
}
