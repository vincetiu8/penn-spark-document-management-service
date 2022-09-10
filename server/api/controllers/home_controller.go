package controllers

import (
	"net/http"
)

// Home allows a user to ping the server to check for connectivity.
func (s *Server) Home(w http.ResponseWriter, _ *http.Request) {
	JSON(w, http.StatusOK, "SBA")
}
