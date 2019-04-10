package main

import (
	"log"
	"net/http"

	"github.com/squanchersquanch/contacts/services/config"
	"github.com/squanchersquanch/contacts/services/postgres"
	"github.com/squanchersquanch/contacts/services/router"
)

const (
	configFile = "development.yaml"
)

func main() {

	// load config
	config := config.NewConfig(configFile)

	// load database
	db := postgres.NewDataBase(config)

	//  create a new http client
	router := router.NewRouter(db, config)

	log.Fatal(http.ListenAndServe(":3000", router))
}
