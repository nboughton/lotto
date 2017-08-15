package main

import (
	"encoding/json"
	"net/http"
)

// JSON is the wrapper for all data being passed out of the API
type JSON struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (j JSON) write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(j)
}
