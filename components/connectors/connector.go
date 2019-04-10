package connectors

import (
	"database/sql"
	"net/http"

	a "github.com/squanchersquanch/contacts/components/actions"
)

const (
	baseURI = ""
)

// Connector manages http communications for the app
type Connector interface {
	NotFound(w http.ResponseWriter, r *http.Request)
	CreateContact(w http.ResponseWriter, r *http.Request)
	GetContacts(w http.ResponseWriter, r *http.Request)
	UpdateContact(w http.ResponseWriter, r *http.Request)
	DeleteContact(w http.ResponseWriter, r *http.Request)
	ImportContacts(w http.ResponseWriter, r *http.Request)
	ExportContacts(w http.ResponseWriter, r *http.Request)
}

// connector is an implementation of the Connector interface
type connector struct {
	actions a.Actions
}

// NewConnector creates a new instance of Connector
func NewConnector(
	db *sql.DB,
) Connector {
	actions := a.NewActions(db)
	return &connector{
		actions: actions,
	}
}

// NotFound calls the NotFound action and returns a StatusNotFound
func (c *connector) NotFound(w http.ResponseWriter, r *http.Request) {
	c.actions.NotFound(w)
}

// CreateContact creates a new contact
func (c *connector) CreateContact(w http.ResponseWriter, r *http.Request) {
	c.actions.CreateRow(w, r)
}

// GetContacts retrieves all contacts or a specifc contact by id
func (c *connector) GetContacts(w http.ResponseWriter, r *http.Request) {
	query := c.getURLQuery(r, "id")
	if query != "" {
		c.actions.ReadRows(w, query)
		return
	}

	c.actions.ReadRows(w)
}

// UpdateContact updates an existing contact
func (c *connector) UpdateContact(w http.ResponseWriter, r *http.Request) {
	c.actions.UpdateRow(w, r)
}

// DeleteContact deletes an existing contacts
func (c *connector) DeleteContact(w http.ResponseWriter, r *http.Request) {
	query := c.getURLQuery(r, "id")
	if query != "" {
		c.actions.DeleteRow(w, query)
	}
}

// ExportContacts exports existing contacts via csv file
func (c *connector) ExportContacts(w http.ResponseWriter, r *http.Request) {
	c.actions.GenerateContactsCSV(w, r)
}

// ImportContacts updates an existing contact via csv file
func (c *connector) ImportContacts(w http.ResponseWriter, r *http.Request) {
	c.actions.ImportContactsCSV(w, r)
}

// getURLQuery returns values of URL query from given key
func (c *connector) getURLQuery(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
