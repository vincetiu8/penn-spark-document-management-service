package controllers

import (
	"errors"
	"net/http"

	"github.com/vincetiu8/penn-spark-server/api/auth"
	"github.com/vincetiu8/penn-spark-server/api/models"
)

// ErrUserUnauthorized is returned if a user has an invalid token.
var ErrUserUnauthorized = errors.New("unauthorized")

// ErrUserForbidden is returned if a user has insufficient access rights.
var ErrUserForbidden = errors.New("forbidden")

// controllerFunc is a wrapper around http.HandlerFunc allowing user information to be passed to the handlers.
type controllerFunc func(w http.ResponseWriter, r *http.Request, user models.User)

// SetMiddlewareJSON sets the header of the outgoing response.
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication authenticates a user from their token.
// It will also verify admins for admin-only functions.
func SetMiddlewareAuthentication(controllerFunc controllerFunc, s *Server, isAdmin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			ERROR(w, http.StatusUnauthorized, ErrUserUnauthorized)
			return
		}

		uid, err := auth.ExtractTokenID(r)
		if err != nil {
			ERROR(w, http.StatusUnauthorized, ErrUserUnauthorized)
			return
		}

		s.Mutex.RLock()
		user, err := models.GetUserByID(s.DB, uid)
		s.Mutex.RUnlock()
		if err != nil {
			ERROR(w, http.StatusUnauthorized, ErrUserUnauthorized)
			return
		}

		if isAdmin && !user.IsAdmin {
			ERROR(w, http.StatusForbidden, ErrUserForbidden)
			return
		}

		controllerFunc(w, r, user)
	}
}
