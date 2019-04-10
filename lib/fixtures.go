package lib

import (
	"encoding/json"
	"io/ioutil"
)

// MustUnMarshalFromFile must unmarshal a json file into an object or panic
func MustUnMarshalFromFile(filePath string, object interface{}) {
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(jsonBytes, object); err != nil {
		panic(err)
	}
	return
}
