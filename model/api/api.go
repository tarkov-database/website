package api

import (
	"time"

	"github.com/tarkov-database/website/model/item"
)

type API struct {
	Modified time.Time `json:"modified"`
	Item     *Endpoint `json:"item"`
}

type Endpoint struct {
	Count    int64     `json:"count"`
	Modified time.Time `json:"modified"`
}

func GetAPI() (*API, error) {
	a := &API{}

	ii, err := item.GetIndex(true)
	if err != nil {
		return a, err
	}

	a.Modified = ii.Modified.Time

	a.Item = &Endpoint{
		Count:    ii.Total,
		Modified: a.Modified,
	}

	return a, nil
}
