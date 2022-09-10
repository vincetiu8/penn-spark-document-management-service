package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/controllers"
	"github.com/invincibot/penn-spark-server/api/models"
	"github.com/invincibot/penn-spark-server/tests/util"
)

func TestCreateFolder(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	invalidID := uint(999)

	require.NoError(t, err)
	folder := models.Folder{
		Name:           "folder name",
		ParentFolderID: &testServer.Data.Folders[0].ID,
	}

	testCases := []struct {
		inputFolder models.Folder
		statusCode  int
		expectedErr error
	}{
		{
			inputFolder: folder,
			statusCode:  http.StatusCreated,
		},
		{
			inputFolder: models.Folder{
				Name:           folder.Name,
				ParentFolderID: &invalidID,
			},
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFolderNotFound,
		},
		{
			inputFolder: models.Folder{
				Name:           folder.Name,
				ParentFolderID: &testServer.Data.Folders[1].ID,
			},
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			inputFolder: folder,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFolderAlreadyExists,
		},
	}

	for _, testCase := range testCases {
		inputJSON := util.FolderToJSON(testCase.inputFolder)
		req, err := http.NewRequest("POST", "/folders", bytes.NewBufferString(inputJSON))
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		testServer.Server.CreateFolder(rr, req, user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusCreated:
				testCase.inputFolder.LastEditorID = user.ID
				util.CheckFoldersEqual(t, testCase.inputFolder, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestGetFolderByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users
	folders := testServer.Data.Folders
	file := testServer.Data.Files[1]
	_, err = models.UpdateAccessRole(testServer.Server.DB, models.AccessRole{
		ID:          users[0].UserRoles[0].AccessRoles[1].ID,
		AccessLevel: models.Publisher,
	})
	require.NoError(t, err)
	_, err = models.CreateAccessRole(testServer.Server.DB, models.AccessRole{
		UserRoleID:  users[2].UserRoles[0].ID,
		FolderID:    folders[1].ID,
		AccessLevel: models.Uploader,
	})
	require.NoError(t, err)
	_, err = models.UpdateAccessRole(testServer.Server.DB, models.AccessRole{
		ID:          users[2].UserRoles[0].AccessRoles[0].ID,
		AccessLevel: models.None,
	})
	require.NoError(t, err)

	testCases := []struct {
		user        models.User
		fid         uint
		statusCode  int
		expectedErr error
	}{
		{
			user:       users[1],
			fid:        folders[1].ID,
			statusCode: http.StatusOK,
		},
		{
			user:        users[1],
			fid:         folders[0].ID,
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			user:        users[1],
			fid:         999,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFolderNotFound,
		},
		{
			user:       users[0],
			fid:        folders[1].ID,
			statusCode: http.StatusOK,
		},
		{
			user:       users[2],
			fid:        folders[1].ID,
			statusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("GET", "/folders", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.fid)})
		rr := httptest.NewRecorder()
		testServer.Server.GetFolderByID(rr, req, testCase.user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				returnedFolders, ok := responseMap["child_folders"].([]interface{})
				if assert.True(t, ok) {
					if testCase.user.ID != 3 {
						if assert.Len(t, returnedFolders, 1) {
							returnedFolder, ok := returnedFolders[0].(map[string]interface{})
							if assert.True(t, ok) {
								util.CheckFoldersEqual(t, folders[2], returnedFolder)
							}
						}
					} else {
						assert.Len(t, returnedFolders, 0)
					}
				}
				returnedFiles, ok := responseMap["files"].([]interface{})
				if assert.True(t, ok) {
					if testCase.user.ID != 3 {
						if assert.Len(t, returnedFiles, 1) {
							returnedFile, ok := returnedFiles[0].(map[string]interface{})
							if assert.True(t, ok) {
								util.CheckFilesEqual(t, file, returnedFile)
							}
						}
					} else {
						assert.Len(t, returnedFiles, 0)
					}
				}

			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestUpdateFolder(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	folders := testServer.Data.Folders
	accessRole := testServer.Data.UserRoles[3]
	user := testServer.Data.Users[1]
	_, err = models.UpdateAccessRole(testServer.Server.DB, models.AccessRole{
		ID:          accessRole.ID,
		AccessLevel: models.Publisher,
	})
	require.NoError(t, err)

	folderUpdate := models.Folder{
		Name:           "new folder name",
		ParentFolderID: &folders[1].ID,
	}

	testCases := []struct {
		id           uint
		folderUpdate models.Folder
		statusCode   int
		expectedErr  error
	}{
		{
			id:           folders[2].ID,
			folderUpdate: folderUpdate,
			statusCode:   http.StatusOK,
		},
		{
			id:           folders[1].ID,
			folderUpdate: folderUpdate,
			statusCode:   http.StatusForbidden,
			expectedErr:  controllers.ErrUserForbidden,
		},
		{
			id: folders[2].ID,
			folderUpdate: models.Folder{
				Name:           "new folder name",
				ParentFolderID: &folders[0].ID,
			},
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			id:           999,
			folderUpdate: folderUpdate,
			statusCode:   http.StatusBadRequest,
			expectedErr:  models.ErrFolderNotFound,
		},
		{
			id: folders[2].ID,
			folderUpdate: models.Folder{
				Name: "new folder name",
			},
			statusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		updateJSON := util.FolderToJSON(testCase.folderUpdate)
		req, err := http.NewRequest("PUT", "/folders", bytes.NewBufferString(updateJSON))
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.id)})
		rr := httptest.NewRecorder()
		testServer.Server.UpdateFolder(rr, req, user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				folderUpdate.ID = folders[1].ID
				folderUpdate.LastEditorID = user.ID
				util.CheckFoldersEqual(t, folderUpdate, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestDeleteFolder(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users
	folders := testServer.Data.Folders

	err = models.DeleteAccessRole(testServer.Server.DB, testServer.Data.AccessRoles[1].ID)
	require.NoError(t, err)

	err = testServer.RefreshTable(&models.File{})
	require.NoError(t, err)

	testCases := []struct {
		user        models.User
		fid         uint
		statusCode  int
		expectedErr error
	}{
		{
			user:       users[0],
			fid:        folders[2].ID,
			statusCode: http.StatusNoContent,
		},
		{
			user:        users[0],
			fid:         folders[0].ID,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrRequiredFolderID,
		},
		{
			user:        users[1],
			fid:         folders[1].ID,
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			user:        users[0],
			fid:         folders[2].ID,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFolderNotFound,
		},
	}
	for _, testCase := range testCases {
		req, err := http.NewRequest("DELETE", "/folders", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.fid)})
		rr := httptest.NewRecorder()
		testServer.Server.DeleteFolder(rr, req, testCase.user)

		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusBadRequest:
				responseMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}
