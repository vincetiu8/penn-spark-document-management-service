package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/invincibot/penn-spark-server/api/models"
)

func AccessRoleToJSON(accessRole models.AccessRole) string {
	params := []string{
		WrapUint("folder_id", accessRole.FolderID),
		WrapUint("user_role_id", accessRole.UserRoleID),
		WrapUint("access_level", uint(accessRole.AccessLevel)),
	}

	return ParamsToJSON(params)
}

func CheckAccessRolesEqual(t *testing.T, expectedAccessRole models.AccessRole, response map[string]interface{}) {
	assert.Equal(t, float64(expectedAccessRole.FolderID), response["folder_id"])
	assert.Equal(t, float64(expectedAccessRole.UserRoleID), response["user_role_id"])
	assert.Equal(t, float64(expectedAccessRole.AccessLevel), response["access_level"])
}

func FileToJSON(file models.File) string {
	params := []string{
		WrapString("name", file.Name),
		WrapUint("folder_id", file.FolderID),
	}

	return ParamsToJSON(params)
}

func CheckFilesEqual(t *testing.T, expectedFile models.File, response map[string]interface{}) {
	assert.Equal(t, expectedFile.Name, response["name"])
	assert.Equal(t, float64(expectedFile.FolderID), response["folder_id"])
	assert.Equal(t, float64(expectedFile.LastEditorID), response["last_editor_id"])
	assert.Nil(t, response["deleted_at"])
	assert.Equal(t, expectedFile.IsPublished, response["is_published"])
}

func FolderToJSON(folder models.Folder) string {
	if folder.ParentFolderID == nil {
		parentFolderID := uint(0)
		folder.ParentFolderID = &parentFolderID
	}

	params := []string{
		WrapString("name", folder.Name),
		WrapUint("parent_folder_id", *folder.ParentFolderID),
	}
	return ParamsToJSON(params)
}

func CheckFoldersEqual(t *testing.T, expectedFolder models.Folder, response map[string]interface{}) {
	assert.Equal(t, expectedFolder.Name, response["name"])
	assert.Equal(t, float64(*expectedFolder.ParentFolderID), response["parent_folder_id"])
	assert.Equal(t, float64(expectedFolder.LastEditorID), response["last_editor_id"])
	assert.Nil(t, response["deleted_at"])
}

func UserToJSON(user models.User, mode string) string {
	var params []string
	switch mode {
	case "login":
		params = append(params,
			WrapString("username", user.Username),
			WrapString("password", user.Password),
		)
	default:
		params = append(params,
			WrapString("username", user.Username),
			WrapString("first_name", user.FirstName),
			WrapString("last_name", user.LastName),
			WrapString("password", user.Password),
			WrapUserRoles(user.UserRoles),
			WrapBool("is_admin", user.IsAdmin),
		)
	}

	return ParamsToJSON(params)
}

func CheckUserInformationEqual(t *testing.T, expectedUser models.User, response map[string]interface{}) {
	assert.Equal(t, expectedUser.Username, response["username"])
	assert.Equal(t, expectedUser.FirstName, response["first_name"])
	assert.Equal(t, expectedUser.LastName, response["last_name"])
	_, ok := response["password"]
	assert.False(t, ok)
	assert.Equal(t, expectedUser.IsAdmin, response["is_admin"])
	assert.Nil(t, response["deleted_at"])
}

func CheckUsersEqual(t *testing.T, expectedUser models.User, response map[string]interface{}) {
	CheckUserInformationEqual(t, expectedUser, response)
	if expectedUser.UserRoles != nil {
		userRoles, ok := response["user_roles"].([]interface{})
		if assert.True(t, ok) {
			for i := range userRoles {
				CheckUserRolesEqual(t, expectedUser.UserRoles[i], userRoles[i].(map[string]interface{}))
			}
		}
	}
}

func UserRoleToJSON(userRole models.UserRole) string {
	params := []string{
		WrapString("name", userRole.Name),
	}

	return ParamsToJSON(params)
}

func CheckUserRolesEqual(t *testing.T, expectedUserRole models.UserRole, response map[string]interface{}) {
	assert.Equal(t, expectedUserRole.Name, response["name"])
}
