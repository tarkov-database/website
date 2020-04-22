package controller

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/model/location"
	"github.com/tarkov-database/website/model/location/feature"
	"github.com/tarkov-database/website/view"

	"github.com/julienschmidt/httprouter"
)

func LocationGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	timeAPI := time.Now()

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

	addTimingHeader(timingMetrics{"api": time.Since(timeAPI)}, w)

	w.Header().Set("Trailer", "Server-Timing")

	timeRender := time.Now()

	view.RenderHTML("location", p.Entity(entity), w)

	addTimingHeader(timingMetrics{"render": time.Since(timeRender)}, w)
}

func LocationsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := make(map[string]string)
	params["available"] = r.URL.Query().Get("available")

	opts := &api.Options{
		Sort:   r.URL.Query().Get("sort"),
		Filter: params,
	}
	opts.Limit, opts.Offset = getLimitOffset(getPage(r))

	timeAPI := time.Now()

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

	addTimingHeader(timingMetrics{"api": time.Since(timeAPI)}, w)

	data := p.Result(result, "location")

	w.Header().Set("Trailer", "Server-Timing")

	timeRender := time.Now()

	view.RenderHTML("list", data, w)

	addTimingHeader(timingMetrics{"render": time.Since(timeRender)}, w)
}

func LocationMapGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	timeAPI := time.Now()

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

	addTimingHeader(timingMetrics{"api": time.Since(timeAPI)}, w)

	w.Header().Set("Trailer", "Server-Timing")

	timeRender := time.Now()
	view.RenderHTML("location_map", p.Entity(entity), w)

	addTimingHeader(timingMetrics{"render": time.Since(timeRender)}, w)
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

	var err error

	lID := ps.ByName("id")

	limit := 100
	if v := r.URL.Query().Get("limit"); v != "" {
		if limit, err = strconv.Atoi(v); err != nil {
			statusServiceBadRequest(w, r)
			return
		}
	}

	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		if offset, err = strconv.Atoi(v); err != nil {
			statusServiceBadRequest(w, r)
			return
		}
	}

	opts := &api.Options{
		Limit:  limit,
		Offset: offset,
		Filter: make(map[string]string),
	}

	var result *feature.FeatureResult
Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "group":
			q, err := url.QueryUnescape(v[0])
			if err != nil {
				statusServiceBadRequest(w, r)
				return
			}

			result, err = feature.GetFeaturesByGroup(q, lID, opts)
			if err != nil {
				getErrorStatus(err, w, r)
				return
			}

			break Loop
		case "text":
			q, err := url.QueryUnescape(v[0])
			if err != nil {
				statusServiceBadRequest(w, r)
				return
			}

			if err := validateKeyword(q); err != nil {
				statusServiceBadRequest(w, r)
				return
			}

			q = cleanupString(q)

			result, err = feature.GetFeaturesByText(q, lID, 50)
			if err != nil {
				getErrorStatus(err, w, r)
				return
			}

			break Loop
		}
	}

	if result == nil {
		result, err = feature.GetFeatures(lID, opts)
		if err != nil {
			getErrorStatus(err, w, r)
			return
		}
	}

	switch v := r.Header.Get("Content-Type"); v {
	case "application/json":
		w.Header().Set("Content-Type", v)
		view.RenderJSON(result, http.StatusOK, w)
	case "application/geo+json":
		w.Header().Set("Content-Type", v)
		view.RenderJSON(result.FeatureCollection(), http.StatusOK, w)
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
