package feature

import (
	"context"
	"fmt"
	"time"

	"github.com/tarkov-database/website/core/api"
)

// Group ...
type Group struct {
	ID          objectID  `json:"_id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Tags        []string  `json:"tags" bson:"tags"`
	Location    objectID  `json:"_location" bson:"_location"`
	Modified    timestamp `json:"_modified" bson:"_modified"`
}

type GroupResult struct {
	Count int64   `json:"total"`
	Items []Group `json:"items"`
}

// const defaultSort = "name"

func GetGroup(lID, fID objectID) (*Group, error) {
	loc := &Group{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/location/%s/featuregroup/%s", lID, fID), &api.Options{}, loc); err != nil {
		return loc, err
	}

	return loc, nil
}

func GetGroups(lID objectID, opts *api.Options) (*GroupResult, error) {
	result := &GroupResult{}

	if opts.Sort == "" {
		opts.Sort = defaultSort
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/location/%s/featuregroup", lID), opts, result); err != nil {
		return result, err
	}

	return result, nil
}
