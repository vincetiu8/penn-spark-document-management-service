package models

import (
	"errors"

	"gorm.io/gorm"
)

// ErrRequiredAccessLevel returned when no AccessRole.AccessLevel is specified.
var ErrRequiredAccessLevel = errors.New("required access level")

// ErrInvalidAccessLevel returned when an invalid AccessLevel is specified.
var ErrInvalidAccessLevel = errors.New("invalid access level")

// ErrRequiredAccessRoleID returned when no AccessRole.ID is specified.
var ErrRequiredAccessRoleID = errors.New("required access role id")

// ErrAccessRoleNotFound returned when no AccessRole matches the required criteria.
var ErrAccessRoleNotFound = errors.New("access role not found")

// ErrAccessRoleAlreadyExists returned when an AccessRole with the given information already exists.
var ErrAccessRoleAlreadyExists = errors.New("access role already exists")

// AccessRole represents a single AccessLevel to a single Folder assigned to a UserRole.
// Only 1 AccessRole can exist per UserRole per Folder.
type AccessRole struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	FolderID    uint        `gorm:"not null" json:"folder_id"`
	UserRoleID  uint        `gorm:"not null" json:"user_role_id"`
	AccessLevel AccessLevel `gorm:"not null" json:"access_level"`
}

// AccessLevel represents the AccessLevel of a User.
type AccessLevel uint

const (
	// Unset represents an unset AccessLevel.
	// Exists to prevent users from accidentally not entering an AccessLevel.
	Unset AccessLevel = iota

	// None represents a lack of access to a Folder.
	// Currently obsolete as by default, users have no access to a child Folder even if they have access to the parent.
	// Left in place in case future upgrades make it necessary to have a None access level.
	None

	// Viewer represents a user that can only view a File inside a Folder, but cannot edit it.
	// In order to access a child Folder, a user must have at least the Viewer role in the parent Folder.
	Viewer

	// Uploader represents a user that has all the permissions of Viewer but can also upload a draft File into a folder.
	// Uploader can view a draft File they own but not others.
	// Uploader cannot edit their File's metadata once uploaded.
	Uploader

	// Publisher represents a user that has all the permissions of Uploader but can also publish a draft File.
	// A Publisher can also edit a File's metadata and delete a File.
	// A Publisher can also create a child Folder in a parent Folder, but by default will not have any rights in them.
	// To interact with created child folders, a Publisher must contact an admin to assign an appropriate AccessRole.
	Publisher
)

// CreateAccessRole creates an AccessRole.
func CreateAccessRole(db *gorm.DB, accessRole AccessRole) (AccessRole, error) {
	// Validate access role information
	if accessRole.FolderID == 0 {
		return AccessRole{}, ErrRequiredFolderID
	}
	if accessRole.UserRoleID == 0 {
		return AccessRole{}, ErrRequiredUserRoleID
	}
	if accessRole.AccessLevel == Unset {
		return AccessRole{}, ErrRequiredAccessLevel
	}
	if accessRole.AccessLevel > Publisher {
		return AccessRole{}, ErrInvalidAccessLevel
	}

	// Check folder referenced in access role exists
	_, err := getFolderByIDRaw(db, accessRole.FolderID)
	if err != nil {
		return AccessRole{}, err
	}

	// Check if referenced user role exists.
	userRole, err := GetUserRoleByID(db, accessRole.UserRoleID)
	if err != nil {
		return AccessRole{}, err
	}

	// Assert that no access role with the same folder id exists in the user role.
	for _, role := range userRole.AccessRoles {
		if role.FolderID == accessRole.FolderID {
			return AccessRole{}, ErrAccessRoleAlreadyExists
		}
	}

	return accessRole, db.Create(&accessRole).Take(&accessRole).Error
}

// GetAccessRoleByID returns an AccessRole by its AccessRole.ID.
func GetAccessRoleByID(db *gorm.DB, roleID uint) (AccessRole, error) {
	if roleID == 0 {
		return AccessRole{}, ErrRequiredAccessRoleID
	}

	accessRole := AccessRole{
		ID: roleID,
	}

	err := db.Where(&accessRole).Take(&accessRole).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return AccessRole{}, ErrAccessRoleNotFound
	}
	return accessRole, err
}

// UpdateAccessRole updates an AccessRole by its AccessRole.ID.
func UpdateAccessRole(db *gorm.DB, accessRole AccessRole) (AccessRole, error) {
	if accessRole.AccessLevel > Publisher {
		return AccessRole{}, ErrInvalidAccessLevel
	}

	oldRole, err := GetAccessRoleByID(db, accessRole.ID)
	if err != nil {
		return AccessRole{}, err
	}

	if accessRole.FolderID == 0 {
		accessRole.FolderID = oldRole.FolderID
	} else {
		_, err = getFolderByIDRaw(db, accessRole.FolderID)
		if err != nil {
			return AccessRole{}, err
		}
	}
	if accessRole.UserRoleID == 0 {
		accessRole.UserRoleID = oldRole.UserRoleID
	}

	// Check if referenced user role exists.
	userRole, err := GetUserRoleByID(db, accessRole.UserRoleID)
	if err != nil {
		return AccessRole{}, err
	}

	// Assert that no access role with the same folder id exists in the user role if the folder id was changed.
	if accessRole.FolderID != oldRole.FolderID {
		for _, role := range userRole.AccessRoles {
			if role.FolderID == accessRole.FolderID {
				return AccessRole{}, ErrAccessRoleAlreadyExists
			}
		}
	}

	if accessRole.AccessLevel == Unset {
		accessRole.AccessLevel = oldRole.AccessLevel
	}

	err = db.Model(&accessRole).Updates(&accessRole).Take(&accessRole).Error
	return accessRole, err
}

// DeleteAccessRole deletes an AccessRole by its AccessRole.ID.
func DeleteAccessRole(db *gorm.DB, roleID uint) error {
	accessRole, err := GetAccessRoleByID(db, roleID)
	if err != nil {
		return err
	}

	return db.Where(&accessRole).Delete(&accessRole).Error
}
