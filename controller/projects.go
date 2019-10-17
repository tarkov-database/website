package controller

import (
	"net/http"

	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/view"

	"github.com/julienschmidt/httprouter"
)

func ProjectsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view.Render("projects", model.CreatePage(r.URL), w)
}
