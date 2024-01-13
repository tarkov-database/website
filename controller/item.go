package controller

import (
	"fmt"
	"net/http"
	"strings"

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

	var tmpl string
	if strings.HasPrefix(kind.String(), "modification") {
		tmpl = "item_modification"
	} else {
		tmpl = fmt.Sprintf("item_%v", kind)
	}

	view.RenderHTML(tmpl, p.Entity(entity), w)
}

func ItemsGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cat := ps.ByName("category")

	kind, err := item.CategoryToKind(cat)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	queryValues := r.URL.Query()

	if err := validateQueryValues(queryValues); err != nil {
		statusBadRequest(w, r)
		return
	}

	params := make(map[string]string)
	switch kind {
	case item.KindAmmunition:
		params["caliber"] = queryValues.Get("caliber")
		params["type"] = queryValues.Get("type")
	case item.KindMagazine:
		params["caliber"] = queryValues.Get("caliber")
	case item.KindFirearm:
		params["manufacturer"] = queryValues.Get("manufacturer")
		params["caliber"] = queryValues.Get("caliber")
		params["type"], params["class"] = queryValues.Get("type"), queryValues.Get("class")
	case item.KindArmor:
		params["armor.material.name"] = queryValues.Get("material")
		params["type"], params["armor.class"] = queryValues.Get("type"), queryValues.Get("class")
	case item.KindTacticalrig:
		params["armor.material.name"] = queryValues.Get("material")
		params["isArmored"], params["armor.class"] = queryValues.Get("armored"), queryValues.Get("class")
		params["isPlateCarrier"] = queryValues.Get("plateCarrier")
	case item.KindMedical, item.KindFood, item.KindGrenade, item.KindClothing, item.KindModificationMuzzle, item.KindModificationDevice, item.KindModificationSight, item.KindModificationSightSpecial, item.KindModificationGoggles, item.KindModificationGogglesSpecial:
		params["type"] = queryValues.Get("type")
	}

	if err := unescapeParams(params); err != nil {
		statusBadRequest(w, r)
		return
	}

	opts := &api.Options{
		Sort:   queryValues.Get("sort"),
		Filter: params,
	}
	opts.Limit, opts.Offset = getLimitOffset(getPage(r.URL))

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

	kw, err := item.CategoryToDisplayName(cat)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	data, err := p.Result(result, kw)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	var tmpl string
	switch r.URL.Query().Get("view") {
	case "table":
		if strings.HasPrefix(kind.String(), "modification") {
			tmpl = "table_modification"
		} else {
			tmpl = fmt.Sprintf("table_%v", kind)
		}
	default:
		tmpl = fmt.Sprintf("grid_%v", kind)
	}

	view.RenderHTML(tmpl, data, w)
}
