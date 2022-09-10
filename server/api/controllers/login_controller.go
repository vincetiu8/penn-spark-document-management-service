package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vincetiu8/penn-spark-server/api/auth"
	"github.com/vincetiu8/penn-spark-server/api/models"
)

// Login allows a user to log into the server and receive a JWT token.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	s.Mutex.RLock()
	matchingUser, err := models.LoginUser(s.DB, user)
	s.Mutex.RUnlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	token, err := auth.CreateToken(matchingUser.ID)

	JSON(w, http.StatusOK, map[string]interface{}{
		"token":     token,
		"user_data": matchingUser,
	})
}
