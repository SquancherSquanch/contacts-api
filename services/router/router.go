package router

import (
	"database/sql"
	"net/http"

	"github.com/squanchersquanch/contacts/components/connectors"

	"github.com/gorilla/mux"
	"github.com/squanchersquanch/contacts/services/config"
	"github.com/squanchersquanch/contacts/services/logger"
	r "github.com/squanchersquanch/contacts/services/routes"
)

// NewRouter ...
func NewRouter(db *sql.DB, config *config.Config) *mux.Router {
	c := connectors.NewConnector(db, config)
	router := mux.NewRouter().StrictSlash(true)
	routes := r.NewRoutes(c)
	for _, route := range routes.RouteList() {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
