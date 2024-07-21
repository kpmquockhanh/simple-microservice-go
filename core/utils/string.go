package utils

import "encoding/json"

func ToJsonString(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}
