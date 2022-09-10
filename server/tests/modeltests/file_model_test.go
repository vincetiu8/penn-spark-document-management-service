package modeltests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/models"
)

func checkFilesEqual(t *testing.T, expectedFile, actualFile models.File) {
	assert.Equal(t, expectedFile.Name, actualFile.Name)
	assert.Equal(t, expectedFile.FolderID, actualFile.FolderID)
	assert.Equal(t, expectedFile.LastEditorID, actualFile.LastEditorID)
	assert.Equal(t, expectedFile.IsPublished, actualFile.IsPublished)
}

func TestCreateFile(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	folder := testServer.Data.Folders[0]

	err = testServer.RefreshTable(&models.File{})
	require.NoError(t, err)

	testCases := []struct {
		file        models.File
		expectedErr error
	}{
		{
			file: models.File{
				Name:         "file",
				FolderID:     folder.ID,
				LastEditorID: user.ID,
			},
			expectedErr: nil,
		},
		{
			file: models.File{
				Name:         "",
				FolderID:     folder.ID,
				LastEditorID: user.ID,
			},
			expectedErr: models.ErrRequiredFileName,
		},
		{
			file: models.File{
				Name:         "file",
				FolderID:     0,
				LastEditorID: user.ID,
			},
			expectedErr: models.ErrRequiredFolderID,
		},
		{
			file: models.File{
				Name:         "file",
				FolderID:     999,
				LastEditorID: user.ID,
			},
			expectedErr: models.ErrFolderNotFound,
		},
		{
			file: models.File{
				Name:         "file",
				FolderID:     folder.ID,
				LastEditorID: 0,
			},
			expectedErr: models.ErrRequiredLastEditorID,
		},
		{
			file: models.File{
				Name:         "file",
				FolderID:     folder.ID,
				LastEditorID: user.ID,
			},
			expectedErr: models.ErrFileAlreadyExists,
		},
		{
			file: models.File{
				Name:         "file",
				FolderID:     testServer.Data.Folders[1].ID,
				LastEditorID: user.ID,
			},
			expectedErr: nil,
		},
		{
			file: models.File{
				Name:         "different file",
				IsPublished:  true,
				FolderID:     folder.ID,
				LastEditorID: user.ID,
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		actualFile, err := models.CreateFile(testServer.Server.DB, testCase.file)

		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			checkFilesEqual(t, testCase.file, actualFile)
		}
	}
}

func TestGetFileByID(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	for _, file := range testServer.Data.Files {
		foundFile, err := models.GetFileByID(testServer.Server.DB, file.ID)
		assert.NoError(t, err)
		checkFilesEqual(t, file, foundFile)
	}

	_, err = models.GetFileByID(testServer.Server.DB, 0)
	assert.Equal(t, models.ErrRequiredFileID, err)

	_, err = models.GetFileByID(testServer.Server.DB, 999)
	assert.Equal(t, models.ErrFileNotFound, err)
}

func TestGetFileByPath(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	for _, file := range testServer.Data.Files {
		foundFile, err := models.GetFileByPath(testServer.Server.DB, file)
		if assert.NoError(t, err) {
			checkFilesEqual(t, file, foundFile)
		}
	}

	_, err = models.GetFileByPath(testServer.Server.DB, models.File{})
	assert.Equal(t, models.ErrRequiredFileName, err)

	_, err = models.GetFileByPath(testServer.Server.DB, models.File{Name: "file"})
	assert.Equal(t, models.ErrRequiredFolderID, err)

	_, err = models.GetFileByPath(testServer.Server.DB, models.File{Name: "file", FolderID: 999})
	assert.Equal(t, models.ErrFileNotFound, err)
}

func TestUpdateFile(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	folders := testServer.Data.Folders
	files := testServer.Data.Files

	testCases := []struct {
		fileUpdate   models.File
		expectedErr  error
		expectedFile models.File
	}{
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				Name:         "new file name",
				FolderID:     folders[0].ID,
				LastEditorID: user.ID,
			},
			expectedErr: nil,
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: 0,
				},
				Name:         "new file name",
				FolderID:     folders[0].ID,
				LastEditorID: user.ID,
			},
			expectedErr: models.ErrRequiredFileID,
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: 999,
				},
				Name:         "new file name",
				FolderID:     folders[0].ID,
				LastEditorID: user.ID,
			},
			expectedErr: models.ErrFileNotFound,
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				Name:         "new file name",
				FolderID:     folders[0].ID,
				LastEditorID: 0,
			},
			expectedErr: models.ErrRequiredLastEditorID,
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				LastEditorID: user.ID,
			},
			expectedErr: nil,
			expectedFile: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				Name:         "new file name",
				FolderID:     files[0].FolderID,
				LastEditorID: user.ID,
			},
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				FolderID:     999,
				LastEditorID: user.ID,
			},
			expectedErr: models.ErrFolderNotFound,
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				Name:         files[1].Name,
				LastEditorID: user.ID,
				FolderID:     files[1].FolderID,
			},
			expectedErr: models.ErrFileAlreadyExists,
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				FolderID:     folders[1].ID,
				LastEditorID: user.ID,
			},
			expectedErr: nil,
			expectedFile: models.File{
				Name:         "new file name",
				FolderID:     folders[1].ID,
				LastEditorID: user.ID,
			},
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				Name:         "file name",
				LastEditorID: user.ID,
			},
			expectedErr: nil,
			expectedFile: models.File{
				Name:         "file name",
				FolderID:     folders[1].ID,
				LastEditorID: user.ID,
			},
		},
		{
			fileUpdate: models.File{
				Model: models.Model{
					ID: files[0].ID,
				},
				Name:         "file name",
				FolderID:     folders[1].ID,
				IsPublished:  true,
				LastEditorID: user.ID,
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		updatedFile, err := models.UpdateFile(testServer.Server.DB, testCase.fileUpdate)
		if assert.Equal(t, testCase.expectedErr, err) && testCase.expectedErr == nil {
			if testCase.expectedFile.Name == "" {
				checkFilesEqual(t, testCase.fileUpdate, updatedFile)
			} else {
				checkFilesEqual(t, testCase.expectedFile, updatedFile)
			}

		}
	}
}

func TestDeleteFile(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	file := testServer.Data.Files[0]

	err = models.DeleteFile(testServer.Server.DB, file.ID, user.ID)
	require.NoError(t, err)

	err = models.DeleteFile(testServer.Server.DB, file.ID, user.ID)
	require.Equal(t, models.ErrFileNotFound, err)
}

func TestGetUserAuthorizationFile(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	user := testServer.Data.Users[0]
	file := testServer.Data.Files[0]

	_, _, err = models.GetUserAuthorizationFile(testServer.Server.DB, user, 0)
	assert.Equal(t, models.ErrRequiredFileID, err)

	foundFile, accessLevel, err := models.GetUserAuthorizationFile(testServer.Server.DB, user, file.ID)
	if assert.NoError(t, err) {
		assert.Equal(t, models.Publisher, accessLevel)
		checkFilesEqual(t, file, foundFile)
	}
}
