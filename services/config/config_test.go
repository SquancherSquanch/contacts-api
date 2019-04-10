package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	configFile  = "../../development.yaml"
	errorConfig = "../../tests/fixtures/errorConfig.yaml"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig(configFile)
	assert.NotNil(t, cfg)
}

func TestNewConfigError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_ = NewConfig("")
}
