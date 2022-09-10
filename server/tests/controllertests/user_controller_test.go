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

func TestCreateUser(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	err = testServer.RefreshTable(&models.User{})
	require.NoError(t, err)
	user := util.Users[0]
	user.UserRoles = []models.UserRole{testServer.Data.UserRoles[0]}

	testCases := []struct {
		inputUser   models.User
		statusCode  int
		expectedErr error
	}{
		{
			inputUser:  user,
			statusCode: http.StatusCreated,
		},
		{
			inputUser:   user,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrUserAlreadyExists,
		},
		{
			inputUser: models.User{
				Username:  "invincibot",
				Password:  user.Password,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				IsAdmin:   user.IsAdmin,
			},
			statusCode: http.StatusCreated,
		},
	}
	for _, testCase := range testCases {
		inputJSON := util.UserToJSON(testCase.inputUser, "")
		req, err := http.NewRequest("POST", "/users/admin", bytes.NewBufferString(inputJSON))
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		testServer.Server.CreateUser(rr, req, models.User{})

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusCreated:
				util.CheckUsersEqual(t, testCase.inputUser, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestGetUsers(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users

	req, err := http.NewRequest("GET", "/users/admin", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	testServer.Server.GetUsers(rr, req, models.User{})

	var returnedUsers []map[string]interface{}
	err = json.Unmarshal([]byte(rr.Body.String()), &returnedUsers)
	require.NoError(t, err)
	require.Len(t, returnedUsers, len(users))

	for i, user := range users {
		util.CheckUsersEqual(t, user, returnedUsers[i])
	}
}

func TestGetUserByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users

	testCases := []struct {
		uid         uint
		user        models.User
		statusCode  int
		expectedErr error
	}{
		{
			uid:        users[0].ID,
			user:       users[0],
			statusCode: http.StatusOK,
		},
		{
			uid:        users[1].ID,
			user:       users[0],
			statusCode: http.StatusOK,
		},
		{
			uid:        users[1].ID,
			user:       users[1],
			statusCode: http.StatusOK,
		},
		{
			uid:         users[0].ID,
			user:        users[1],
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			uid:         999,
			statusCode:  http.StatusForbidden,
			expectedErr: models.ErrUserNotFound,
		},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("GET", "/users", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.uid)})
		rr := httptest.NewRecorder()
		testServer.Server.GetUserByID(rr, req, testCase.user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				util.CheckUsersEqual(t, users[testCase.uid-1], responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestUpdateUser(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users
	userUpdate := models.User{
		FirstName: "new first name",
		IsAdmin:   true,
	}

	testCases := []struct {
		id           uint
		user         models.User
		userUpdate   models.User
		statusCode   int
		expectedErr  error
		expectedUser models.User
	}{
		{
			id:         users[0].ID,
			user:       users[0],
			userUpdate: userUpdate,
			statusCode: http.StatusOK,
			expectedUser: models.User{
				FirstName: userUpdate.FirstName,
				LastName:  users[0].LastName,
				Username:  users[0].Username,
				IsAdmin:   userUpdate.IsAdmin,
			},
		},
		{
			id:         users[1].ID,
			user:       users[1],
			userUpdate: userUpdate,
			statusCode: http.StatusOK,
			expectedUser: models.User{
				FirstName: userUpdate.FirstName,
				LastName:  users[1].LastName,
				Username:  users[1].Username,
				IsAdmin:   users[1].IsAdmin,
			},
		},
		{
			id:          users[0].ID,
			user:        users[1],
			userUpdate:  userUpdate,
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			id:          999,
			userUpdate:  userUpdate,
			statusCode:  http.StatusForbidden,
			expectedErr: models.ErrUserNotFound,
		},
		{
			id:         users[1].ID,
			user:       users[0],
			userUpdate: userUpdate,
			statusCode: http.StatusOK,
			expectedUser: models.User{
				FirstName: userUpdate.FirstName,
				LastName:  users[1].LastName,
				Username:  users[1].Username,
				IsAdmin:   userUpdate.IsAdmin,
			},
		},
	}

	for _, testCase := range testCases {
		updateJSON := util.UserToJSON(testCase.userUpdate, "update")
		req, err := http.NewRequest("PUT", "/users", bytes.NewBufferString(updateJSON))
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.id)})
		rr := httptest.NewRecorder()
		testServer.Server.UpdateUser(rr, req, testCase.user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				util.CheckUserInformationEqual(t, testCase.expectedUser, responseMap)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestDeleteUser(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]

	testCases := []struct {
		uid         uint
		statusCode  int
		expectedErr error
	}{
		{
			uid:        user.ID,
			statusCode: http.StatusNoContent,
		},
		{
			uid:         user.ID,
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrUserNotFound,
		},
	}
	for _, testCase := range testCases {
		req, err := http.NewRequest("DELETE", "/users/admin", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.uid)})
		rr := httptest.NewRecorder()
		testServer.Server.DeleteUser(rr, req, models.User{})

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
