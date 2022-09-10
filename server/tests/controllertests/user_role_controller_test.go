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

func TestCreateUserRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	err = testServer.RefreshTable(&models.UserRole{})
	require.NoError(t, err)
	userRole := util.UserRoles[0]
	userRole.Name = "new user role"

	testCases := []struct {
		inputUserRole models.UserRole
		statusCode    int
		expectedErr   error
	}{
		{
			inputUserRole: userRole,
			statusCode:    http.StatusCreated,
		},
		{
			inputUserRole: models.UserRole{
				Name: "",
			},
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrRequiredUserRoleName,
		},
	}

	for _, testCase := range testCases {
		inputJSON := util.UserRoleToJSON(testCase.inputUserRole)
		req, err := http.NewRequest("POST", "/user-roles", bytes.NewBufferString(inputJSON))
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		testServer.Server.CreateUserRole(rr, req, models.User{})

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusCreated:
				util.CheckUserRolesEqual(t, testCase.inputUserRole, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestGetUserRoleByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	userRole := testServer.Data.UserRoles[0]

	testCases := []struct {
		fid         uint
		statusCode  int
		expectedErr error
	}{
		{
			fid:        userRole.ID,
			statusCode: http.StatusOK,
		},
		{
			fid:         999,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrUserRoleNotFound,
		},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("GET", "/user-roles", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.fid)})
		rr := httptest.NewRecorder()
		testServer.Server.GetUserRoleByID(rr, req, models.User{})

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				util.CheckUserRolesEqual(t, userRole, responseMap)

			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestUpdateUserRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	userRole := testServer.Data.UserRoles[0]

	userRoleUpdate := models.UserRole{
		Name: "updated user role name",
	}

	testCases := []struct {
		id             uint
		userRoleUpdate models.UserRole
		statusCode     int
		expectedErr    error
	}{
		{
			id:             userRole.ID,
			userRoleUpdate: userRoleUpdate,
			statusCode:     http.StatusOK,
		},
		{
			id:             999,
			userRoleUpdate: userRoleUpdate,
			statusCode:     http.StatusBadRequest,
			expectedErr:    models.ErrUserRoleNotFound,
		},
	}

	for _, testCase := range testCases {
		updateJSON := util.UserRoleToJSON(testCase.userRoleUpdate)
		req, err := http.NewRequest("PUT", "/user-roles", bytes.NewBufferString(updateJSON))
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.id)})
		rr := httptest.NewRecorder()
		testServer.Server.UpdateUserRole(rr, req, models.User{})

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				util.CheckUserRolesEqual(t, userRoleUpdate, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestDeleteUserRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	userRole := testServer.Data.UserRoles[0]

	testCases := []struct {
		fid         uint
		statusCode  int
		expectedErr error
	}{
		{
			fid:        userRole.ID,
			statusCode: http.StatusNoContent,
		},
		{
			fid:         userRole.ID,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrUserRoleNotFound,
		},
	}
	for _, testCase := range testCases {
		req, err := http.NewRequest("DELETE", "/user-roles", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.fid)})
		rr := httptest.NewRecorder()
		testServer.Server.DeleteUserRole(rr, req, models.User{})

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
