package item

import (
	"context"
	"time"

	"github.com/tarkov-database/website/core/api"
)

type Index struct {
	Total    int64                 `json:"total"`
	Modified timestamp             `json:"modified"`
	Kinds    map[string]*KindStats `json:"kinds"`
}

type KindStats struct {
	Count    int64     `json:"count"`
	Modified timestamp `json:"modified"`
}

func GetIndex(skipKinds bool) (*Index, error) {
	idx := &Index{}

	skip := "0"
	if skipKinds {
		skip = "1"
	}

	opts := &api.Options{
		Filter: map[string]string{
			"skipKinds": skip,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.GET(ctx, "/item", opts, idx); err != nil {
		return idx, err
	}

	return idx, nil
}
