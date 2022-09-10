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

// CreateFolder creates a folder.
func (s *Server) CreateFolder(w http.ResponseWriter, r *http.Request, user models.User) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
	}
	folder := models.Folder{}
	err = json.Unmarshal(body, &folder)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// All folders must have a parent folder.
	// The root folder is created manually and all folders are child folders of the root.
	if folder.ParentFolderID == nil {
		ERROR(w, http.StatusBadRequest, models.ErrRequiredParentFolderID)
		return
	}

	s.Mutex.Lock()
	// Verify user is able to create folders in parent folder.
	_, accessLevel, err := models.GetUserAuthorizationFolder(s.DB, user, *folder.ParentFolderID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	} else if accessLevel < models.Publisher {
		s.Mutex.Unlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	folder.LastEditorID = user.ID
	folder, err = models.CreateFolder(s.DB, folder)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, folder.ID))
	JSON(w, http.StatusCreated, folder)
}

// GetFolderByID gets a folder by its id.
func (s *Server) GetFolderByID(w http.ResponseWriter, r *http.Request, user models.User) {
	// Get folder id
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	folderID := uint(fid)

	// Get a read lock
	s.Mutex.RLock()

	// Verify user has access to the folder
	folder, accessLevel, err := models.GetUserAuthorizationFolder(s.DB, user, folderID)
	if err != nil {
		s.Mutex.RUnlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	} else if accessLevel < models.Viewer {
		s.Mutex.RUnlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	// Remove draft files so user can't see them
	if accessLevel < models.Publisher {
		numFiles := len(folder.Files)

		// Counts number of deleted files
		deleted := 0
		for i := 0; i < numFiles; i++ {
			// Need to subtract the number of deleted files from the current index
			// This is because we decrease the length of the slice by 1 each time we delete a file
			// In order for the index to properly line up we make this correction
			file := folder.Files[i-deleted]

			// If user is the uploader (last editor of a draft file), let them see the file
			if !file.IsPublished && file.LastEditorID != user.ID {

				// Delete the file from the slice
				// We swap the file to be deleted with the last element in the slice and then cut the last element
				// This is more efficient from a memory standpoint
				lastIndex := numFiles - deleted - 1

				// No need to swap elements if the file is already the last element
				if lastIndex > 0 {
					folder.Files[i-deleted] = folder.Files[lastIndex]
				}
				folder.Files = folder.Files[:lastIndex]
				deleted += 1
			}
		}
	}

	// Remove child folders user doesn't have access to
	numFolders := len(folder.ChildFolders)
	deleted := 0
	for i := 0; i < numFolders; i++ {
		// Get the access level of the user to the child folder
		_, childAccessLevel, err := models.GetUserAuthorizationFolder(s.DB, user, folder.ChildFolders[i-deleted].ID)
		if err != nil {
			s.Mutex.RUnlock()
			ERROR(w, http.StatusInternalServerError, err)
			return

			// If user doesn't have access to the child folder, remove it from the list of returned folders
		} else if childAccessLevel == models.None {

			// Delete the folder from the slice
			// We swap the folder to be deleted with the last element in the slice and then cut the last element
			// This is more efficient from a memory standpoint
			lastIndex := numFolders - deleted - 1

			// No need to swap elements if the folder is already the last element
			if lastIndex > 0 {
				folder.ChildFolders[i-deleted] = folder.ChildFolders[lastIndex]
			}
			folder.ChildFolders = folder.ChildFolders[:lastIndex]
			deleted += 1
		}
	}
	s.Mutex.RUnlock()

	// We return a custom struct that also contains the access level of the user in the folder
	// This allows the GUI to render appropriate functions based on the access level
	JSON(w, http.StatusOK, struct {
		models.Folder
		AccessLevel models.AccessLevel `json:"access_level"`
	}{
		Folder:      folder,
		AccessLevel: accessLevel,
	})
}

// UpdateFolder updates a folder based on its id.
func (s *Server) UpdateFolder(w http.ResponseWriter, r *http.Request, user models.User) {
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
	folder := models.Folder{}
	err = json.Unmarshal(body, &folder)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	folder.ID = uint(fid)

	s.Mutex.Lock()
	// Verify folder exists.
	currentFolder, err := models.GetFolderByID(s.DB, folder.ID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Verify user has appropriate access rights to the folder.
	_, accessLevel, err := models.GetUserAuthorizationFolder(s.DB, user, *currentFolder.ParentFolderID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusInternalServerError, err)
		return
	} else if accessLevel < models.Publisher {
		s.Mutex.Unlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	// Check user is authorized if folder is being moved to different parent folder.
	if folder.ParentFolderID != nil &&
		*folder.ParentFolderID != *currentFolder.ParentFolderID &&
		*folder.ParentFolderID != 0 {
		_, accessLevel, err = models.GetUserAuthorizationFolder(s.DB, user, *folder.ParentFolderID)
		if err != nil {
			s.Mutex.Unlock()
			ERROR(w, http.StatusInternalServerError, err)
			return
		} else if accessLevel < models.Publisher {
			s.Mutex.Unlock()
			ERROR(w, http.StatusForbidden, ErrUserForbidden)
			return
		}
	}

	folder.LastEditorID = user.ID
	folder, err = models.UpdateFolder(s.DB, folder)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	JSON(w, http.StatusOK, folder)
}

// DeleteFolder deletes a folder by its id.
func (s *Server) DeleteFolder(w http.ResponseWriter, r *http.Request, user models.User) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}
	folderID := uint(fid)

	s.Mutex.Lock()
	currentFolder, err := models.GetFolderByID(s.DB, folderID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Verify the user is authorized to delete the folder.
	_, accessLevel, err := models.GetUserAuthorizationFolder(s.DB, user, *currentFolder.ParentFolderID)
	if err != nil {
		s.Mutex.Unlock()
		ERROR(w, http.StatusBadRequest, err)
		return
	} else if accessLevel < models.Publisher {
		s.Mutex.Unlock()
		ERROR(w, http.StatusForbidden, ErrUserForbidden)
		return
	}

	err = models.DeleteFolder(s.DB, folderID, user.ID)
	s.Mutex.Unlock()
	if err != nil {
		ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", fid))
	JSON(w, http.StatusNoContent, "")
}
