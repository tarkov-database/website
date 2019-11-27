package api

import (
	"sync"
	"time"

	client "github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model/item"
	"github.com/tarkov-database/website/model/location"

	"github.com/google/logger"
)

type API struct {
	Modified time.Time `json:"modified"`
	Item     *Endpoint `json:"item"`
	Location *Endpoint `json:"location"`
}

type Endpoint struct {
	Count    int64     `json:"count"`
	Modified time.Time `json:"modified"`
}

func GetAPI() (*API, error) {
	a := &API{}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		a.Item = getItemEndpoint()
	}(wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		a.Location = getLocationEndpoint()
	}(wg)

	wg.Wait()

	if a.Location.Modified.After(a.Item.Modified) {
		a.Modified = a.Location.Modified
	} else {
		a.Modified = a.Item.Modified
	}

	return a, nil
}

func getItemEndpoint() *Endpoint {
	ep := &Endpoint{}

	ii, err := item.GetIndex(true)
	if err != nil {
		logger.Error(err)
		return ep
	}

	ep.Count, ep.Modified = ii.Total, ii.Modified.Time

	return ep
}

func getLocationEndpoint() *Endpoint {
	ep := &Endpoint{}

	loc, err := location.GetLocations(&client.Options{Sort: "-_modified", Limit: 1})
	if err != nil {
		logger.Error(err)
		return ep
	}

	ep.Count, ep.Modified = loc.Count, loc.Items[0].Modified.Time

	return ep
}
