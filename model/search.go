package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tarkov-database/website/core/search"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := &search.Query{
		Term:  so.Term,
		Index: search.IndexItem,
	}

	if c := so.Filter.Category; c != "" {
		k, err := item.CategoryToKind(c)
		if err != nil {
			so.Lock()
			so.Error = err
			so.Unlock()
			return
		}

		query.Term = fmt.Sprintf("kind:%s AND %s", k, query.Term)
	}

	opts := &search.Options{
		Limit: so.Limit,
	}

	result, err := search.Search(ctx, query, opts)
	if err != nil {
		so.Lock()
		so.Error = err
		so.Unlock()
		return
	}

	rs := make([]*SearchResult, len(result.Data))
	for i, r := range result.Data {
		cat, err := item.KindToCategory(item.Kind(r.Kind))
		if err != nil {
			so.Lock()
			so.Error = err
			so.Unlock()
			return
		}

		rs[i] = &SearchResult{
			ID:        r.ID,
			Name:      r.Name,
			ShortName: r.ShortName,
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
