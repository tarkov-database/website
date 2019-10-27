package view

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model/item"
	"github.com/tarkov-database/website/version"

	"github.com/google/logger"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func addFloat(a, b float64) float64 {
	return a + b
}

func subtractFloat(a, b float64) float64 {
	return a - b
}

func multiplyFloat(a, b float64) float64 {
	return a * b
}

func appendStaticHash(p string) string {
	if sum, ok := version.StaticSums[strings.TrimPrefix(p, "/")]; ok {
		p += fmt.Sprintf("?v=%s", sum[:8])
	}

	return p
}

func formatTime(format string, t time.Time) string {
	switch format {
	case "RFC3339":
		return t.Format(time.RFC3339)
	default:
		return t.Format(format)
	}
}

func setQuery(path, key string, val interface{}) string {
	u, err := url.Parse(path)
	if err != nil {
		logger.Error(err)
		return path
	}

	q := u.Query()
	q.Set(key, fmt.Sprintf("%v", val))

	return fmt.Sprintf("%s?%s", u.Path, q.Encode())
}

func hasQuery(path, key string, val interface{}) bool {
	u, err := url.Parse(path)
	if err != nil {
		logger.Error(err)
		return false
	}

	q := u.Query()
	if q.Get(key) == fmt.Sprintf("%v", val) {
		return true
	}

	return false
}

func decimalToPercent(d float64) float64 {
	return math.Round((d*100)*100) / 100
}

func localeString(v interface{}) string {
	return message.NewPrinter(language.English).Sprint(v)
}

func getItem(id string, kind item.Kind) item.Entity {
	result, err := item.GetItem(id, kind)
	if err != nil {
		logger.Error(err)
	}

	return result
}

type resolvedItemList map[item.Kind][]item.Entity

func resolveItemList(list item.ItemList) resolvedItemList {
	return item.GetItemList(list)
}

func getAmmunitionByCaliber(caliber, sort string, limit int) item.EntityResult {
	opts := &api.Options{
		Limit: limit,
		Sort:  sort,
		Filter: map[string]string{
			"caliber": caliber,
		},
	}

	result, err := item.GetItems(item.KindAmmunition, opts)
	if err != nil {
		logger.Error(err)
	}

	return result
}

type resolvedSlots map[string]resolvedItemList

func resolveSlots(slots item.Slots) resolvedSlots {
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(slots))

	rs := make(resolvedSlots)
	for k, v := range slots {
		go func(n string, s item.Slot) {
			defer wg.Done()

			switch n {
			case "magazine":
				return
			case "barrel":
				n = "Barrel"
			case "muzzle", "muzzle_00":
				n = "Muzzle"
			case "muzzle_01":
				n = "Muzzle 2"
			case "receiver":
				n = "Receiver"
			case "gasBlock":
				n = "Gas Block"
			case "pistolGrip":
				n = "Pistol Grip"
			case "stock":
				n = "Stock"
			case "charge":
				n = "Charge"
			case "handguard":
				n = "Handguard"
			case "bipod":
				n = "Bipod"
			case "launcher":
				n = "Launcher"
			case "equipment_00":
				n = "Equipment 1"
			case "equipment_01":
				n = "Equipment 2"
			case "equipment_02":
				n = "Equipment 3"
			case "tactical_00":
				n = "Tactical 1"
			case "tactical_01":
				n = "Tactical 2"
			case "tactical_02":
				n = "Tactical 3"
			case "tactical_03":
				n = "Tactical 4"
			case "mount_00":
				n = "Mount 1"
			case "mount_01":
				n = "Mount 2"
			case "mount_02":
				n = "Mount 3"
			case "mount_03":
				n = "Mount 4"
			case "mount_04":
				n = "Mount 5"
			case "mount_05":
				n = "Mount 6"
			case "flashlight":
				n = "Flashlight"
			case "foregrip":
				n = "Foregrip"
			case "nvg":
				n = "Night Vision"
			case "scope_00":
				n = "Scope 1"
			case "scope_01":
				n = "Scope 2"
			case "scope_02":
				n = "Scope 3"
			case "scope_03":
				n = "Scope 4"
			case "sightFront":
				n = "Front Sight"
			case "sightRear":
				n = "Rear Sight"
			default:
				logger.Warningf("Unknown slot name \"%s\"", n)
			}

			res := item.GetItemList(s.Filter)
			mutex.Lock()
			rs[n] = res
			mutex.Unlock()
		}(k, v)
	}

	wg.Wait()

	return rs
}
