package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/vincetiu8/penn-spark-server/api/models"
)

// CreateAccessRole creates an access role.
func (s *Server) CreateAccessRole(w http.ResponseWriter, r *http.Request, _ models.User) {
	// Read http body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
	}

	// Unmarshal json into access role struct
	accessRole := models.AccessRole{}
	err = json.Unmarshal(body, &accessRole)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Get mutex lock
	s.Mutex.Lock()

	// Create the access role
	accessRole, err = models.CreateAccessRole(s.DB, accessRole)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, accessRole.ID))
	JSON(w, http.StatusCreated, accessRole)
}

// GetAccessRoleByID gets an access role by its id.
func (s *Server) GetAccessRoleByID(w http.ResponseWriter, r *http.Request, _ models.User) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	s.Mutex.RLock()
	accessRole, err := models.GetAccessRoleByID(s.DB, uint(uid))
	s.Mutex.RUnlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	JSON(w, http.StatusOK, accessRole)
}

// UpdateAccessRole updates an existing access role based on its id.
func (s *Server) UpdateAccessRole(w http.ResponseWriter, r *http.Request, _ models.User) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	accessRole := models.AccessRole{}
	err = json.Unmarshal(body, &accessRole)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	accessRole.ID = uint(uid)

	s.Mutex.Lock()
	accessRole, err = models.UpdateAccessRole(s.DB, accessRole)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	JSON(w, http.StatusOK, accessRole)
}

// DeleteAccessRole deletes an access role based on its id.
func (s *Server) DeleteAccessRole(w http.ResponseWriter, r *http.Request, _ models.User) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	s.Mutex.Lock()
	err = models.DeleteAccessRole(s.DB, uint(uid))
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	JSON(w, http.StatusNoContent, "")
}
