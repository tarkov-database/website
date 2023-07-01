package view

import (
	"encoding/base64"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model/item"
	"github.com/tarkov-database/website/model/statistic"
	"github.com/tarkov-database/website/version"

	"github.com/goccy/go-json"
	"github.com/google/logger"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func addFloat(a, b interface{}) float64 {
	var af, bf float64

	switch v := a.(type) {
	case float64:
		af = v
	case int64:
		af = float64(v)
	case int:
		af = float64(v)
	}

	switch v := b.(type) {
	case float64:
		bf = v
	case int64:
		bf = float64(v)
	case int:
		bf = float64(v)
	}

	return af + bf
}

func subtractFloat(a, b interface{}) float64 {
	var af, bf float64

	switch v := a.(type) {
	case float64:
		af = v
	case int64:
		af = float64(v)
	case int:
		af = float64(v)
	}

	switch v := b.(type) {
	case float64:
		bf = v
	case int64:
		bf = float64(v)
	case int:
		bf = float64(v)
	}

	return af - bf
}

func multiplyFloat(a, b interface{}) float64 {
	var af, bf float64

	switch v := a.(type) {
	case float64:
		af = v
	case int64:
		af = float64(v)
	case int:
		af = float64(v)
	}

	switch v := b.(type) {
	case float64:
		bf = v
	case int64:
		bf = float64(v)
	case int:
		bf = float64(v)
	}

	return af * bf
}

func divideFloat(a, b interface{}) float64 {
	var af, bf float64

	switch v := a.(type) {
	case float64:
		af = v
	case int64:
		af = float64(v)
	case int:
		af = float64(v)
	}

	switch v := b.(type) {
	case float64:
		bf = v
	case int64:
		bf = float64(v)
	case int:
		bf = float64(v)
	}

	return af / bf
}

func decimalToPercent(d float64) float64 {
	return math.Round((d*100)*100) / 100
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
	if sum, err := version.SumOf(strings.TrimPrefix(p, "/")); err == nil {
		p += "?v=" + sum[:8]
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

	return u.Path + "?" + q.Encode()
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

func localeString(v interface{}) string {
	return message.NewPrinter(language.English).Sprint(v)
}

var regexUpperCase = regexp.MustCompile(`([A-Z])`)

func camelToTitle(s string) string {
	return strings.Title(regexUpperCase.ReplaceAllString(s, " $1"))
}

func toBase64(v interface{}) string {
	j, err := json.Marshal(v)
	if err != nil {
		logger.Error(err)
	}

	return base64.StdEncoding.EncodeToString(j)
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

func getAmmunitionRangeData(ammo *item.Ammunition) []statistic.AmmoDistanceStatistics {
	max := 1000

	switch {
	case ammo.Type == "buckshot":
		max = 100
	case ammo.Subsonic:
		max = 500
	}

	ids := []string{ammo.ID}

	result, err := statistic.GetAmmoDistanceStatistics(ids, 0, uint64(max), 100)
	if err != nil {
		logger.Error(err)
		return nil
	}

	if result.Count == 0 {
		return nil
	}

	switch {
	case result.Count < 30:
	case max > 300 && result.Items[30].Drop < -16:
		max = 300
	case max > 500 && result.Items[50].Drop < -24:
		max = 500
	}

	if max/10 > len(result.Items) {
		max = len(result.Items) * 10
	}

	steps := int(max / 10)
	items := make([]statistic.AmmoDistanceStatistics, 0, 10)
	for _, v := range result.Items[:steps] {
		if int(v.Distance)%steps == 0 {
			items = append(items, v)
		}
	}

	return items
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
