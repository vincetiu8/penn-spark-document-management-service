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

	"github.com/invincibot/penn-spark-server/api/models"
	"github.com/invincibot/penn-spark-server/tests/util"
)

func TestCreateAccessRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	folder := testServer.Data.Folders[0]
	userRole := testServer.Data.UserRoles[0]

	err = testServer.RefreshTable(&models.AccessRole{})
	require.NoError(t, err)
	accessRole := util.AccessRoles[0]
	accessRole.FolderID = folder.ID
	accessRole.UserRoleID = userRole.ID

	testCases := []struct {
		inputAccessRole models.AccessRole
		statusCode      int
		expectedErr     error
	}{
		{
			inputAccessRole: accessRole,
			statusCode:      http.StatusCreated,
		},
		{
			inputAccessRole: models.AccessRole{
				FolderID:    0,
				AccessLevel: models.Publisher,
			},
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrRequiredFolderID,
		},
	}

	for _, testCase := range testCases {
		inputJSON := util.AccessRoleToJSON(testCase.inputAccessRole)
		req, err := http.NewRequest("POST", "/access-roles", bytes.NewBufferString(inputJSON))
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		testServer.Server.CreateAccessRole(rr, req, models.User{})

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusCreated:
				util.CheckAccessRolesEqual(t, testCase.inputAccessRole, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestGetAccessRoleByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	accessRole := testServer.Data.AccessRoles[0]

	testCases := []struct {
		fid         uint
		statusCode  int
		expectedErr error
	}{
		{
			fid:        accessRole.ID,
			statusCode: http.StatusOK,
		},
		{
			fid:         999,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrAccessRoleNotFound,
		},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("GET", "/access-roles", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.fid)})
		rr := httptest.NewRecorder()
		testServer.Server.GetAccessRoleByID(rr, req, models.User{})

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				util.CheckAccessRolesEqual(t, accessRole, responseMap)

			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestUpdateAccessRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	folder := testServer.Data.Folders[0]
	userRole := testServer.Data.UserRoles[2]
	accessRole := testServer.Data.AccessRoles[0]

	accessRoleUpdate := models.AccessRole{
		FolderID:    folder.ID,
		UserRoleID:  userRole.ID,
		AccessLevel: models.Viewer,
	}

	testCases := []struct {
		id               uint
		accessRoleUpdate models.AccessRole
		statusCode       int
		expectedErr      error
	}{
		{
			id:               accessRole.ID,
			accessRoleUpdate: accessRoleUpdate,
			statusCode:       http.StatusOK,
		},
		{
			id:               999,
			accessRoleUpdate: accessRoleUpdate,
			statusCode:       http.StatusBadRequest,
			expectedErr:      models.ErrAccessRoleNotFound,
		},
	}

	for _, testCase := range testCases {
		updateJSON := util.AccessRoleToJSON(testCase.accessRoleUpdate)
		req, err := http.NewRequest("PUT", "/access-roles", bytes.NewBufferString(updateJSON))
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.id)})
		rr := httptest.NewRecorder()
		testServer.Server.UpdateAccessRole(rr, req, models.User{})

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				util.CheckAccessRolesEqual(t, accessRoleUpdate, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestDeleteAccessRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	accessRole := testServer.Data.AccessRoles[0]

	testCases := []struct {
		fid         uint
		statusCode  int
		expectedErr error
	}{
		{
			fid:        accessRole.ID,
			statusCode: http.StatusNoContent,
		},
		{
			fid:         accessRole.ID,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrAccessRoleNotFound,
		},
	}
	for _, testCase := range testCases {
		req, err := http.NewRequest("DELETE", "/access-roles", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.fid)})
		rr := httptest.NewRecorder()
		testServer.Server.DeleteAccessRole(rr, req, models.User{})

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
