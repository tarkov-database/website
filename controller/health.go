package controller

import (
	"net/http"

	"github.com/tarkov-database/website/model/health"
	"github.com/tarkov-database/website/view"

	"github.com/julienschmidt/httprouter"
)

func HealthGET(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	h := health.GetHealth()

	if !h.OK {
		view.RenderJSON(h, http.StatusInternalServerError, w)
	} else {
		view.RenderJSON(h, http.StatusOK, w)
	}
}
