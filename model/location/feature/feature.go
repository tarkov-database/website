package feature

import (
	"context"
	"fmt"
	"time"

	"github.com/tarkov-database/website/core/api"
)

type objectID = string

type timestamp = api.Timestamp

// Feature ...
type Feature struct {
	ID         objectID               `json:"_id"`
	Name       string                 `json:"name"`
	Type       FeatureType            `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
	Group      objectID               `json:"group"`
	Location   objectID               `json:"_location"`
	Modified   timestamp              `json:"_modified"`
}

type FeatureResult struct {
	Count int64     `json:"total"`
	Items []Feature `json:"items"`
}

func (r *FeatureResult) FeatureCollection() *FeatureCollection {
	return &FeatureCollection{FeatureCollectionType, r.Items}
}

const defaultSort = "name"

func GetFeature(lID, fID objectID) (*Feature, error) {
	f := &Feature{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/location/%s/feature/%s", lID, fID), &api.Options{}, f); err != nil {
		return f, err
	}

	return f, nil
}

func GetFeatures(lID objectID, opts *api.Options) (*FeatureResult, error) {
	result := &FeatureResult{}

	if opts.Sort == "" {
		opts.Sort = defaultSort
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/location/%s/feature", lID), opts, result); err != nil {
		return result, err
	}

	return result, nil
}

func GetFeaturesByText(txt string, lID objectID, limit int) (*FeatureResult, error) {
	opts := &api.Options{
		Sort:  "",
		Limit: limit,
		Filter: map[string]string{
			"text": txt,
		},
	}

	result := &FeatureResult{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/location/%s/feature", lID), opts, result); err != nil {
		return result, err
	}

	return result, nil
}
