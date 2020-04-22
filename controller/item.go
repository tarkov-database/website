package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/model/item"
	"github.com/tarkov-database/website/view"

	"github.com/julienschmidt/httprouter"
)

func ItemGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	kind, err := item.CategoryToKind(ps.ByName("category"))
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	timeAPI := time.Now()

	entity, err := item.GetItem(ps.ByName("id"), kind)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	addTimingHeader(timingMetrics{"api": time.Since(timeAPI)}, w)

	var tmpl string
	if strings.HasPrefix(kind.String(), "modification") {
		tmpl = "item_modification"
	} else {
		tmpl = fmt.Sprintf("item_%v", kind)
	}

	w.Header().Set("Trailer", "Server-Timing")

	timeRender := time.Now()

	view.RenderHTML(tmpl, p.Entity(entity), w)

	addTimingHeader(timingMetrics{"render": time.Since(timeRender)}, w)
}

func ItemsGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := ps.ByName("category")

	kind, err := item.CategoryToKind(c)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	params := make(map[string]string)
	switch kind {
	case item.KindAmmunition:
		params["type"], params["caliber"] = r.URL.Query().Get("type"), r.URL.Query().Get("caliber")
	case item.KindMagazine:
		params["caliber"] = r.URL.Query().Get("caliber")
	case item.KindFirearm:
		params["type"], params["caliber"], params["class"] = r.URL.Query().Get("type"), r.URL.Query().Get("caliber"), r.URL.Query().Get("class")
	case item.KindArmor:
		params["type"], params["class"] = r.URL.Query().Get("type"), r.URL.Query().Get("class")
	case item.KindTacticalrig:
		params["class"] = r.URL.Query().Get("class")
	case item.KindMedical, item.KindFood, item.KindGrenade, item.KindClothing, item.KindModificationMuzzle, item.KindModificationDevice, item.KindModificationSight, item.KindModificationSightSpecial, item.KindModificationGoggles, item.KindModificationGogglesSpecial:
		params["type"] = r.URL.Query().Get("type")
	}

	opts := &api.Options{
		Sort:   r.URL.Query().Get("sort"),
		Filter: params,
	}
	opts.Limit, opts.Offset = getLimitOffset(getPage(r))

	timeAPI := time.Now()

	result, err := item.GetItems(kind, opts)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	addTimingHeader(timingMetrics{"api": time.Since(timeAPI)}, w)

	cat, err := item.KindToCategory(kind)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	kw, err := item.CategoryToDisplayName(cat)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	data := p.Result(result, kw)

	var tmpl string
	switch r.URL.Query().Get("view") {
	case "table":
		tmpl = fmt.Sprintf("table_%v", kind)
	default:
		tmpl = "list"
	}

	w.Header().Set("Trailer", "Server-Timing")

	timeRender := time.Now()

	view.RenderHTML(tmpl, data, w)

	addTimingHeader(timingMetrics{"render": time.Since(timeRender)}, w)
}
