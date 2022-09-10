package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/vincetiu8/penn-spark-server/api/models"
)

const maxUploadSize = 8 << (10 * 3)

// CreateFileData saves a file to the server, and overwrites existing file data.
// If associated file metadata doesn't exist, the operation will fail.
func (s *Server) CreateFileData(w http.ResponseWriter, r *http.Request, user models.User) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	fileID := uint(fid)

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Read the file in to memory.
	fileData, _, err := r.FormFile("file")
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	defer fileData.Close()

	s.Mutex.Lock()
	// Verify file metadata for this file exists.
	file, err := models.GetFileByID(s.DB, fileID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	if user.ID != file.LastEditorID {
		s.Mutex.Unlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	// Upsert the file.
	err = s.FileSystem.UpsertFileRaw(fileID, fileData)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusInternalServerError, err)
		return
	}

	JSON(w, http.StatusOK, file)
}

// GetFileData gets a file's data based on its id.
func (s *Server) GetFileData(w http.ResponseWriter, r *http.Request, user models.User) {
	// Get file id
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	fileID := uint(fid)

	// Get read lock
	s.Mutex.RLock()

	// Verify user has access rights to file
	file, accessLevel, err := models.GetUserAuthorizationFile(s.DB, user, fileID)
	if err != nil {
		s.Mutex.RUnlock()
		ERROR(w, http.StatusBadRequest, models.ErrFileNotFound)
		return
	} else if accessLevel < models.Viewer {
		s.Mutex.RUnlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	// Allow publishers and file owner to see draft files
	if !file.IsPublished && accessLevel < models.Publisher && file.LastEditorID != user.ID {
		s.Mutex.RUnlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	// Load the file data into memory
	fileData, err := s.FileSystem.GetFileRaw(fileID)
	s.Mutex.RUnlock()
	if err != nil {
		ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Send the file content back to the client
	http.ServeContent(w, r, file.Name, time.Now(), fileData)
}
