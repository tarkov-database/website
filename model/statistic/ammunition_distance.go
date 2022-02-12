package statistic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tarkov-database/website/core/api"
)

// AmmoDistanceStatistics describes the entity of a ammo distance statistics
type AmmoDistanceStatistics struct {
	ID               objectID  `json:"_id"`
	Reference        objectID  `json:"ammo"`
	Distance         uint64    `json:"distance"`
	Velocity         float64   `json:"velocity"`
	Damage           float64   `json:"damage"`
	PenetrationPower float64   `json:"penetrationPower"`
	TimeOfFlight     float64   `json:"timeOfFlight"`
	Drop             float64   `json:"drop"`
	Modified         timestamp `json:"_modified"`
}

type AmmoDistanceStatisticsResult struct {
	Count int64                    `json:"total"`
	Items []AmmoDistanceStatistics `json:"items"`
}

func GetAmmoDistanceStatistic(id objectID) (*AmmoDistanceStatistics, error) {
	stats := &AmmoDistanceStatistics{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/statistic/ammunition/distance/%s", id), &api.Options{}, stats); err != nil {
		return stats, err
	}

	return stats, nil
}

func GetAmmoDistanceStatistics(ids []objectID, gte, lte uint64, limit int) (*AmmoDistanceStatisticsResult, error) {
	opts := &api.Options{
		Limit: limit,
		Filter: map[string]string{
			"ammo":  strings.Join(ids, ","),
			"range": fmt.Sprintf("%v,%v", gte, lte),
		},
	}

	result := &AmmoDistanceStatisticsResult{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, "/statistic/ammunition/distance", opts, result); err != nil {
		return result, err
	}

	return result, nil
}
