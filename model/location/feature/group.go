package feature

import (
	"context"
	"fmt"
	"time"

	"github.com/tarkov-database/website/core/api"
)

// Group ...
type Group struct {
	ID          objectID  `json:"_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Location    objectID  `json:"_location"`
	Modified    timestamp `json:"_modified"`
}

type GroupResult struct {
	Count int64   `json:"total"`
	Items []Group `json:"items"`
}

// const defaultSort = "name"

func GetGroup(lID, gID objectID) (*Group, error) {
	group := &Group{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/location/%s/featuregroup/%s", lID, gID), &api.Options{}, group); err != nil {
		return group, err
	}

	return group, nil
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
