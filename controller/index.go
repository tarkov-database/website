package controller

import (
	"net/http"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/view"

	"github.com/julienschmidt/httprouter"
)

func IndexGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil && err != api.ErrUnreachable {
		getErrorStatus(err, w, r)
		return
	}

	view.Render("index", p.GetIndex(), w)
}
