package location

import (
	"context"
	"fmt"
	"time"

	"github.com/tarkov-database/website/core/api"
)

type objectID = string

type timestamp = api.Timestamp

// Location describes the entity of a location
type Location struct {
	ID             objectID  `json:"_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	MinimumPlayers int64     `json:"minPlayers"`
	MaximumPlayers int64     `json:"maxPlayers"`
	EscapeTime     int64     `json:"escapeTime"`
	Insurance      bool      `json:"insurance"`
	Available      bool      `json:"available"`
	Exits          []Exit    `json:"exits"`
	Bosses         []Boss    `json:"bosses"`
	Modified       timestamp `json:"_modified"`
}

type LocationResult struct {
	Count int64      `json:"total"`
	Items []Location `json:"items"`
}

// Exit describes an exit of a location
type Exit struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Chance           float64 `json:"chance"`
	MinimumTime      int64   `json:"minTime"`
	MaximumTime      int64   `json:"maxTime"`
	ExfiltrationTime int64   `json:"exfilTime"`
	Requirement      string  `json:"requirement,omitempty"`
}

// Boss describes a boss of a location
type Boss struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Chance      float64 `json:"chance"`
	Followers   int64   `json:"followers"`
}

const defaultSort = "name"

func GetLocation(id objectID) (*Location, error) {
	loc := &Location{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/location/%s", id), &api.Options{}, loc); err != nil {
		return loc, err
	}

	return loc, nil
}

func GetLocations(opts *api.Options) (*LocationResult, error) {
	result := &LocationResult{}

	if opts.Sort == "" {
		opts.Sort = defaultSort
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, "/location", opts, result); err != nil {
		return result, err
	}

	return result, nil
}

func GetLocationsByText(txt string, limit int) (*LocationResult, error) {
	opts := &api.Options{
		Sort:  "",
		Limit: limit,
		Filter: map[string]string{
			"text": txt,
		},
	}

	result := &LocationResult{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, "/location", opts, result); err != nil {
		return result, err
	}

	return result, nil
}

type Filter map[string][]string

func (f Filter) GetAll() map[string][]string {
	return f
}

func (f Filter) Get(k string) []string {
	return f[k]
}

func GetFilter() Filter {
	return Filter{
		"available": {
			"true",
			"false",
		},
	}
}
