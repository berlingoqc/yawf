package utility

import (
	"encoding/json"
)

// JSONToMap convert a json string to a map of interface
func JSONToMap(data string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	return m, json.Unmarshal([]byte(data), m)
}
