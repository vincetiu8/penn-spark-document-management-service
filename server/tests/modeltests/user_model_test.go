package modeltests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/models"
	"github.com/invincibot/penn-spark-server/tests/util"
)

func checkUsersInformationEqual(t *testing.T, expectedUser, actualUser models.User) {
	assert.Equal(t, expectedUser.Username, actualUser.Username)
	assert.Equal(t, expectedUser.FirstName, actualUser.FirstName)
	assert.Equal(t, expectedUser.LastName, actualUser.LastName)
}

func checkUsersEqual(t *testing.T, expectedUser, actualUser models.User) {
	checkUsersInformationEqual(t, expectedUser, actualUser)
	assert.Equal(t, expectedUser.IsAdmin, actualUser.IsAdmin)
	if assert.Len(t, actualUser.UserRoles, len(expectedUser.UserRoles)) {
		for i := range expectedUser.UserRoles {
			checkUserRolesInformationEqual(t, expectedUser.UserRoles[i], actualUser.UserRoles[i])
		}
	}
}

func TestCreateUser(t *testing.T) {
	err := testServer.RefreshTable(&models.User{})
	require.NoError(t, err)
	newUser := util.Users[0]

	testCases := []struct {
		user        models.User
		expectedErr error
	}{
		{
			user:        newUser,
			expectedErr: nil,
		},
		{
			user: models.User{
				Username:  "",
				IsAdmin:   newUser.IsAdmin,
				FirstName: newUser.FirstName,
				LastName:  newUser.LastName,
				Password:  newUser.Password,
			},
			expectedErr: models.ErrRequiredUserUsername,
		},
		{
			user: models.User{
				Username:  util.Users[1].Username,
				IsAdmin:   newUser.IsAdmin,
				FirstName: "",
				LastName:  newUser.LastName,
				Password:  newUser.Password,
			},
			expectedErr: models.ErrRequiredFirstName,
		},
		{
			user: models.User{
				Username:  util.Users[1].Username,
				IsAdmin:   newUser.IsAdmin,
				FirstName: newUser.FirstName,
				LastName:  "",
				Password:  newUser.Password,
			},
			expectedErr: models.ErrRequiredLastName,
		},
		{
			user: models.User{
				Username:  util.Users[1].Username,
				IsAdmin:   newUser.IsAdmin,
				FirstName: newUser.FirstName,
				LastName:  newUser.LastName,
				Password:  "",
			},
			expectedErr: models.ErrRequiredUserPassword,
		},
		{
			user:        newUser,
			expectedErr: models.ErrUserAlreadyExists,
		},
		{
			user: models.User{
				Username:  util.Users[1].Username,
				IsAdmin:   newUser.IsAdmin,
				FirstName: newUser.FirstName,
				LastName:  newUser.LastName,
				Password:  newUser.Password,
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		actualUser, err := models.CreateUser(testServer.Server.DB, testCase.user)

		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			checkUsersInformationEqual(t, testCase.user, actualUser)
		}
	}
}

func TestGetAllUsers(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	foundUsers, err := models.GetAllUsers(testServer.Server.DB)
	require.NoError(t, err)

	assert.Equal(t, len(testServer.Data.Users), len(foundUsers))
	for i := range testServer.Data.Users {
		checkUsersInformationEqual(t, testServer.Data.Users[i], foundUsers[i])
	}
}

func TestGetUserByID(t *testing.T) {
	err := testServer.SeedData()

	for _, user := range testServer.Data.Users {
		foundUser, err := models.GetUserByID(testServer.Server.DB, user.ID)
		if assert.NoError(t, err) {
			checkUsersEqual(t, user, foundUser)
		}
	}

	_, err = models.GetUserByID(testServer.Server.DB, 0)
	assert.Equal(t, models.ErrRequiredUserID, err)

	_, err = models.GetUserByID(testServer.Server.DB, 999)
	assert.Equal(t, models.ErrUserNotFound, err)
}

func TestUpdateUser(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	userUpdate := models.User{
		Model: models.Model{
			ID: testServer.Data.Users[0].ID,
		},
		Username:  "new username",
		IsAdmin:   false,
		FirstName: "new first name",
		LastName:  "new last name",
		Password:  "new password",
		UserRoles: []models.UserRole{
			testServer.Data.UserRoles[0],
		},
	}
	expectedUserWithoutAdmin := userUpdate
	expectedUserWithoutAdmin.IsAdmin = true

	testCases := []struct {
		userUpdate   models.User
		isAdmin      bool
		expectedUser models.User
		expectedErr  error
	}{
		{
			userUpdate:   userUpdate,
			expectedErr:  nil,
			expectedUser: expectedUserWithoutAdmin,
		},
		{
			userUpdate: models.User{
				Model: models.Model{
					ID: userUpdate.ID,
				},
				Username:  "",
				FirstName: userUpdate.FirstName,
				LastName:  userUpdate.LastName,
				Password:  userUpdate.Password,
				IsAdmin:   userUpdate.IsAdmin,
				UserRoles: userUpdate.UserRoles,
			},
			expectedUser: expectedUserWithoutAdmin,
		},
		{
			userUpdate: models.User{
				Model: models.Model{
					ID: userUpdate.ID,
				},
				Username:  testServer.Data.Users[1].Username,
				FirstName: userUpdate.FirstName,
				LastName:  userUpdate.LastName,
				Password:  userUpdate.Password,
				IsAdmin:   userUpdate.IsAdmin,
				UserRoles: userUpdate.UserRoles,
			},
			expectedErr: models.ErrUserAlreadyExists,
		},
		{
			userUpdate: models.User{
				Model: models.Model{
					ID: userUpdate.ID,
				},
				Username:  userUpdate.Username,
				FirstName: "",
				LastName:  userUpdate.LastName,
				Password:  userUpdate.Password,
				IsAdmin:   userUpdate.IsAdmin,
				UserRoles: userUpdate.UserRoles,
			},
			expectedUser: expectedUserWithoutAdmin,
		},
		{
			userUpdate: models.User{
				Model: models.Model{
					ID: userUpdate.ID,
				},
				Username:  userUpdate.Username,
				FirstName: userUpdate.FirstName,
				LastName:  "",
				Password:  userUpdate.Password,
				IsAdmin:   userUpdate.IsAdmin,
				UserRoles: userUpdate.UserRoles,
			},
			expectedUser: expectedUserWithoutAdmin,
		},
		{
			userUpdate: models.User{
				Model: models.Model{
					ID: userUpdate.ID,
				},
				Username:  userUpdate.Username,
				FirstName: userUpdate.FirstName,
				LastName:  userUpdate.LastName,
				Password:  "",
				IsAdmin:   userUpdate.IsAdmin,
				UserRoles: userUpdate.UserRoles,
			},
			expectedUser: expectedUserWithoutAdmin,
		},
		{
			userUpdate: models.User{
				Model: models.Model{
					ID: 0,
				},
				Username:  userUpdate.Username,
				FirstName: userUpdate.FirstName,
				LastName:  userUpdate.LastName,
				Password:  userUpdate.Password,
				IsAdmin:   userUpdate.IsAdmin,
				UserRoles: userUpdate.UserRoles,
			},
			expectedErr: models.ErrRequiredUserID,
		},
		{
			userUpdate: models.User{
				Model: models.Model{
					ID: 999,
				},
				Username:  userUpdate.Username,
				FirstName: userUpdate.FirstName,
				LastName:  userUpdate.LastName,
				Password:  userUpdate.Password,
				IsAdmin:   userUpdate.IsAdmin,
				UserRoles: userUpdate.UserRoles,
			},
			expectedErr: models.ErrUserNotFound,
		},
		{
			userUpdate: userUpdate,
			isAdmin:    true,
		},
		{
			userUpdate: models.User{
				Model: models.Model{
					ID: userUpdate.ID,
				},
				Username:  userUpdate.Username,
				FirstName: userUpdate.FirstName,
				LastName:  userUpdate.LastName,
				Password:  userUpdate.Password,
				IsAdmin:   userUpdate.IsAdmin,
				UserRoles: []models.UserRole{},
			},
			expectedUser: userUpdate,
		},
	}

	for _, testCase := range testCases {
		actualUser, err := models.UpdateUser(testServer.Server.DB, testCase.userUpdate, testCase.isAdmin)
		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			if testCase.expectedUser.Username != "" {
				checkUsersEqual(t, testCase.expectedUser, actualUser)
			} else {
				checkUsersEqual(t, testCase.userUpdate, actualUser)
			}
		}
	}
}

func TestDeleteUser(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	err = models.DeleteUser(testServer.Server.DB, testServer.Data.Users[0].ID)
	require.NoError(t, err)

	err = models.DeleteUser(testServer.Server.DB, testServer.Data.Users[0].ID)
	require.Equal(t, models.ErrUserNotFound, err)
}

func TestLoginUser(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)
	user := util.Users[0]

	testCases := []struct {
		login       models.User
		expectedErr error
	}{
		{
			login:       user,
			expectedErr: nil,
		},
		{
			login: models.User{
				Username: "",
				Password: user.Password,
			},
			expectedErr: models.ErrRequiredUserUsername,
		},
		{
			login: models.User{
				Username: "incorrect username",
				Password: user.Password,
			},
			expectedErr: models.ErrUserNotFound,
		},
		{
			login: models.User{
				Username: user.Username,
				Password: "",
			},
			expectedErr: models.ErrRequiredUserPassword,
		},
		{
			login: models.User{
				Username: user.Username,
				Password: "incorrect password",
			},
			expectedErr: models.ErrIncorrectPassword,
		},
	}

	for _, testCase := range testCases {
		actualUser, err := models.LoginUser(testServer.Server.DB, testCase.login)
		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			checkUsersEqual(t, user, actualUser)
		}
	}
}
