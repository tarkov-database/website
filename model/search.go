package model

import (
	"sync"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model/item"
	"github.com/tarkov-database/website/model/location"
	"github.com/tarkov-database/website/model/location/feature"
)

type SearchResult struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	ShortName string     `json:"shortName,omitempty"`
	Parent    string     `json:"parent,omitempty"`
	Type      EntityType `json:"type"`
}

type SearchOperation struct {
	Term    string
	Filter  *SearchFilter
	Limit   int
	Results chan []*SearchResult
	Error   error
	Tasks   sync.WaitGroup
	sync.RWMutex
}

type SearchFilter struct {
	Category string
	Location string
}

func NewSearch(term string, filter *SearchFilter, limit int) *SearchOperation {
	return &SearchOperation{
		Term:    term,
		Filter:  filter,
		Limit:   limit,
		Results: make(chan []*SearchResult, 1),
	}
}

func (so *SearchOperation) Close() {
	go func() {
		so.Tasks.Wait()
		close(so.Results)
	}()
}

func (so *SearchOperation) Items() {
	defer so.Tasks.Done()

	var err error
	var result item.EntityResult

	if c := so.Filter.Category; c != "" {
		var k item.Kind
		if k, err = item.CategoryToKind(c); err != nil {
			so.Lock()
			so.Error = err
			so.Unlock()
			return
		}

		result, err = item.GetItems(k, &api.Options{
			Filter: map[string]string{"text": so.Term},
			Limit:  so.Limit,
		})
	} else {
		result, err = item.GetItemsBySearch(so.Term, so.Limit)
	}
	if err != nil {
		so.Lock()
		so.Error = err
		so.Unlock()
		return
	}

	items := result.GetEntities()

	rs := make([]*SearchResult, len(items))
	for i, r := range items {
		cat, err := item.KindToCategory(r.GetKind())
		if err != nil {
			so.Lock()
			so.Error = err
			so.Unlock()
			return
		}

		rs[i] = &SearchResult{
			ID:        r.GetID(),
			Name:      r.GetName(),
			ShortName: r.GetShortName(),
			Parent:    cat,
			Type:      TypeItem,
		}
	}

	so.Results <- rs
}

func (so *SearchOperation) Locations() {
	defer so.Tasks.Done()

	result, err := location.GetLocationsByText(so.Term, so.Limit)
	if err != nil {
		so.Lock()
		so.Error = err
		so.Unlock()
		return
	}

	items := result.Items

	rs := make([]*SearchResult, len(items))
	for i, r := range items {
		rs[i] = &SearchResult{
			ID:   r.ID,
			Name: r.Name,
			Type: TypeLocation,
		}
	}

	so.Results <- rs
}

func (so *SearchOperation) Features() {
	defer so.Tasks.Done()

	result, err := feature.GetFeaturesByText(so.Term, so.Filter.Location, so.Limit)
	if err != nil {
		so.Lock()
		so.Error = err
		so.Unlock()
		return
	}

	items := result.Items

	rs := make([]*SearchResult, len(items))
	for i, r := range items {
		rs[i] = &SearchResult{
			ID:     r.ID,
			Name:   r.Name,
			Parent: r.Group,
			Type:   TypeFeature,
		}
	}

	so.Results <- rs
}
