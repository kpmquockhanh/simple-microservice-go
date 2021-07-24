package transhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Json(w http.ResponseWriter, data interface{})  {
	payload, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
	w.WriteHeader(200)
	_, _ = w.Write(payload)
}