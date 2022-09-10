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

// CreateFile creates a file metadata reference.
func (s *Server) CreateFile(w http.ResponseWriter, r *http.Request, user models.User) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
	}
	file := models.File{}
	err = json.Unmarshal(body, &file)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	s.Mutex.Lock()
	// Verify the user has sufficient permissions to create the file.
	_, accessLevel, err := models.GetUserAuthorizationFolder(s.DB, user, file.FolderID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	} else if accessLevel < models.Uploader {
		s.Mutex.Unlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	file.LastEditorID = user.ID
	file.IsPublished = false
	file, err = models.CreateFile(s.DB, file)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, file.ID))
	JSON(w, http.StatusCreated, file)
}

// GetFileByID gets a file's metadata based on its id.
func (s *Server) GetFileByID(w http.ResponseWriter, r *http.Request, user models.User) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	fileID := uint(fid)

	s.Mutex.RLock()
	// Verify the user has access rights to the file.
	file, accessLevel, err := models.GetUserAuthorizationFile(s.DB, user, fileID)
	s.Mutex.RUnlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	} else if accessLevel < models.Viewer && file.LastEditorID != user.ID {
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	if !file.IsPublished && file.LastEditorID != user.ID && accessLevel < models.Publisher {
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	JSON(w, http.StatusOK, file)
}

// UpdateFile updates a file's metadata based on its id.
func (s *Server) UpdateFile(w http.ResponseWriter, r *http.Request, user models.User) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	file := models.File{}
	err = json.Unmarshal(body, &file)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	file.ID = uint(fid)

	s.Mutex.Lock()
	// Check if user is authorized to update the file.
	currentFile, accessLevel, err := models.GetUserAuthorizationFile(s.DB, user, file.ID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	} else if accessLevel < models.Publisher {
		s.Mutex.Unlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	// Check if user is authorized if the file is being moved to another folder.
	if file.FolderID != currentFile.FolderID && file.FolderID != 0 {
		_, accessLevel, err = models.GetUserAuthorizationFolder(s.DB, user, file.FolderID)
		if err != nil {
			s.Mutex.Unlock()
			ERROR(w, http.StatusBadRequest, err)
			return
		} else if accessLevel < models.Publisher {
			s.Mutex.Unlock()
			ERROR(w, http.StatusForbidden, ErrUserForbidden)
			return
		}
	}

	file.LastEditorID = user.ID
	file, err = models.UpdateFile(s.DB, file)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	JSON(w, http.StatusOK, file)
}

// DeleteFile deletes a file's metadata by its id.
func (s *Server) DeleteFile(w http.ResponseWriter, r *http.Request, user models.User) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	fileID := uint(fid)

	s.Mutex.Lock()
	// Verify the user is authorized to delete the file.
	_, accessLevel, err := models.GetUserAuthorizationFile(s.DB, user, fileID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, models.ErrFileNotFound)
		return
	} else if accessLevel < models.Publisher {
		s.Mutex.Unlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	err = models.DeleteFile(s.DB, fileID, user.ID)
	s.Mutex.Unlock()
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", fid))
	JSON(w, http.StatusNoContent, "")
}
