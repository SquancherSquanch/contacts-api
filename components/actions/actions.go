package actions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/squanchersquanch/contacts/models"
	"github.com/squanchersquanch/contacts/services/config"
)

// sql constants
const (
	selectFrom      = "SELECT * FROM %s;"
	selectFromWhere = "SELECT * FROM %s WHERE id='%s';"
	deleteFrom      = "DELETE FROM %s WHERE id=$1;"
	update          = `UPDATE %s SET firstName='%s', lastName='%s', email='%s', phone='%s' WHERE id='%s';`

	insertInto = `INSERT INTO %s (firstName, lastName, email, phone)
					VALUES ($1, $2, $3, $4)
					RETURNING id;`

	commandQuery    = "Query"
	commandQueryRow = "QueryRow"
)

// messaging constants
const (
	numberOfEntries = "%d entries found"
	duplicateEntry  = `pq: duplicate key value violates unique constraint "entries_email_key"`
	invalidID       = "invalid id provided"
	invalidFileType = "invalid files type"
	notFound        = "not found"
)

// header constants
const (
	contentTypeHeader        = "Content-Type"
	contentDispositionHeader = "Content-Disposition"

	csvContentType        = "text/csv"
	csvContentDisposition = "attachment; filename=contacts.csv"
)

// Actions manages http requests from the connector
// in return providing a response along with interacting with the database for the app
type Actions interface {
	NotFound(w http.ResponseWriter)
	CreateRow(w http.ResponseWriter, r *http.Request)
	ReadRows(w http.ResponseWriter, id ...string)
	UpdateRow(w http.ResponseWriter, r *http.Request)
	DeleteRow(w http.ResponseWriter, urlQuearies string)
	GenerateContactsCSV(w http.ResponseWriter, r *http.Request)
	ImportContactsCSV(w http.ResponseWriter, r *http.Request)
}

// actions is the implementation of the Actions interface
type actions struct {
	db     *sql.DB
	config *config.Config
}

// NewActions creates a new action interface
func NewActions(
	db *sql.DB,
	config *config.Config,
) Actions {
	return &actions{
		db:     db,
		config: config,
	}
}

// NotFound action that returns a StatusNotFound
func (a *actions) NotFound(w http.ResponseWriter) {
	a.handleError(w, errors.New(notFound), http.StatusNotFound)
}

// CreateRow action creates a row in the entries database for new contacts
func (a *actions) CreateRow(w http.ResponseWriter, r *http.Request) {
	contact, err := a.getContactFromRequest(r)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
	sqlStatement := fmt.Sprintf(insertInto, a.config.Service.DB)
	id, err := a.doCreateEntry(sqlStatement, contact)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
	a.ReadRows(w, *id)
}

// ReadRows action retreives contact(s) information depending if an id from the entries database
func (a *actions) ReadRows(w http.ResponseWriter, urlQuearies ...string) {
	entries := models.Entries{}
	if len(urlQuearies) == 0 || urlQuearies == nil {
		sqlStatement := fmt.Sprintf(selectFrom, a.config.Service.DB)
		entries = a.doGetEntries(w, sqlStatement, commandQuery)
	} else {
		// TODO: validate id (this could be set to a helper method)
		_, err := strconv.Atoi(urlQuearies[0])
		if err != nil {
			a.handleError(w, errors.New(invalidID), http.StatusBadRequest)
			return
		}

		sqlStatement := fmt.Sprintf(selectFromWhere, a.config.Service.DB, urlQuearies[0])
		entries = a.doGetEntries(w, sqlStatement, commandQueryRow)
	}

	if len(entries.Contacts) > 0 {
		a.encodeJSON(w, entries.Contacts)
		return
	}
	a.encodeJSON(w, fmt.Sprintf(numberOfEntries, len(entries.Contacts)))
}

// UpdateRow action updates a contact in entries database
func (a *actions) UpdateRow(w http.ResponseWriter, r *http.Request) {
	contact, err := a.getContactFromRequest(r)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
	sqlStatement := fmt.Sprintf(update, a.config.Service.DB, contact.FirstName, contact.LastName, contact.Email, contact.Phone, contact.ID)

	err = a.doUpdateEntry(sqlStatement)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
	a.ReadRows(w, contact.ID)
}

// DeleteRow action deletes a contact from the entries database
func (a *actions) DeleteRow(w http.ResponseWriter, urlQuearies string) {
	sqlStatement := fmt.Sprintf(deleteFrom, a.config.Service.DB)
	row, err := a.db.Exec(sqlStatement, urlQuearies)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(row)
}

// GenerateContactsCSV action adapts contacts from entries database to a http response
func (a *actions) GenerateContactsCSV(w http.ResponseWriter, r *http.Request) {
	res, err := a.doExportContacts(w)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
	defer func() {
		res.Close()
		os.Remove(res.Name())
	}()

	w.Header().Set(contentDispositionHeader, csvContentDisposition)
	w.Header().Set(contentTypeHeader, csvContentType)
	w.WriteHeader(http.StatusOK)
	res.Seek(0, 0)
	io.Copy(w, res)
}

// ImportContactsCSV action adapts a csv file from http request and adds the contacts to entries database
func (a *actions) ImportContactsCSV(w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("file")
	if err != nil {
		a.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := handle.Header.Get(contentTypeHeader)
	if mimeType != csvContentType {
		a.handleError(w, errors.New(invalidFileType), http.StatusBadRequest)
		return
	}
	tempFile, err := ioutil.TempFile(os.TempDir(), "tmp.*.csv")
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	_, err = io.Copy(tempFile, file)
	res, err := a.doImportContacts(tempFile)
	if err != nil {
		switch err {
		case errors.New(duplicateEntry):
			errRes := &struct {
				Err            string
				InvalidEntries []*models.Contact
			}{
				Err:            err.Error(),
				InvalidEntries: res,
			}
			a.encodeJSON(w, errRes)
			return
		default:
			a.handleError(w, err, http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusAccepted)
}

// doImportContacts is a helper function that adapts the csv to contacts
func (a *actions) doImportContacts(file *os.File) ([]*models.Contact, error) {
	contacts := []*models.Contact{}
	invalidEntries := []*models.Contact{}
	entryFile, err := os.Open(file.Name())
	if err != nil {
		return nil, err
	}
	defer func() {
		entryFile.Close()
		os.Remove(file.Name())
	}()

	err = gocsv.UnmarshalFile(entryFile, &contacts)
	if err != nil {
		return nil, err
	}
	for _, contact := range contacts {
		if contact.ID != "" {
			sqlStatement := fmt.Sprintf(update, a.config.Service.DB, contact.FirstName, contact.LastName, contact.Email, contact.Phone, contact.ID)
			err = a.doUpdateEntry(sqlStatement)
			if err != nil {
				switch err {
				case errors.New(duplicateEntry):
					invalidEntries = append(invalidEntries, contact)
				default:
					return invalidEntries, err
				}

			}
		} else {
			sqlStatement := fmt.Sprintf(insertInto, a.config.Service.DB)
			_, err = a.doCreateEntry(sqlStatement, *contact)
			if err != nil {
				switch err {
				case errors.New(duplicateEntry):
					invalidEntries = append(invalidEntries, contact)
				default:
					return invalidEntries, err
				}
			}
		}
	}

	if len(invalidEntries) > 0 {
		err = errors.New(duplicateEntry)
	}

	return invalidEntries, err
}

// doExportContacts is a helper function that adapts contacts to a csv file
func (a *actions) doExportContacts(w http.ResponseWriter) (*os.File, error) {
	sqlStatement := fmt.Sprintf(selectFrom, a.config.Service.DB)
	entries := a.doGetEntries(w, sqlStatement, commandQuery)

	contactsFile, err := ioutil.TempFile(os.TempDir(), "tmp.*.csv")
	if err != nil {
		return nil, err
	}
	err = gocsv.MarshalFile(&entries.Contacts, contactsFile)
	if err != nil {
		return nil, err
	}

	return contactsFile, nil
}

// doGetEntries is a helper function for ReadRows action
func (a *actions) doGetEntries(w http.ResponseWriter, sqlStatement, command string) models.Entries {
	entries := models.Entries{}
	var contact models.Contact
	switch command {
	case "Query":
		rows, err := a.db.Query(sqlStatement)
		if err != nil {
			a.handleError(w, err, http.StatusInternalServerError)
			break
		}

		for rows.Next() {
			err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.Email, &contact.Phone)
			if err != nil {
				a.handleError(w, err, http.StatusInternalServerError)
				break
			}
			entries.Contacts = append(entries.Contacts, contact)
		}
		break
	case "QueryRow":
		err := a.db.QueryRow(sqlStatement).Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.Email, &contact.Phone)
		if err != nil {
			a.handleError(w, err, http.StatusBadRequest)
			break
		}
		entries.Contacts = append(entries.Contacts, contact)
		break
	}
	return entries
}

func (a *actions) doUpdateEntry(sqlStatement string) error {
	_, err := a.db.Exec(sqlStatement)
	if err != nil {
		return err
	}
	return nil
}

func (a *actions) doCreateEntry(sqlStatement string, contact models.Contact) (*string, error) {
	var id string
	err := a.db.QueryRow(sqlStatement, contact.FirstName, contact.LastName, contact.Email, contact.Phone).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// getContactFromRequest tries to unmarshal json request into a contact
func (a *actions) getContactFromRequest(r *http.Request) (models.Contact, error) {
	var contact models.Contact
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return contact, err
	}

	err = r.Body.Close()
	if err != nil {
		return contact, err
	}

	if err := json.Unmarshal(body, &contact); err != nil {
		return contact, err
	}
	return contact, nil
}

// encodeJSON is a helper function that encodes data into json for response
func (a *actions) encodeJSON(w http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

// handleError is a helper function that handles errors for http response
func (a *actions) handleError(w http.ResponseWriter, err error, code int) {
	log.Printf("http error: %s (code=%d)", err, code)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&models.HTTPErrorResponse{Error: err.Error()})
}
