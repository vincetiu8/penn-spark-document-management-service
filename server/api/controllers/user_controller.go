package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/vincetiu8/penn-spark-server/api/models"
)

// CreateUser creates a user.
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request, _ models.User) {
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

	s.Mutex.Lock()
	user, err = models.CreateUser(s.DB, user)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, user.ID))
	JSON(w, http.StatusCreated, user)
}

// GetAllUsers returns a list of all activated users.
func (s *Server) GetAllUsers(w http.ResponseWriter, _ *http.Request, _ models.User) {
	s.Mutex.RLock()
	users, err := models.GetAllUsers(s.DB)
	s.Mutex.RUnlock()
	if err != nil {
		ERROR(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, users)
}

// GetUserByID gets a user by their id.
func (s *Server) GetUserByID(w http.ResponseWriter, r *http.Request, user models.User) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	userID := uint(uid)

	if userID == user.ID {
		JSON(w, http.StatusOK, user)
		return
	}

	s.Mutex.RLock()
	foundUser, err := models.GetUserByID(s.DB, userID)
	s.Mutex.RUnlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	JSON(w, http.StatusOK, foundUser)
}

// UpdateUser updates a user based on their id.
func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request, user models.User) {
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
	userUpdate := models.User{}
	err = json.Unmarshal(body, &userUpdate)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userUpdate.Model = models.Model{
		ID: uint(uid),
	}

	// Non-admin users should only be able to update their own profile.
	if !user.IsAdmin && user.ID != userUpdate.ID {
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	s.Mutex.Lock()
	userUpdate, err = models.UpdateUser(s.DB, userUpdate, user.IsAdmin)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	JSON(w, http.StatusOK, userUpdate)
}

// DeleteUser deletes a user by their id.
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request, _ models.User) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	s.Mutex.Lock()
	err = models.DeleteUser(s.DB, uint(uid))
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	JSON(w, http.StatusNoContent, "")
}

// AddUserRole adds a user role to a user.
func (s *Server) AddUserRole(w http.ResponseWriter, r *http.Request, _ models.User) {
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

	s.Mutex.Lock()
	userUpdate, err := models.AddUserRole(s.DB, uint(uid), userRole)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	JSON(w, http.StatusOK, userUpdate)
}

// RemoveUserRole removes a user role from a user.
func (s *Server) RemoveUserRole(w http.ResponseWriter, r *http.Request, _ models.User) {
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

	s.Mutex.Lock()
	userUpdate, err := models.RemoveUserRole(s.DB, uint(uid), userRole)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	JSON(w, http.StatusOK, userUpdate)
}

// Reactivate user reactivates a deleted user by their username.
func (s *Server) ReactivateUser(w http.ResponseWriter, r *http.Request, _ models.User) {
	// Read http body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Unmarshal json into user struct
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Invalidate user deletion, effectively reactivating them
	// This will allow the user to be found in the database
	user = models.User{
		Model: models.Model{
			DeletedAt: gorm.DeletedAt{
				Valid: false,
			},
		},
		Username: user.Username,
	}

	// Get mutex lock
	s.Mutex.Lock()

	// Update the user's details accordingly
	user, err = models.UpdateUser(s.DB, user, true)
	s.Mutex.Unlock()

	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	JSON(w, http.StatusOK, user)
}
