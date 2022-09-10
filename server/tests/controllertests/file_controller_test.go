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

func TestCreateFile(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	folders := testServer.Data.Folders

	err = testServer.RefreshTable(&models.File{})
	require.NoError(t, err)
	file := util.Files[0]
	file.FolderID = folders[0].ID

	testCases := []struct {
		inputFile   models.File
		statusCode  int
		expectedErr error
	}{
		{
			inputFile:  file,
			statusCode: http.StatusCreated,
		},
		{
			inputFile: models.File{
				FolderID: folders[2].ID,
				Name:     "file name",
			},
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			inputFile:   file,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFileAlreadyExists,
		},
	}
	for _, testCase := range testCases {
		inputJSON := util.FileToJSON(testCase.inputFile)
		req, err := http.NewRequest("POST", "/files", bytes.NewBufferString(inputJSON))
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		testServer.Server.CreateFile(rr, req, user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusCreated:
				testCase.inputFile.LastEditorID = user.ID
				util.CheckFilesEqual(t, testCase.inputFile, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestGetFileByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	file := testServer.Data.Files[1]
	users := testServer.Data.Users
	_, err = models.UpdateAccessRole(testServer.Server.DB, models.AccessRole{
		ID:          users[0].UserRoles[0].AccessRoles[1].ID,
		AccessLevel: models.Publisher,
	})

	testCases := []struct {
		user        models.User
		fid         uint
		statusCode  int
		expectedErr error
	}{
		{
			user:       users[1],
			fid:        file.ID,
			statusCode: http.StatusOK,
		},
		{
			user:        users[2],
			fid:         file.ID,
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			user:        users[1],
			fid:         999,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFileNotFound,
		},
		{
			user:       users[0],
			fid:        file.ID,
			statusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("GET", "/files", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.fid)})
		rr := httptest.NewRecorder()
		testServer.Server.GetFileByID(rr, req, testCase.user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				util.CheckFilesEqual(t, file, responseMap)

			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestUpdateFile(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	folder := testServer.Data.Folders[0]
	files := testServer.Data.Files
	require.NoError(t, err)

	fileUpdate := models.File{
		Name: "new file name",
	}

	testCases := []struct {
		id          uint
		fileUpdate  models.File
		statusCode  int
		expectedErr error
	}{
		{
			id:         files[0].ID,
			fileUpdate: fileUpdate,
			statusCode: http.StatusOK,
		},
		{
			id:          files[1].ID,
			fileUpdate:  fileUpdate,
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			id:          999,
			fileUpdate:  fileUpdate,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFileNotFound,
		},
	}

	for _, testCase := range testCases {
		updateJSON := util.FileToJSON(testCase.fileUpdate)
		req, err := http.NewRequest("PUT", "/files", bytes.NewBufferString(updateJSON))
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.id)})
		rr := httptest.NewRecorder()
		testServer.Server.UpdateFile(rr, req, user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				fileUpdate.FolderID = folder.ID
				fileUpdate.LastEditorID = user.ID
				util.CheckFilesEqual(t, fileUpdate, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestDeleteFile(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	files := testServer.Data.Files

	testCases := []struct {
		fid         uint
		statusCode  int
		expectedErr error
	}{
		{
			fid:        files[0].ID,
			statusCode: http.StatusNoContent,
		},
		{
			fid:         files[1].ID,
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			fid:         files[0].ID,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFileNotFound,
		},
	}
	for _, testCase := range testCases {
		req, err := http.NewRequest("DELETE", "/files", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.fid)})
		rr := httptest.NewRecorder()
		testServer.Server.DeleteFile(rr, req, user)

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
