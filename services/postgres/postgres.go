package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/squanchersquanch/contacts/services/config"

	// Import blank needed just incase
	_ "github.com/lib/pq"
)

const (
	driverName           = "postgres"
	psqlInfoFormatString = "host='%s' port=%d user='%s' password='%s' dbname='%s' sslmode=disable"
)

// NewDataBase creates a DataBase interface for with appropriated configuration info
func NewDataBase(cfg *config.Config) *sql.DB {
	s := cfg.Service
	psqlInfo := fmt.Sprintf(psqlInfoFormatString,
		s.Host, s.Port, s.User, s.Password, s.Name)
	db, err := sql.Open(driverName, psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected!")
	return db
}
