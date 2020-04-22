package controller

import (
	"net/http"
	"time"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/model"
	"github.com/tarkov-database/website/view"

	"github.com/julienschmidt/httprouter"
)

func IndexGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	timeAPI := time.Now()

	p, err := model.CreatePageWithAPI(r.URL)
	if err != nil && err != api.ErrUnreachable {
		getErrorStatus(err, w, r)
		return
	}

	addTimingHeader(timingMetrics{"api": time.Since(timeAPI)}, w)

	w.Header().Set("Trailer", "Server-Timing")

	timeRender := time.Now()

	view.RenderHTML("index", p.GetIndex(), w)

	addTimingHeader(timingMetrics{"render": time.Since(timeRender)}, w)
}
