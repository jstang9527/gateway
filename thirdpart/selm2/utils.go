package selm2

import (
	"encoding/json"
)

// ToJSON ...
func ToJSON(v interface{}) (js string) {
	b, _ := json.Marshal(v)
	return string(b)
}
