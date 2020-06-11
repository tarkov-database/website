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

	if err := validateQueryValues(r.URL.Query()); err != nil {
		statusBadRequest(w, r)
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
		params["type"], params["armor.class"], params["armor.material.name"] = r.URL.Query().Get("type"), r.URL.Query().Get("class"), r.URL.Query().Get("material")
	case item.KindTacticalrig:
		params["armor.class"], params["armor.material.name"] = r.URL.Query().Get("class"), r.URL.Query().Get("material")
	case item.KindMedical, item.KindFood, item.KindGrenade, item.KindClothing, item.KindModificationMuzzle, item.KindModificationDevice, item.KindModificationSight, item.KindModificationSightSpecial, item.KindModificationGoggles, item.KindModificationGogglesSpecial:
		params["type"] = r.URL.Query().Get("type")
	}

	opts := &api.Options{
		Sort:   r.URL.Query().Get("sort"),
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
		tmpl = fmt.Sprintf("table_%v", kind)
	default:
		tmpl = "list"
	}

	view.RenderHTML(tmpl, data, w)
}
