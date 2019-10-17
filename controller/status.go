package controller

import (
	"net/http"
	"strings"

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

	w.WriteHeader(http.StatusNotFound)
	view.Render("status_404", p, w)
}

func StatusNotFoundHandler() http.Handler {
	return http.HandlerFunc(statusNotFound)
}

func statusInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	view.Render("status_500", model.CreatePage(r.URL), w)
}

func statusServiceUnavailable(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	view.Render("status_503", model.CreatePage(r.URL), w)
}

func getErrorStatus(err error, w http.ResponseWriter, r *http.Request) {
	code := strings.Split(err.Error(), ":")[0]
	switch {
	case code == "404", err == item.ErrInvalidCategory:
		statusNotFound(w, r)
	case err == api.ErrUnreachable:
		logger.Error(err)
		statusServiceUnavailable(w, r)
	default:
		logger.Error(err)
		statusInternalServerError(w, r)
	}
}
