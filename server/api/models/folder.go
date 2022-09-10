package models

import (
	"errors"

	"gorm.io/gorm"
)

// ErrRequiredFolderID returned when no Model.ID is specified on a Folder
var ErrRequiredFolderID = errors.New("required folder id")

// ErrFolderNotFound returned when no Folder matches the given criteria.
var ErrFolderNotFound = errors.New("folder not found")

// ErrFolderAlreadyExists returned when a Folder with the given information already exists.
var ErrFolderAlreadyExists = errors.New("folder already exists")

// ErrRequiredFolderName returned when no Folder.Name is specified.
var ErrRequiredFolderName = errors.New("required folder name")

// ErrRequiredLastEditorID returned when no LastEditorID is specified.
var ErrRequiredLastEditorID = errors.New("required last editor id")

// ErrFolderNotEmpty returned when a Folder contains Folder.Files.
var ErrFolderNotEmpty = errors.New("cannot delete non-empty folder")

// ErrRequiredParentFolderID returned when no Folder.ParentFolderID is specified.
var ErrRequiredParentFolderID = errors.New("required parent folder id")

// ErrInvalidParentFolderID returned when an invalid Folder.ParentFolderID is specified.
var ErrInvalidParentFolderID = errors.New("invalid parent folder id")

// Folder represents a single Folder in the document system.
// Is merely a representation - a File is not stored by its true path when saved.
// The folder path is reconstructed by the browser client.
// Folder names are unique per ParentFolderID.
type Folder struct {
	Model
	Name           string       `gorm:"not null" json:"name"`
	ParentFolderID *uint        `gorm:"not_null" json:"parent_folder_id"`
	ChildFolders   []Folder     `gorm:"foreignKey:ParentFolderID" json:"child_folders"`
	Files          []File       `gorm:"foreignKey:FolderID" json:"files"`
	LastEditorID   uint         `gorm:"not null" json:"last_editor_id"`
	LastEditor     User         `gorm:"foreignKey:LastEditorID" json:"last_editor"`
	AccessRoles    []AccessRole `gorm:"foreignKey:FolderID" json:"access_roles"`
}

// prepare escapes Folder.Name before processing.
func (folder *Folder) prepare() {
	folder.Name = prepareString(folder.Name)
}

// formatFolderContents properly formats the return values of a gorm query to standardize empty arrays.
func (folder *Folder) formatFolderContents() {
	if folder.ChildFolders == nil {
		folder.ChildFolders = []Folder{}
	}
	if folder.Files == nil {
		folder.Files = []File{}
	}
}

// CreateFolder creates a Folder.
func CreateFolder(db *gorm.DB, folder Folder) (Folder, error) {
	folder.prepare()
	if folder.LastEditorID == 0 {
		return Folder{}, ErrRequiredLastEditorID
	}
	if folder.ParentFolderID == nil {
		return Folder{}, ErrRequiredParentFolderID
	}
	if *folder.ParentFolderID == 0 {
		return Folder{}, ErrRequiredParentFolderID
	}

	// Make sure no folder exists with the same name in the parent folder.
	_, err := GetFolderByPath(db, folder)
	if err == nil {
		return Folder{}, ErrFolderAlreadyExists
	} else if err != ErrFolderNotFound {
		return Folder{}, err
	}

	_, err = getFolderByIDRaw(db, *folder.ParentFolderID)
	if err != nil {
		return Folder{}, err
	}

	folder.ID = 0
	err = db.Create(&folder).Take(&folder).Error
	folder.formatFolderContents()
	return folder, err
}

// getFolderByIDRaw gets a Folder by its Model.ID.
func getFolderByIDRaw(db *gorm.DB, folderID uint) (Folder, error) {
	folder := Folder{
		Model: Model{ID: folderID},
	}

	err := db.Where(&folder).Take(&folder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Folder{}, ErrFolderNotFound
	}

	if folder.ParentFolderID == nil {
		parentID := uint(0)
		folder.ParentFolderID = &parentID
	}

	return folder, err
}

// GetFolderByID adds a wrapper around getFolderByIDRaw.
func GetFolderByID(db *gorm.DB, folderID uint) (Folder, error) {
	if folderID == 0 {
		return Folder{}, ErrRequiredFolderID
	}

	folder, err := getFolderByIDRaw(db, folderID)
	if err != nil {
		return Folder{}, err
	}

	err = db.Model(&folder).Association("AccessRoles").Find(&folder.AccessRoles)
	if err != nil {
		return Folder{}, err
	}

	err = db.Model(&folder).Association("ChildFolders").Find(&folder.ChildFolders)
	if err != nil {
		return Folder{}, err
	}

	err = db.Model(&folder).Association("Files").Find(&folder.Files)
	if err != nil {
		return Folder{}, err
	}

	return folder, err
}

// GetFolderByPath gets a Folder by its Folder.Name and Folder.ParentFolderID.
func GetFolderByPath(db *gorm.DB, folder Folder) (Folder, error) {
	folder.prepare()
	if folder.Name == "" {
		return Folder{}, ErrRequiredFolderName
	}

	folder.ID = 0
	folder.LastEditorID = 0
	err := db.Where(&folder).Take(&folder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Folder{}, ErrFolderNotFound
	}
	return folder, err
}

// UpdateFolder updates a Folder based on its Model.ID
func UpdateFolder(db *gorm.DB, folder Folder) (Folder, error) {
	folder.prepare()
	if folder.ID == 0 {
		return Folder{}, ErrRequiredFolderID
	}
	if folder.LastEditorID == 0 {
		return Folder{}, ErrRequiredLastEditorID
	}

	oldFolder, err := getFolderByIDRaw(db, folder.ID)
	if err != nil {
		return Folder{}, err
	}

	// Makes sure there are no circular folder references.
	if folder.ParentFolderID == nil || *folder.ParentFolderID == 0 {
		folder.ParentFolderID = oldFolder.ParentFolderID
	} else {
		folderID := folder.ParentFolderID
		for {
			if *folderID == 0 {
				break
			} else if *folderID == folder.ID {
				// Circular reference found.
				return Folder{}, ErrInvalidParentFolderID
			}

			parentFolder, err := getFolderByIDRaw(db, *folderID)
			if err != nil {
				return Folder{}, err
			}

			folderID = parentFolder.ParentFolderID
		}
	}
	if folder.Name == "" {
		folder.Name = oldFolder.Name
	}

	// Check that no existing Folder with the same name in the parent folder.
	existingFolder, err := GetFolderByPath(db, folder)
	if err == nil {
		if existingFolder.ID != folder.ID {
			return Folder{}, ErrFolderAlreadyExists
		}
	} else if err != ErrFolderNotFound {
		return Folder{}, err
	}

	err = db.Model(&folder).Updates(&folder).Take(&folder).Error
	folder.formatFolderContents()
	return folder, err
}

// DeleteFolder deletes a Folder by its Model.ID.
func DeleteFolder(db *gorm.DB, folderID, userID uint) error {
	folder, err := GetFolderByID(db, folderID)
	if err != nil {
		return err
	}
	if len(folder.ChildFolders) > 0 || len(folder.Files) > 0 {
		return ErrFolderNotEmpty
	}

	folder.LastEditorID = userID
	folder, err = UpdateFolder(db, folder)
	return db.Select("AccessRoles").Delete(&folder).Error
}

// GetUserAuthorizationFolder gets a User's AccessLevel in a Folder.
func GetUserAuthorizationFolder(db *gorm.DB, user User, folderID uint) (Folder, AccessLevel, error) {
	// Get the folder
	folder, err := GetFolderByID(db, folderID)
	if err != nil {
		return Folder{}, Unset, err
	}

	// Get all the user's access roles from the database
	var accessRoles []AccessRole
	err = db.Model(&user.UserRoles).Association("AccessRoles").Find(&accessRoles)
	if err != nil {
		return Folder{}, Unset, err
	}

	// Default access level is none
	accessLevel := None
	// Admins need viewer permissions to all folders by default
	// This allows them to grant roles to other users
	if user.IsAdmin {
		accessLevel = Viewer
	}

	// Loop through all the user's access roles
	// The user can have multiple user roles with different access roles in the same folder
	// Therefore, we take the one granting the most privileges
	for _, accessRole := range accessRoles {
		if accessRole.FolderID == folderID && accessRole.AccessLevel > accessLevel {
			accessLevel = accessRole.AccessLevel

			// No point continuing loop if highest access role is achieved
			if accessLevel == Publisher {
				break
			}
		}
	}

	return folder, accessLevel, nil
}
