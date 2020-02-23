package controller

import (
	"net/http"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/model/location"
	"github.com/tarkov-database/website/model/location/feature"
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

	view.RenderHTML("location", p.Location(entity), w)
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

	data := p.LocationResult(result, "location", false)

	view.RenderHTML("list", data, w)
}

func LocationMapGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	view.RenderHTML("location_map", p.Location(entity), w)
}

func LocationFeatureGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		statusUnsupportedMediaType(w, r)
		return
	}

	entity, err := feature.GetFeature(ps.ByName("id"), ps.ByName("fid"))
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	view.RenderJSON(entity, http.StatusOK, w)
}

func LocationFeaturesGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		statusUnsupportedMediaType(w, r)
		return
	}

	params := make(map[string]string)
	params["group"] = r.URL.Query().Get("group")

	opts := &api.Options{
		Sort:   r.URL.Query().Get("sort"),
		Filter: params,
	}

	entity, err := feature.GetFeatures(ps.ByName("id"), opts)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	switch v := r.Header.Get("Content-Type"); v {
	case "application/json":
		w.Header().Set("Content-Type", v)
		view.RenderJSON(entity, http.StatusOK, w)
	case "application/geo+json":
		w.Header().Set("Content-Type", v)
		view.RenderJSON(entity.FeatureCollection(), http.StatusOK, w)
	}
}

func LocationFeatureGroupGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		statusUnsupportedMediaType(w, r)
		return
	}

	entity, err := feature.GetGroup(ps.ByName("id"), ps.ByName("gid"))
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	view.RenderJSON(entity, http.StatusOK, w)
}

func LocationFeatureGroupsGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		statusUnsupportedMediaType(w, r)
		return
	}

	params := make(map[string]string)

	opts := &api.Options{
		Sort:   r.URL.Query().Get("sort"),
		Filter: params,
	}

	entity, err := feature.GetGroups(ps.ByName("id"), opts)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	view.RenderJSON(entity, http.StatusOK, w)
}
