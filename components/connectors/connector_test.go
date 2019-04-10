package connectors

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"

	"testing"

	"github.com/squanchersquanch/contacts/components/actions"
	"github.com/squanchersquanch/contacts/services/config"
	"github.com/squanchersquanch/contacts/services/logger"
	"github.com/squanchersquanch/contacts/services/postgres"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

const (
	configFile            = "../../development.yaml"
	newContactFilePath    = "../../tests/fixtures/connectors/new_contact.json"
	updateContactFilePath = "../../tests/fixtures/connectors/update_contact.json"
	testContactsFilePath  = "../../tests/fixtures/test_import_contacts.csv"
	testImportFile        = "test_import_contacts.csv"
)

type routes interface {
	RouteList() []route
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type connectorSuite struct {
	suite.Suite
	db      *sql.DB
	actions actions.Actions

	router    *mux.Router
	routes    []route
	connector Connector
}

func TestConnectorSuite(t *testing.T) {
	suite.Run(t, &connectorSuite{})
}

func (s *connectorSuite) SetupTest() {
	config := config.NewConfig(configFile)

	s.db = postgres.NewDataBase(config)

	s.actions = actions.NewActions(s.db)

	s.connector = &connector{
		actions: s.actions,
	}
	s.router = mux.NewRouter().StrictSlash(true)
	s.routes = []route{
		route{
			"CreateContact",
			"POST",
			"/api/entry",
			s.connector.CreateContact,
		},
		route{
			"GetContact",
			"GET",
			"/api/entry",
			s.connector.GetContacts,
		},
		route{
			"UpdateContact",
			"PUT",
			"/api/entry",
			s.connector.UpdateContact,
		},
		route{
			"DeleteContact",
			"DELETE",
			"/api/entry",
			s.connector.DeleteContact,
		},
		route{
			"ExportContacts",
			"GET",
			"/api/entry/export",
			s.connector.ExportContacts,
		},
		route{
			"ImportContacts",
			"POST",
			"/api/entry/import",
			s.connector.ImportContacts,
		},
		route{
			"NotFound",
			"",
			"",
			s.connector.NotFound,
		},
	}
	for _, route := range s.routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		s.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}

func (s *connectorSuite) TestNotFound() {
	req, err := http.NewRequest("GET", "/api/failhard", nil)
	s.NoError(err)
	s.NotNil(req)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(rr.Code, http.StatusNotFound)
}

func (s *connectorSuite) TestCreateContact() {
	data, err := ioutil.ReadFile(newContactFilePath)
	s.NoError(err)
	req, err := http.NewRequest("POST", "/api/entry", bytes.NewBuffer(data))
	s.NoError(err)
	s.NotNil(req)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(rr.Code, http.StatusOK)
}

func (s *connectorSuite) TestGetContacts() {
	req, err := http.NewRequest("GET", "/api/entry", nil)
	s.NoError(err)
	s.NotNil(req)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(rr.Code, http.StatusOK)

	req, err = http.NewRequest("GET", "/api/entry?id=2", nil)
	s.NoError(err)
	s.NotNil(req)

	rr = httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(rr.Code, http.StatusOK)
}

func (s *connectorSuite) TestUpdateContact() {
	data, err := ioutil.ReadFile(updateContactFilePath)
	s.NoError(err)
	req, err := http.NewRequest("PUT", "/api/entry", bytes.NewBuffer(data))
	s.NoError(err)
	s.NotNil(req)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(rr.Code, http.StatusOK)
}

func (s *connectorSuite) TestDeleteContact() {
	req, err := http.NewRequest("DELETE", "/api/entry?id=2", nil)
	s.NoError(err)
	s.NotNil(req)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(rr.Code, http.StatusOK)
}

func (s *connectorSuite) TestImportContacts() {
	req, err := http.NewRequest("GET", "/api/entry/export", nil)
	s.NoError(err)
	s.NotNil(req)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(rr.Code, http.StatusOK)
}

func (s *connectorSuite) TestExportContacts() {
	testImportCSV, err := os.Open(testContactsFilePath)
	s.NoError(err)
	defer testImportCSV.Close()

	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)
	formFile, _ := generateCSVFile(bodyWriter, testImportFile)
	io.Copy(formFile, testImportCSV)
	bodyWriter.Close()

	req, err := http.NewRequest("POST", "/api/entry/import", bodyBuffer)
	s.NoError(err)
	s.NotNil(req)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(rr.Code, http.StatusAccepted)
}

func generateCSVFile(w *multipart.Writer, filename string) (io.Writer, error) {
	mh := make(textproto.MIMEHeader)
	mh.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	mh.Set("Content-Type", "text/csv")
	return w.CreatePart(mh)
}
