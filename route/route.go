package route

import (
	"net/http"

	cntrl "github.com/tarkov-database/website/controller"

	"github.com/julienschmidt/httprouter"
)

func Load() *httprouter.Router {
	return routes()
}

func routes() *httprouter.Router {
	r := httprouter.New()

	// Index
	r.GET("/", cntrl.IndexGET)

	// Item
	r.GET("/item/:category", cntrl.ItemsGET)
	r.GET("/item/:category/:id", cntrl.ItemGET)

	// Location
	r.GET("/location", cntrl.LocationsGET)
	r.GET("/location/:id", cntrl.LocationGET)

	// Search
	r.GET("/search", cntrl.SearchGET)
	r.GET("/search/ws", cntrl.SearchWS)

	// About
	r.GET("/about", cntrl.AboutGET)

	// Projects
	r.GET("/projects", cntrl.ProjectsGET)

	// Static
	r.ServeFiles("/resources/*filepath", http.Dir("static/dist/resources"))

	// Health
	r.GET("/health", cntrl.HealthGET)

	// Status
	r.NotFound = cntrl.StatusNotFoundHandler()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	return r
}
