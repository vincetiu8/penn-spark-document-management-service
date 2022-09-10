package modeltests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/models"
)

func checkUserRolesInformationEqual(t *testing.T, expectedRole, actualRole models.UserRole) {
	assert.Equal(t, expectedRole.Name, actualRole.Name)
}

func checkUserRolesEqual(t *testing.T, expectedRole, actualRole models.UserRole) {
	checkUserRolesInformationEqual(t, expectedRole, actualRole)
	if assert.Len(t, actualRole.AccessRoles, len(expectedRole.AccessRoles)) {
		for i := range expectedRole.AccessRoles {
			checkAccessRolesEqual(t, expectedRole.AccessRoles[i], actualRole.AccessRoles[i])
		}
	}
}

func TestCreateUserRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	accessRoles := testServer.Data.AccessRoles

	err = testServer.RefreshTable(&models.UserRole{})
	require.NoError(t, err)

	testCases := []struct {
		userRole    models.UserRole
		expectedErr error
	}{
		{
			userRole: models.UserRole{
				Name:        "user role",
				AccessRoles: accessRoles,
			},
		},
		{
			userRole: models.UserRole{
				Name:        "",
				AccessRoles: accessRoles,
			},
			expectedErr: models.ErrRequiredUserRoleName,
		},
		{
			userRole: models.UserRole{
				Name:        "user role",
				AccessRoles: accessRoles,
			},
			expectedErr: models.ErrUserRoleAlreadyExists,
		},
		{
			userRole: models.UserRole{
				Name:        "different user role",
				AccessRoles: nil,
			},
		},
	}

	for _, testCase := range testCases {
		createdUserRole, err := models.CreateUserRole(testServer.Server.DB, testCase.userRole)
		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			checkUserRolesEqual(t, testCase.userRole, createdUserRole)
		}
	}
}

func TestGetUserRoleByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	for _, userRole := range testServer.Data.UserRoles {
		foundRole, err := models.GetUserRoleByID(testServer.Server.DB, userRole.ID)
		if assert.NoError(t, err) {
			checkUserRolesEqual(t, userRole, foundRole)
		}
	}

	_, err = models.GetUserRoleByID(testServer.Server.DB, 0)
	assert.Equal(t, models.ErrRequiredUserRoleID, err)

	_, err = models.GetUserRoleByID(testServer.Server.DB, 999)
	assert.Equal(t, models.ErrUserRoleNotFound, err)
}

func TestUpdateUserRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	userRoles := testServer.Data.UserRoles
	accessRoles := testServer.Data.AccessRoles

	accessRolesUpdate := []models.AccessRole{
		accessRoles[0],
	}

	testCases := []struct {
		roleUpdate   models.UserRole
		expectedErr  error
		expectedRole models.UserRole
	}{
		{
			roleUpdate: models.UserRole{
				ID:          userRoles[0].ID,
				Name:        "user role name",
				AccessRoles: accessRolesUpdate,
			},
		},
		{
			roleUpdate: models.UserRole{
				ID:          0,
				Name:        "user role name",
				AccessRoles: accessRolesUpdate,
			},
			expectedErr: models.ErrRequiredUserRoleID,
		},
		{
			roleUpdate: models.UserRole{
				ID:          999,
				Name:        "user role name",
				AccessRoles: accessRolesUpdate,
			},
			expectedErr: models.ErrUserRoleNotFound,
		},
		{
			roleUpdate: models.UserRole{
				ID:          userRoles[0].ID,
				Name:        "",
				AccessRoles: accessRolesUpdate,
			},
			expectedRole: models.UserRole{
				ID:          userRoles[0].ID,
				Name:        "user role name",
				AccessRoles: accessRolesUpdate,
			},
		},
		{
			roleUpdate: models.UserRole{
				ID:          userRoles[0].ID,
				Name:        "user role name",
				AccessRoles: []models.AccessRole{},
			},
		},
	}

	for _, testCase := range testCases {
		updatedRole, err := models.UpdateUserRole(testServer.Server.DB, testCase.roleUpdate)

		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			if testCase.expectedRole.Name == "" {
				checkUserRolesEqual(t, testCase.roleUpdate, updatedRole)
			} else {
				checkUserRolesEqual(t, testCase.expectedRole, updatedRole)
			}
		}
	}
}

func TestDeleteUserRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	userRole := testServer.Data.UserRoles[0]

	user := testServer.Data.Users[0]
	user, err = models.GetUserByID(testServer.Server.DB, user.ID)
	require.NoError(t, err)
	require.Len(t, user.UserRoles, 1)
	checkUserRolesInformationEqual(t, userRole, user.UserRoles[0])

	accessRole := testServer.Data.AccessRoles[0]
	accessRole, err = models.GetAccessRoleByID(testServer.Server.DB, accessRole.ID)
	require.NoError(t, err)
	require.Equal(t, userRole.ID, accessRole.UserRoleID)

	err = models.DeleteUserRole(testServer.Server.DB, userRole.ID)
	require.NoError(t, err)

	err = models.DeleteUserRole(testServer.Server.DB, userRole.ID)
	require.Equal(t, models.ErrUserRoleNotFound, err)

	user, err = models.GetUserByID(testServer.Server.DB, user.ID)
	require.NoError(t, err)
	require.Len(t, user.UserRoles, 0)

	_, err = models.GetAccessRoleByID(testServer.Server.DB, accessRole.ID)
	require.Equal(t, models.ErrAccessRoleNotFound, err)
}
