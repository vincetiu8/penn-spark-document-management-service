package models

import (
	"errors"

	"gorm.io/gorm"
)

// ErrRequiredUserRoleName returned when no UserRole.Name is specified.
var ErrRequiredUserRoleName = errors.New("required user role name")

// ErrRequiredUserRoleID returned when no Model.ID is specified on a UserRole.
var ErrRequiredUserRoleID = errors.New("required user role id")

// ErrUserRoleNotFound returned when no UserRole matches the given criteria.
var ErrUserRoleNotFound = errors.New("user role not found")

// ErrUserRoleAlreadyExists returned when a UserRole with the given information already exists.
var ErrUserRoleAlreadyExists = errors.New("user role already exists")

// UserRole represents a role assignable to a User.
// It contains a set of AccessRoles giving the User access to each of the Folder specified.
// Different UserRole can be stacked with the User inheriting the most powerful permissions of each role.
type UserRole struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"not null;uniqueIndex" json:"name"`
	AccessRoles []AccessRole `json:"access_roles"`
}

// prepare escapes UserRole.Name before processing.
func (role *UserRole) prepare() {
	role.Name = prepareString(role.Name)
}

// CreateUserRole creates a UserRole.
func CreateUserRole(db *gorm.DB, userRole UserRole) (UserRole, error) {
	userRole.prepare()

	if userRole.Name == "" {
		return UserRole{}, ErrRequiredUserRoleName
	}
	_, err := getUserRoleByName(db, userRole.Name)
	if err == nil {
		return UserRole{}, ErrUserRoleAlreadyExists
	} else if err != ErrUserRoleNotFound {
		return UserRole{}, err
	}

	err = db.Create(&userRole).Take(&userRole).Error
	if err != nil {
		return UserRole{}, err
	}

	if userRole.AccessRoles == nil {
		userRole.AccessRoles = []AccessRole{}
	}
	return userRole, nil
}

// GetAllUserRoles returns a list of all present UserRole.
func GetAllUserRoles(db *gorm.DB) ([]UserRole, error) {
	var userRoles []UserRole
	err := db.Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	for i := range userRoles {
		err = db.Model(&userRoles[i]).Association("AccessRoles").Find(&userRoles[i].AccessRoles)
	}
	return userRoles, nil
}

// GetUserRoleByID gets a UserRole by its Model.ID.
func GetUserRoleByID(db *gorm.DB, roleID uint) (UserRole, error) {
	if roleID == 0 {
		return UserRole{}, ErrRequiredUserRoleID
	}

	userRole := UserRole{
		ID: roleID,
	}

	err := db.Where(&userRole).Take(&userRole).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserRole{}, ErrUserRoleNotFound
	} else if err != nil {
		return UserRole{}, err
	}

	err = db.Model(&userRole).Association("AccessRoles").Find(&userRole.AccessRoles)
	if err != nil {
		return UserRole{}, err
	}

	return userRole, nil
}

// getUserRoleByName gets a UserRole by its UserRole.Name.
func getUserRoleByName(db *gorm.DB, name string) (UserRole, error) {
	role := UserRole{
		Name: name,
	}
	err := db.Where(&role).Take(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserRole{}, ErrUserRoleNotFound
	}
	return role, err
}

// UpdateUserRole updates a UserRole based on its Model.ID.
func UpdateUserRole(db *gorm.DB, userRole UserRole) (UserRole, error) {
	userRole.prepare()

	oldRole, err := GetUserRoleByID(db, userRole.ID)
	if err != nil {
		return UserRole{}, err
	}

	if userRole.Name == "" {
		userRole.Name = oldRole.Name
	}

	err = db.Model(&userRole).Updates(&userRole).Take(&userRole).Error
	if err != nil {
		return UserRole{}, err
	}

	err = db.Model(&userRole).Association("AccessRoles").Find(&userRole.AccessRoles)
	return userRole, err
}

// DeleteUserRole deletes a UserRole by its Model.ID.
func DeleteUserRole(db *gorm.DB, roleID uint) error {
	userRole, err := GetUserRoleByID(db, roleID)
	if err != nil {
		return err
	}

	return db.Select("AccessRoles").Delete(&userRole).Error
}
