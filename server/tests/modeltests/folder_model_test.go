package modeltests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/models"
)

func checkFolderInformationEqual(t *testing.T, expectedFolder, actualFolder models.Folder) {
	assert.Equal(t, expectedFolder.Name, actualFolder.Name)
	assert.Equal(t, *expectedFolder.ParentFolderID, *actualFolder.ParentFolderID)
	assert.Equal(t, expectedFolder.LastEditorID, actualFolder.LastEditorID)
}

func checkFoldersEqual(t *testing.T, expectedFolder, actualFolder models.Folder) {
	checkFolderInformationEqual(t, expectedFolder, actualFolder)
	if assert.Len(t, actualFolder.AccessRoles, len(expectedFolder.AccessRoles)) {
		for i := range expectedFolder.AccessRoles {
			checkAccessRolesEqual(t, expectedFolder.AccessRoles[i], actualFolder.AccessRoles[i])
		}
	}
	if assert.Len(t, actualFolder.ChildFolders, len(expectedFolder.ChildFolders)) {
		for i := range expectedFolder.ChildFolders {
			checkFoldersEqual(t, expectedFolder.ChildFolders[i], actualFolder.ChildFolders[i])
		}
	}
	if assert.Len(t, actualFolder.Files, len(expectedFolder.Files)) {
		for i := range expectedFolder.Files {
			checkFilesEqual(t, expectedFolder.Files[i], actualFolder.Files[i])
		}
	}
}

func TestCreateFolder(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	folder := testServer.Data.Folders[0]
	invalidID := uint(999)

	testCases := []struct {
		folder      models.Folder
		expectedErr error
	}{
		{
			folder: models.Folder{
				Name:           "folder",
				LastEditorID:   user.ID,
				ParentFolderID: &folder.ID,
			},
			expectedErr: nil,
		},
		{
			folder: models.Folder{
				Name:           "",
				LastEditorID:   user.ID,
				ParentFolderID: &folder.ID,
			},
			expectedErr: models.ErrRequiredFolderName,
		},
		{
			folder: models.Folder{
				Name:           "folder",
				LastEditorID:   0,
				ParentFolderID: &folder.ID,
			},
			expectedErr: models.ErrRequiredLastEditorID,
		},
		{
			folder: models.Folder{
				Name:           "folder",
				LastEditorID:   user.ID,
				ParentFolderID: nil,
			},
			expectedErr: models.ErrRequiredParentFolderID,
		},
		{
			folder: models.Folder{
				Name:           "folder",
				LastEditorID:   user.ID,
				ParentFolderID: &invalidID,
			},
			expectedErr: models.ErrFolderNotFound,
		},
		{
			folder: models.Folder{
				Name:           "folder",
				LastEditorID:   user.ID,
				ParentFolderID: &folder.ID,
			},
			expectedErr: models.ErrFolderAlreadyExists,
		},
		{
			folder: models.Folder{
				Name:           "folder",
				LastEditorID:   testServer.Data.Users[1].ID,
				ParentFolderID: &folder.ID,
			},
			expectedErr: models.ErrFolderAlreadyExists,
		},
	}

	for _, testCase := range testCases {
		actualFolder, err := models.CreateFolder(testServer.Server.DB, testCase.folder)

		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			checkFoldersEqual(t, testCase.folder, actualFolder)
		}
	}
}

func TestGetFolderByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	for _, folder := range testServer.Data.Folders {
		foundFolder, err := models.GetFolderByID(testServer.Server.DB, folder.ID)
		require.NoError(t, err)
		checkFoldersEqual(t, folder, foundFolder)
	}

	_, err = models.GetFolderByID(testServer.Server.DB, 0)
	assert.Equal(t, models.ErrRequiredFolderID, err)

	_, err = models.GetFolderByID(testServer.Server.DB, 999)
	assert.Equal(t, models.ErrFolderNotFound, err)
}

func TestGetFolderByPath(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	for _, folder := range testServer.Data.Folders {
		foundFolder, err := models.GetFolderByPath(testServer.Server.DB, folder)
		require.NoError(t, err)
		checkFoldersEqual(t, folder, foundFolder)
	}

	_, err = models.GetFolderByPath(testServer.Server.DB, models.Folder{})
	assert.Equal(t, models.ErrRequiredFolderName, err)

	_, err = models.GetFolderByPath(testServer.Server.DB, models.Folder{
		Name: "not a folder",
	})
	assert.Equal(t, models.ErrFolderNotFound, err)
}

func TestUpdateFolder(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	folders := testServer.Data.Folders
	invalidID := uint(999)

	testCases := []struct {
		folderUpdate   models.Folder
		expectedErr    error
		expectedFolder models.Folder
	}{
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[1].ID,
				},
				Name:           "new folder name",
				ParentFolderID: &folders[0].ID,
				LastEditorID:   user.ID,
			},
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[1].ID,
				},
				Name:           "new folder name",
				ParentFolderID: &folders[0].ID,
				LastEditorID:   0,
			},
			expectedErr: models.ErrRequiredLastEditorID,
		},
		{
			folderUpdate: models.Folder{
				Name:           "new folder name",
				ParentFolderID: nil,
				LastEditorID:   user.ID,
			},
			expectedErr: models.ErrRequiredFolderID,
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: 999,
				},
				Name:           "new folder name",
				ParentFolderID: &folders[1].ID,
				LastEditorID:   user.ID,
			},
			expectedErr: models.ErrFolderNotFound,
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[1].ID,
				},
				Name:           "new folder name",
				ParentFolderID: &folders[1].ID,
				LastEditorID:   user.ID,
			},
			expectedErr: models.ErrInvalidParentFolderID,
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[1].ID,
				},
				LastEditorID: user.ID,
			},
			expectedErr: nil,
			expectedFolder: models.Folder{
				Model: models.Model{
					ID: folders[1].ID,
				},
				Name:           "new folder name",
				ParentFolderID: folders[1].ParentFolderID,
				LastEditorID:   user.ID,
			},
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[1].ID,
				},
				ParentFolderID: &invalidID,
				LastEditorID:   user.ID,
			},
			expectedErr: models.ErrFolderNotFound,
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[2].ID,
				},
				Name:           "new folder name",
				ParentFolderID: &folders[0].ID,
				LastEditorID:   user.ID,
			},
			expectedErr: models.ErrFolderAlreadyExists,
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[1].ID,
				},
				Name:           folders[1].Name,
				ParentFolderID: nil,
				LastEditorID:   user.ID,
			},
			expectedErr: nil,
			expectedFolder: models.Folder{
				Model: models.Model{
					ID: folders[1].ID,
				},
				Name:           folders[1].Name,
				ParentFolderID: folders[1].ParentFolderID,
				LastEditorID:   user.ID,
			},
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[2].ID,
				},
				Name:           "",
				ParentFolderID: &folders[0].ID,
				LastEditorID:   user.ID,
			},
			expectedErr: nil,
			expectedFolder: models.Folder{
				Model: models.Model{
					ID: folders[2].ID,
				},
				Name:           folders[2].Name,
				ParentFolderID: &folders[0].ID,
				LastEditorID:   user.ID,
			},
		},
		{
			folderUpdate: models.Folder{
				Model: models.Model{
					ID: folders[0].ID,
				},
				ParentFolderID: &folders[2].ID,
				LastEditorID:   user.ID,
			},
			expectedErr: models.ErrInvalidParentFolderID,
		},
	}

	for _, testCase := range testCases {
		updatedFolder, err := models.UpdateFolder(testServer.Server.DB, testCase.folderUpdate)
		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			if testCase.expectedFolder.Name != "" {
				checkFolderInformationEqual(t, testCase.expectedFolder, updatedFolder)
			} else {
				checkFolderInformationEqual(t, testCase.folderUpdate, updatedFolder)
			}
		}
	}
}

func TestDeleteFolder(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	folder := testServer.Data.Folders[2]
	file := testServer.Data.Files[2]
	accessRole := testServer.Data.AccessRoles[5]

	accessRole, err = models.GetAccessRoleByID(testServer.Server.DB, accessRole.ID)
	require.NoError(t, err)
	require.Equal(t, folder.ID, accessRole.FolderID)

	err = models.DeleteFolder(testServer.Server.DB, folder.ID, user.ID)
	require.Equal(t, models.ErrFolderNotEmpty, err)

	err = models.DeleteFile(testServer.Server.DB, file.ID, user.ID)
	require.NoError(t, err)

	err = models.DeleteFolder(testServer.Server.DB, folder.ID, user.ID)
	require.NoError(t, err)

	err = models.DeleteFolder(testServer.Server.DB, folder.ID, testServer.Data.Users[1].ID)
	require.Equal(t, models.ErrFolderNotFound, err)

	_, err = models.GetAccessRoleByID(testServer.Server.DB, accessRole.ID)
	require.Equal(t, models.ErrAccessRoleNotFound, err)
}

func TestGetUserAuthorizationFolder(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users
	folders := testServer.Data.Folders

	testCases := []struct {
		user                models.User
		folderID            uint
		expectedErr         error
		expectedAccessLevel models.AccessLevel
	}{
		{
			user:        users[0],
			folderID:    0,
			expectedErr: models.ErrRequiredFolderID,
		},
		{
			user:        users[0],
			folderID:    999,
			expectedErr: models.ErrFolderNotFound,
		},
		{
			user:                users[0],
			folderID:            folders[0].ID,
			expectedAccessLevel: models.Publisher,
		},
		{
			user:                users[1],
			folderID:            folders[0].ID,
			expectedAccessLevel: models.Unset,
		},
		{
			user:                users[1],
			folderID:            folders[1].ID,
			expectedAccessLevel: models.Uploader,
		},
		{
			user:                users[0],
			folderID:            folders[1].ID,
			expectedAccessLevel: models.Uploader,
		},
		{
			user:                users[3],
			folderID:            folders[0].ID,
			expectedAccessLevel: models.Unset,
		},
		{
			user:                users[3],
			folderID:            folders[2].ID,
			expectedAccessLevel: models.None,
		},
	}

	for _, testCase := range testCases {
		actualFolder, actualAccessLevel, err := models.GetUserAuthorizationFolder(testServer.Server.DB, testCase.user, testCase.folderID)
		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			assert.Equal(t, testCase.expectedAccessLevel, actualAccessLevel)
			checkFoldersEqual(t, folders[testCase.folderID-1], actualFolder)
		}
	}
}
