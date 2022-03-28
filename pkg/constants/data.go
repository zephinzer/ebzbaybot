package constants

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data.json
var dataSource []byte
var data Data

type Data struct {
	Collections collections `json:"collections"`
}
type collections map[string]string

var Collection map[string]string

func init() {
	if err := json.Unmarshal(dataSource, &data); err != nil {
		panic(fmt.Errorf("failed to unmarshal data.json into data: %s", err))
	}
	Collection = data.Collections
}
