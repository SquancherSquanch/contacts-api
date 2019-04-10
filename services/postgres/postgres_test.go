package postgres

import (
	"testing"

	"github.com/squanchersquanch/contacts/services/config"
	"github.com/stretchr/testify/assert"
)

const (
	configFile = "../../development.yaml"
)

func TestNewDataBase(t *testing.T) {
	cfg := config.NewConfig(configFile)
	ps := NewDataBase(cfg)
	assert.NotNil(t, ps)
}

func TestInitDB(t *testing.T) {
	cfg := config.NewConfig(configFile)
	ps := NewDataBase(cfg)
	db := ps.InitDB()
	assert.NotNil(t, db)
}

func TestInitDBPingError(t *testing.T) {
	cfg := config.NewConfig(configFile)
	cfg.Service.Password = ""
	ps := NewDataBase(cfg)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	db := ps.InitDB()
	assert.Nil(t, db)
}
