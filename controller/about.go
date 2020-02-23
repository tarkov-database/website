package controller

import (
	"net/http"

	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/view"

	"github.com/julienschmidt/httprouter"
)

func AboutGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view.RenderHTML("about", model.CreatePage(r.URL), w)
}
