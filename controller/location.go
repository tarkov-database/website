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

	view.Render("location", p.Location(entity), w)
}

func LocationsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := make(map[string]string)
	params["available"] = r.URL.Query().Get("available")

	opts := &api.Options{
		Sort:   r.URL.Query().Get("sort"),
		Filter: params,
	}
	opts.Limit, opts.Offset = getLimitOffset(getPage(r))

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

	data := p.LocationResult(result, "Locations", false)

	view.Render("list", data, w)
}
