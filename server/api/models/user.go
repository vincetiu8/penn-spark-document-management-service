package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ErrRequiredUserID returned when no Model.ID is specified on a User.
var ErrRequiredUserID = errors.New("required user id")

// ErrRequiredUsername returned when no User.Username is specified.
var ErrRequiredUsername = errors.New("required username")

// ErrUserAlreadyExists returned when a User with the given information already exists.
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrRequiredFirstName returned when no User.FirstName is specified.
var ErrRequiredFirstName = errors.New("required first name")

// ErrRequiredLastName returned when no User.LastName is specified.
var ErrRequiredLastName = errors.New("required last name")

// ErrRequiredPassword returned when no User.Password is specified.
var ErrRequiredPassword = errors.New("required password")

// ErrUserNotFound returned when no User matches the given criteria.
var ErrUserNotFound = errors.New("user not found")

// ErrIncorrectPassword returned when the given password doesn't match the User.Password.
var ErrIncorrectPassword = errors.New("incorrect password")

// User represents a User in the system.
// Each user has a unique Username and Model.ID.
type User struct {
	Model
	Username  string     `gorm:"not null;uniqueIndex" json:"username"`
	IsAdmin   bool       `json:"is_admin"`
	FirstName string     `gorm:"not null" json:"first_name"`
	LastName  string     `gorm:"not null" json:"last_name"`
	Password  string     `gorm:"not null" json:"password,omitempty"`
	UserRoles []UserRole `gorm:"many2many:assigned_user_roles" json:"user_roles"`
}

// hash hashes a password for storage.
func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// prepare escapes a User's fields before processing.
func (user *User) prepare() {
	user.Username = prepareString(user.Username)
	user.FirstName = prepareString(user.FirstName)
	user.LastName = prepareString(user.LastName)
}

// beforeSave hashes the User.Password before saving.
func (user *User) beforeSave() error {
	hashedPassword, err := hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// CreateUser creates a User.
func CreateUser(db *gorm.DB, user User) (User, error) {
	user.prepare()
	if user.Username == "" {
		return User{}, ErrRequiredUsername
	}
	if user.FirstName == "" {
		return User{}, ErrRequiredFirstName
	}
	if user.LastName == "" {
		return User{}, ErrRequiredLastName
	}
	if user.Password == "" {
		return User{}, ErrRequiredPassword
	}

	_, err := GetUserByUsername(db, user.Username)
	if err == nil {
		return User{}, ErrUserAlreadyExists
	} else if err != ErrUserNotFound {
		return User{}, err
	}

	err = user.beforeSave()
	if err != nil {
		return User{}, err
	}

	user.ID = 0
	err = db.Create(&user).Take(&user).Error
	if err != nil {
		return User{}, err
	}
	user.Password = ""
	if user.UserRoles == nil {
		user.UserRoles = []UserRole{}
	}
	return user, nil
}

// GetAllUsers returns a list of all present User.
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Password = ""
		err = db.Model(&users[i]).Association("UserRoles").Find(&users[i].UserRoles)
	}
	return users, nil
}

// getUserByIDRaw gets a User by its Model.ID.
func getUserByIDRaw(db *gorm.DB, uid uint) (User, error) {
	if uid == 0 {
		return User{}, ErrRequiredUserID
	}

	user := User{
		Model: Model{ID: uid},
	}

	err := db.Where(&user).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
	}

	err = db.Model(&user).Association("UserRoles").Find(&user.UserRoles)
	return user, err
}

// GetUserByID is a wrapper around getUserByIDRaw.
func GetUserByID(db *gorm.DB, uid uint) (User, error) {
	user, err := getUserByIDRaw(db, uid)
	if err != nil {
		return User{}, err
	}

	user.Password = ""
	return user, nil
}

// GetUserByUsername gets a User by its User.Username
func GetUserByUsername(db *gorm.DB, username string) (User, error) {
	user := User{Username: username}
	err := db.Where(&user).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}
	}
	return user, err
}

// UpdateUser updates a user based on its Model.ID.
func UpdateUser(db *gorm.DB, user User, isAdmin bool) (User, error) {
	user.prepare()

	oldUser, err := getUserByIDRaw(db, user.ID)
	if err != nil {
		if !isAdmin {
			return User{}, err
		}

		err = db.Unscoped().Where(&user).Take(&oldUser).Error
		if err == gorm.ErrRecordNotFound {
			err = ErrUserNotFound
		}
		if err != nil {
			return User{}, err
		}
		user.ID = oldUser.ID
		user.IsAdmin = oldUser.IsAdmin
	}

	if !isAdmin || user.Username == "" {
		user.Username = oldUser.Username
	} else {
		foundUser, err := GetUserByUsername(db, user.Username)
		if err == nil {
			if foundUser.ID != user.ID {
				return User{}, ErrUserAlreadyExists
			}
		} else if err != ErrUserNotFound {
			return User{}, err
		}
	}

	if !isAdmin || user.FirstName == "" {
		user.FirstName = oldUser.FirstName
	}

	if !isAdmin || user.LastName == "" {
		user.LastName = oldUser.LastName
	}

	if user.Password == "" {
		user.Password = oldUser.Password
	} else {
		err = user.beforeSave()
		if err != nil {
			return User{}, err
		}

		// Prevents user admin privilege from being changed with password updates
		user.IsAdmin = oldUser.IsAdmin
	}

	user.UserRoles = nil
	if !isAdmin {
		user.IsAdmin = oldUser.IsAdmin
		user.DeletedAt = oldUser.DeletedAt
	}

	err = db.Unscoped().Where("id = " + fmt.Sprint(user.ID)).Select("*").Updates(&user).Take(&user).Error
	if err != nil {
		return User{}, err
	}

	err = user.getUserRoles(db)
	if err != nil {
		return User{}, err
	}

	user.Password = ""
	if user.UserRoles == nil {
		user.UserRoles = []UserRole{}
	}
	return user, err
}

// DeleteUser deletes a User by its Model.ID.
func DeleteUser(db *gorm.DB, uid uint) error {
	user, err := getUserByIDRaw(db, uid)
	if err != nil {
		return err
	}

	err = db.Delete(&user).Error
	if err != nil {
		return err
	}

	return db.Model(&user).Association("UserRoles").Clear()
}

// getUserRoles gets a User's UserRoles.
func (user *User) getUserRoles(db *gorm.DB) error {
	return db.Model(&user).Association("UserRoles").Find(&user.UserRoles)
}

// checkUserRolePresent checks if a User has a UserRole.
func (user *User) checkUserRolePresent(db *gorm.DB, userRole UserRole) error {
	if userRole.ID == 0 {
		return ErrRequiredUserRoleID
	}

	err := user.getUserRoles(db)
	if err != nil {
		return err
	}

	for _, role := range user.UserRoles {
		if role.ID == userRole.ID {
			return nil
		}
	}

	return ErrUserRoleNotFound
}

// AddUserRole adds a UserRole to a User.
func AddUserRole(db *gorm.DB, uid uint, userRole UserRole) (User, error) {
	user, err := getUserByIDRaw(db, uid)
	if err != nil {
		return User{}, err
	}

	err = user.checkUserRolePresent(db, userRole)
	if err == nil {
		return User{}, ErrUserRoleAlreadyExists
	}
	if err != ErrUserRoleNotFound {
		return User{}, err
	}

	err = db.Model(&user).Association("UserRoles").Append(&userRole)
	if err != nil {
		return User{}, err
	}

	return user, user.getUserRoles(db)
}

// RemoveUserRole removes a UserRole from a User.
func RemoveUserRole(db *gorm.DB, uid uint, userRole UserRole) (User, error) {
	user, err := getUserByIDRaw(db, uid)
	if err != nil {
		return User{}, err
	}

	err = user.checkUserRolePresent(db, userRole)
	if err != nil {
		return User{}, err
	}

	err = db.Model(&user).Association("UserRoles").Delete(&userRole)
	if err != nil {
		return User{}, err
	}

	return user, user.getUserRoles(db)
}

// LoginUser validates a User login.
func LoginUser(db *gorm.DB, user User) (User, error) {
	user.prepare()
	if user.Username == "" {
		return User{}, ErrRequiredUsername
	}
	if user.Password == "" {
		return User{}, ErrRequiredPassword
	}

	matchingUser, err := GetUserByUsername(db, user.Username)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(matchingUser.Password), []byte(user.Password))
	if err != nil {
		return User{}, ErrIncorrectPassword
	}
	matchingUser.Password = ""
	return matchingUser, err
}
