package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type FixturesTestStruct struct {
	Field string `json:"field"`
}

func TestMustUnmarshalFromFile(t *testing.T) {
	f := FixturesTestStruct{}
	defer func() {
		if r := recover(); r != nil {
			assert.Fail(t, "json unmarshal panic")
		}
	}()

	MustUnMarshalFromFile("../tests/fixtures/fixtures_test.json", &f)
}
