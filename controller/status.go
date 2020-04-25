package controller

import (
	"errors"
	"net/http"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/model/item"
	"github.com/tarkov-database/website/view"

	"github.com/google/logger"
)

func statusNotFound(w http.ResponseWriter, r *http.Request) {
	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil {
		getErrorStatus(err, w, r)
		return
	}

	status := http.StatusNotFound

	switch r.Header.Get("Content-Type") {
	case "application/json", "application/geo+json":
		view.RenderJSON(model.NewResponse("Entity not found", status), status, w)
	default:
		w.WriteHeader(status)
		view.RenderHTML("status_404", p, w)
	}
}

func StatusNotFoundHandler() http.Handler {
	return http.HandlerFunc(statusNotFound)
}

func statusInternalServerError(w http.ResponseWriter, r *http.Request) {
	status := http.StatusInternalServerError

	switch r.Header.Get("Content-Type") {
	case "application/json", "application/geo+json":
		view.RenderJSON(model.NewResponse("Internal Server Error", status), status, w)
	default:
		w.WriteHeader(status)
		view.RenderHTML("status_500", model.CreatePage(r.URL), w)
	}
}

func statusServiceUnavailable(w http.ResponseWriter, r *http.Request) {
	status := http.StatusServiceUnavailable

	switch r.Header.Get("Content-Type") {
	case "application/json", "application/geo+json":
		view.RenderJSON(model.NewResponse("API is not available", status), status, w)
	default:
		w.WriteHeader(status)
		view.RenderHTML("status_503", model.CreatePage(r.URL), w)
	}
}

func statusServiceBadRequest(w http.ResponseWriter, r *http.Request) {
	status := http.StatusBadRequest

	switch r.Header.Get("Content-Type") {
	case "application/json", "application/geo+json":
		view.RenderJSON(model.NewResponse("Bad Request", status), status, w)
	default:
		w.WriteHeader(status)
		view.RenderHTML("status_400", model.CreatePage(r.URL), w)
	}
}

func statusUnsupportedMediaType(w http.ResponseWriter, _ *http.Request) {
	status := http.StatusUnsupportedMediaType
	view.RenderJSON(model.NewResponse("Content type is not supported", status), status, w)
}

func getErrorStatus(err error, w http.ResponseWriter, r *http.Request) {
	var status int
	var apiResponse *api.Response
	if errors.As(err, &apiResponse) {
		status = apiResponse.StatusCode
	}

	switch {
	case status == 404, errors.Is(err, item.ErrInvalidCategory):
		statusNotFound(w, r)
	case errors.Is(err, api.ErrUnreachable):
		logger.Error(err)
		statusServiceUnavailable(w, r)
	default:
		logger.Error(err)
		statusInternalServerError(w, r)
	}
}
