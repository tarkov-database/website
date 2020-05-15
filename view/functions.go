package view

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
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

func hasPrefix(v interface{}, p string) bool {
	var s string

	switch v := v.(type) {
	case string:
		s = v
	case item.Kind:
		s = v.String()
	}

	return strings.HasPrefix(s, p)
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

	a := u.Query().Get(key)
	b := fmt.Sprintf("%v", val)
	if a != "" && (b == a || b == "*") {
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

func resolveItemList(list item.ItemList, sort string) resolvedItemList {
	if sort == "" {
		sort = "name"
	}

	return item.GetItemList(list, sort)
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

func resolveSlots(slots item.Slots, sort string) resolvedSlots {
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	length := len(slots)

	wg.Add(length)

	rs := make(resolvedSlots, length)
	for k, v := range slots {
		go func(n string, s item.Slot) {
			defer wg.Done()

			p := strings.SplitN(n, "_", 2)

			switch t := p[0]; t {
			case "magazine":
				return
			case "barrel":
				n = "Barrel"
			case "muzzle":
				n = "Muzzle"
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
			case "equipment":
				n = "Equipment"
			case "tactical":
				n = "Tactical"
			case "mount":
				n = "Mount"
			case "flashlight":
				n = "Flashlight"
			case "foregrip":
				n = "Foregrip"
			case "nvg":
				n = "Night Vision"
			case "scope":
				n = "Scope"
			case "sightFront":
				n = "Front Sight"
			case "sightRear":
				n = "Rear Sight"
			case "trigger":
				n = "Trigger"
			case "hammer":
				n = "Hammer"
			case "catch":
				n = "Catch"
			default:
				logger.Warningf("Unknown slot name \"%s\"", t)
			}

			if len(p) > 1 {
				if i, err := strconv.ParseInt(p[1], 10, 64); err == nil {
					if i != 0 {
						n += fmt.Sprintf(" %v", i+1)
					}
				} else {
					logger.Error(err)
				}
			}

			res := item.GetItemList(s.Filter, sort)
			mutex.Lock()
			rs[n] = res
			mutex.Unlock()
		}(k, v)
	}

	wg.Wait()

	return rs
}
