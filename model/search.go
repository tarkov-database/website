package model

import (
	"sync"

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

func NewSearch(q string, limit int) *SearchOperation {
	return &SearchOperation{
		Keyword: q,
		Limit:   limit,
		Results: make(chan []*SearchResult, 1),
	}
}

type SearchOperation struct {
	Keyword string
	Limit   int
	Results chan []*SearchResult
	Error   error
	Tasks   sync.WaitGroup
	sync.RWMutex
}

func (so *SearchOperation) Close() {
	go func() {
		so.Tasks.Wait()
		close(so.Results)
	}()
}

func (so *SearchOperation) Items() {
	defer so.Tasks.Done()

	result, err := item.GetItemsBySearch(so.Keyword, so.Limit)
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

	result, err := location.GetLocationsByText(so.Keyword, so.Limit)
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

func (so *SearchOperation) Features(lID string) {
	defer so.Tasks.Done()

	result, err := feature.GetFeaturesByText(so.Keyword, lID, so.Limit)
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
