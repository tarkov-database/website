package controller

import (
	"net/http"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/model/location"
	"github.com/tarkov-database/website/view"

	"github.com/julienschmidt/httprouter"
)

func LocationGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	entity, err := location.GetLocation(ps.ByName("id"))
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	view.RenderHTML("location", p.Entity(entity), w)
}

func LocationsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := validateQueryValues(r.URL.Query()); err != nil {
		statusBadRequest(w, r)
		return
	}

	params := make(map[string]string)
	params["available"] = r.URL.Query().Get("available")

	opts := &api.Options{
		Sort:   r.URL.Query().Get("sort"),
		Filter: params,
	}
	opts.Limit, opts.Offset = getLimitOffset(getPage(r.URL))

	result, err := location.GetLocations(opts)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	data, err := p.Result(result, "location")
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	var tmpl string
	switch r.URL.Query().Get("view") {
	case "table":
		tmpl = "table_location"
	default:
		tmpl = "grid_location"
	}

	view.RenderHTML(tmpl, data, w)
}
