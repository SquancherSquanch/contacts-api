package routes

import (
	"net/http"

	"github.com/squanchersquanch/contacts/components/connectors"
)

// Routes manages routes for the http calls and API endpoints
type Routes interface {
	RouteList() []Route
}

// Route is a single route or endpoint
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// routes is an implementation of the Routes interface
type routes struct {
	connector connectors.Connector
}

// NewRoutes creates a new routes interface with connectors
func NewRoutes(
	connector connectors.Connector,
) Routes {
	return &routes{
		connector: connector,
	}
}

// RouteList returns an array of Routes
func (r *routes) RouteList() []Route {
	return []Route{
		Route{
			"CreateContact",
			"POST",
			"/api/entry",
			r.connector.CreateContact,
		},
		Route{
			"GetContact",
			"GET",
			"/api/entry",
			r.connector.GetContacts,
		},
		Route{
			"UpdateContact",
			"PUT",
			"/api/entry",
			r.connector.UpdateContact,
		},
		Route{
			"DeleteContact",
			"DELETE",
			"/api/entry",
			r.connector.DeleteContact,
		},
		Route{
			"ExportContacts",
			"GET",
			"/api/entry/export",
			r.connector.ExportContacts,
		},
		Route{
			"ImportContacts",
			"POST",
			"/api/entry/import",
			r.connector.ImportContacts,
		},
		Route{
			"NotFound",
			"",
			"",
			r.connector.NotFound,
		},
	}
}
