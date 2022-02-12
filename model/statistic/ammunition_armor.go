package statistic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model/item"
)

// AmmoArmorStatistics describes the entity of a ammo against armor statistics
type AmmoArmorStatistics struct {
	ID                        objectID   `json:"_id"`
	Ammo                      objectID   `json:"ammo"`
	Armor                     ItemRef    `json:"armor"`
	Distance                  uint64     `json:"distance"`
	PenetrationChance         [4]float64 `json:"penetrationChance"`
	AverageShotsToDestruction Statistics `json:"avgShotsToDestruct"`
	AverageShotsTo50Damage    Statistics `json:"avgShotsTo50Damage"`
	Modified                  timestamp  `json:"_modified"`
}

// ItemRef refers to an item entity
type ItemRef struct {
	ID   objectID  `json:"id"`
	Kind item.Kind `json:"kind"`
}

// Statistics describes the statistical values
type Statistics struct {
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	StdDev float64 `json:"stdDev"`
}

type AmmoArmorStatisticsResult struct {
	Count int64                 `json:"total"`
	Items []AmmoArmorStatistics `json:"items"`
}

func GetAmmoArmorStatistic(id objectID) (*AmmoArmorStatistics, error) {
	stats := &AmmoArmorStatistics{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/statistic/ammunition/armor/%s", id), &api.Options{}, stats); err != nil {
		return stats, err
	}

	return stats, nil
}

func GetAmmoArmorStatistics(ammoIDs, armorIDs []objectID, gte, lte uint64, limit int) (*AmmoArmorStatisticsResult, error) {
	opts := &api.Options{
		Limit: limit,
		Filter: map[string]string{
			"ammo":  strings.Join(ammoIDs, ","),
			"armor": strings.Join(armorIDs, ","),
			"range": fmt.Sprintf("%v,%v", gte, lte),
		},
	}

	result := &AmmoArmorStatisticsResult{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, "/statistic/ammunition/armor", opts, result); err != nil {
		return result, err
	}

	return result, nil
}

func GetAmmoArmorStatisticsByChunks(ammoIDs, armorIDs []objectID, gte, lte uint64, limit int) (*AmmoArmorStatisticsResult, error) {
	opts := &api.Options{
		Limit: limit,
		Filter: map[string]string{
			"ammo":  strings.Join(ammoIDs, ","),
			"armor": strings.Join(armorIDs, ","),
			"range": fmt.Sprintf("%v,%v", gte, lte),
		},
	}

	result := &AmmoArmorStatisticsResult{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, "/statistic/ammunition/armor", opts, result); err != nil {
		return result, err
	}

	return result, nil
}
