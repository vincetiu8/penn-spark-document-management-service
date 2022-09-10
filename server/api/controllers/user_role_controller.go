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

// CreateUserRole creates a user role.
func (s *Server) CreateUserRole(w http.ResponseWriter, r *http.Request, _ models.User) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
	}

	userRole := models.UserRole{}
	err = json.Unmarshal(body, &userRole)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	s.Mutex.Lock()
	userRole, err = models.CreateUserRole(s.DB, userRole)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userRole.ID))
	JSON(w, http.StatusCreated, userRole)
}

// GetAllUserRoles returns a list of all the user roles.
func (s *Server) GetAllUserRoles(w http.ResponseWriter, _ *http.Request, _ models.User) {
	s.Mutex.RLock()
	users, err := models.GetAllUserRoles(s.DB)
	s.Mutex.RUnlock()
	if err != nil {
		ERROR(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, users)
}

// GetUserRolesByID gets a user role by its id.
func (s *Server) GetUserRoleByID(w http.ResponseWriter, r *http.Request, _ models.User) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	s.Mutex.RLock()
	userRole, err := models.GetUserRoleByID(s.DB, uint(uid))
	s.Mutex.RUnlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	JSON(w, http.StatusOK, userRole)
}

// UpdateUserRole updates a user role based on its id.
func (s *Server) UpdateUserRole(w http.ResponseWriter, r *http.Request, _ models.User) {
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

	userRole := models.UserRole{}
	err = json.Unmarshal(body, &userRole)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userRole.ID = uint(uid)

	s.Mutex.Lock()
	userRole, err = models.UpdateUserRole(s.DB, userRole)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	JSON(w, http.StatusOK, userRole)
}

// DeleteUserRole deletes a user role by its id.
func (s *Server) DeleteUserRole(w http.ResponseWriter, r *http.Request, _ models.User) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	s.Mutex.Lock()
	err = models.DeleteUserRole(s.DB, uint(uid))
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	JSON(w, http.StatusNoContent, "")
}
