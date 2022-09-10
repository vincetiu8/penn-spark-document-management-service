package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON encodes an interface to a JSON string.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%s", err.Error())
	}
}

// ERROR encodes an error into a JSON string.
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
