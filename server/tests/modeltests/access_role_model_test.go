package modeltests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/models"
)

func checkAccessRolesEqual(t *testing.T, expectedRole, actualRole models.AccessRole) {
	assert.Equal(t, expectedRole.FolderID, actualRole.FolderID)
	assert.Equal(t, expectedRole.UserRoleID, actualRole.UserRoleID)
	assert.Equal(t, expectedRole.AccessLevel, actualRole.AccessLevel)
}

func TestCreateAccessRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	folders := testServer.Data.Folders
	userRole := testServer.Data.UserRoles[0]

	err = testServer.RefreshTable(&models.AccessRole{})
	require.NoError(t, err)

	testCases := []struct {
		role        models.AccessRole
		expectedErr error
	}{
		{
			role: models.AccessRole{
				FolderID:    folders[0].ID,
				UserRoleID:  userRole.ID,
				AccessLevel: models.Publisher,
			},
		},
		{
			role: models.AccessRole{
				FolderID:    0,
				UserRoleID:  userRole.ID,
				AccessLevel: models.Publisher,
			},
			expectedErr: models.ErrRequiredFolderID,
		},
		{
			role: models.AccessRole{
				FolderID:    999,
				UserRoleID:  userRole.ID,
				AccessLevel: models.Publisher,
			},
			expectedErr: models.ErrFolderNotFound,
		},
		{
			role: models.AccessRole{
				FolderID:    folders[0].ID,
				UserRoleID:  0,
				AccessLevel: models.Publisher,
			},
			expectedErr: models.ErrRequiredUserRoleID,
		},
		{
			role: models.AccessRole{
				FolderID:    folders[0].ID,
				UserRoleID:  999,
				AccessLevel: models.Publisher,
			},
			expectedErr: models.ErrUserRoleNotFound,
		},
		{
			role: models.AccessRole{
				FolderID:    folders[0].ID,
				UserRoleID:  userRole.ID,
				AccessLevel: models.Unset,
			},
			expectedErr: models.ErrRequiredAccessLevel,
		},
		{
			role: models.AccessRole{
				FolderID:    folders[0].ID,
				UserRoleID:  userRole.ID,
				AccessLevel: 999,
			},
			expectedErr: models.ErrInvalidAccessLevel,
		},
		{
			role: models.AccessRole{
				FolderID:    folders[0].ID,
				UserRoleID:  userRole.ID,
				AccessLevel: models.Publisher,
			},
			expectedErr: models.ErrAccessRoleAlreadyExists,
		},
		{
			role: models.AccessRole{
				FolderID:    folders[1].ID,
				UserRoleID:  userRole.ID,
				AccessLevel: models.Publisher,
			},
		},
	}

	for _, testCase := range testCases {
		createdRole, err := models.CreateAccessRole(testServer.Server.DB, testCase.role)

		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			checkAccessRolesEqual(t, testCase.role, createdRole)
		}
	}
}

func TestGetAccessRoleByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	for _, role := range testServer.Data.AccessRoles {
		foundRole, err := models.GetAccessRoleByID(testServer.Server.DB, role.ID)
		if assert.NoError(t, err) {
			checkAccessRolesEqual(t, role, foundRole)
		}
	}

	_, err = models.GetAccessRoleByID(testServer.Server.DB, 0)
	assert.Equal(t, models.ErrRequiredAccessRoleID, err)

	_, err = models.GetAccessRoleByID(testServer.Server.DB, 999)
	assert.Equal(t, models.ErrAccessRoleNotFound, err)
}

func TestUpdateAccessRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	folders := testServer.Data.Folders
	accessRoles := testServer.Data.AccessRoles
	userRoles := testServer.Data.UserRoles

	testCases := []struct {
		roleUpdate   models.AccessRole
		expectedErr  error
		expectedRole models.AccessRole
	}{
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    folders[1].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.Viewer,
			},
		},
		{
			roleUpdate: models.AccessRole{
				ID:          0,
				FolderID:    folders[1].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.Viewer,
			},
			expectedErr: models.ErrRequiredAccessRoleID,
		},
		{
			roleUpdate: models.AccessRole{
				ID:          999,
				FolderID:    folders[1].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.Viewer,
			},
			expectedErr: models.ErrAccessRoleNotFound,
		},
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    0,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.Uploader,
			},
			expectedRole: models.AccessRole{
				FolderID:    folders[1].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.Uploader,
			},
		},
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    999,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.Uploader,
			},
			expectedErr: models.ErrFolderNotFound,
		},
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    folders[0].ID,
				UserRoleID:  0,
				AccessLevel: models.Uploader,
			},
			expectedRole: models.AccessRole{
				FolderID:    folders[0].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.Uploader,
			},
		},
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    folders[0].ID,
				UserRoleID:  999,
				AccessLevel: models.Uploader,
			},
			expectedErr: models.ErrUserRoleNotFound,
		},
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    accessRoles[1].FolderID,
				UserRoleID:  userRoles[1].ID,
				AccessLevel: models.Uploader,
			},
			expectedErr: models.ErrAccessRoleAlreadyExists,
		},
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    folders[0].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: 0,
			},
			expectedRole: models.AccessRole{
				FolderID:    folders[0].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.Uploader,
			},
		},
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    folders[0].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: 999,
			},
			expectedErr: models.ErrInvalidAccessLevel,
		},
		{
			roleUpdate: models.AccessRole{
				ID:          accessRoles[5].ID,
				FolderID:    folders[0].ID,
				UserRoleID:  userRoles[2].ID,
				AccessLevel: models.None,
			},
		},
	}

	for _, testCase := range testCases {
		updatedRole, err := models.UpdateAccessRole(testServer.Server.DB, testCase.roleUpdate)
		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			if testCase.expectedRole.FolderID == 0 {
				checkAccessRolesEqual(t, testCase.roleUpdate, updatedRole)
			} else {
				checkAccessRolesEqual(t, testCase.expectedRole, updatedRole)
			}
		}
	}
}

func TestDeleteAccessRole(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	accessRole := testServer.Data.AccessRoles[5]

	userRole := testServer.Data.UserRoles[2]
	userRole, err = models.GetUserRoleByID(testServer.Server.DB, userRole.ID)
	require.NoError(t, err)
	require.Len(t, userRole.AccessRoles, 1)
	checkAccessRolesEqual(t, accessRole, userRole.AccessRoles[0])

	folder := testServer.Data.Folders[2]
	folder, err = models.GetFolderByID(testServer.Server.DB, folder.ID)
	require.NoError(t, err)
	require.Len(t, folder.AccessRoles, 4)
	checkAccessRolesEqual(t, accessRole, folder.AccessRoles[2])

	err = models.DeleteAccessRole(testServer.Server.DB, accessRole.ID)
	require.NoError(t, err)

	err = models.DeleteAccessRole(testServer.Server.DB, accessRole.ID)
	require.Equal(t, models.ErrAccessRoleNotFound, err)

	userRole, err = models.GetUserRoleByID(testServer.Server.DB, userRole.ID)
	require.NoError(t, err)
	require.Len(t, userRole.AccessRoles, 0)

	folder, err = models.GetFolderByID(testServer.Server.DB, folder.ID)
	require.NoError(t, err)
	require.Len(t, folder.AccessRoles, 3)
}
