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
	ID             objectID  `json:"_id" bson:"_id"`
	Name           string    `json:"name" bson:"name"`
	Description    string    `json:"description" bson:"description"`
	MinimumPlayers int64     `json:"minPlayers" bson:"minPlayers"`
	MaximumPlayers int64     `json:"maxPlayers" bson:"maxPlayers"`
	EscapeTime     int64     `json:"escapeTime" bson:"escapeTime"`
	Insurance      bool      `json:"insurance" bson:"insurance"`
	Available      bool      `json:"available" bson:"available"`
	Exits          []Exit    `json:"exits" bson:"exits"`
	Bosses         []Boss    `json:"bosses" bson:"bosses"`
	Modified       timestamp `json:"_modified" bson:"_modified"`
}

type LocationResult struct {
	Count int64      `json:"total"`
	Items []Location `json:"items"`
}

// Exit describes an exit of a location
type Exit struct {
	Name             string  `json:"name" bson:"name"`
	Description      string  `json:"description" bson:"description"`
	Chance           float64 `json:"chance" bson:"chance"`
	MinimumTime      int64   `json:"minTime" bson:"minTime"`
	MaximumTime      int64   `json:"maxTime" bson:"maxTime"`
	ExfiltrationTime int64   `json:"exfilTime" bson:"exfilTime"`
	Requirement      string  `json:"requirement,omitempty" bson:"requirement,omitempty"`
}

// Boss describes a boss of a location
type Boss struct {
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Chance      float64 `json:"chance" bson:"chance"`
	Followers   int64   `json:"followers" bson:"followers"`
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
