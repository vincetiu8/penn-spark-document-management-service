package models

import (
	"errors"

	"gorm.io/gorm"
)

// ErrRequiredFileID returned when no Model.ID is specified on a File.
var ErrRequiredFileID = errors.New("required file id")

// ErrRequiredFileName returned when no File.Name is specified.
var ErrRequiredFileName = errors.New("required file name")

// ErrFileNotFound returned when no File matches the given criteria.
var ErrFileNotFound = errors.New("file not found")

// ErrFileAlreadyExists returned when a File with the given information already exists.
var ErrFileAlreadyExists = errors.New("file already exists")

// File represents a single file's metadata in a Folder.
// Doesn't contain the actual file in memory - files can be accessed by querying the file system.
// A File object doesn't guarantee a file exists, as the file data needs to be uploaded after the creation of a File.
// File names are unique per FolderID.
type File struct {
	Model
	Name         string `gorm:"not null" json:"name"`
	FolderID     uint   `gorm:"not null" json:"folder_id"`
	LastEditorID uint   `gorm:"not null" json:"last_editor_id"`
	LastEditor   User   `gorm:"foreignKey:LastEditorID" json:"-"`
	IsPublished  bool   `json:"is_published"`
}

// prepare escapes File.Name before processing.
func (file *File) prepare() {
	file.Name = prepareString(file.Name)
}

// CreateFile creates a File.
func CreateFile(db *gorm.DB, file File) (File, error) {
	file.prepare()
	if file.LastEditorID == 0 {
		return File{}, ErrRequiredLastEditorID
	}

	_, err := getFolderByIDRaw(db, file.FolderID)
	if err != nil {
		return File{}, err
	}

	// Make sure no file exists with the same name in the same folder.
	_, err = GetFileByPath(db, file)
	if err == nil {
		return File{}, ErrFileAlreadyExists
	} else if err != ErrFileNotFound {
		return File{}, err
	}

	err = db.Create(&file).Take(&file).Error
	return file, err
}

// GetFileByID gets a File by its Model.ID.
func GetFileByID(db *gorm.DB, fileID uint) (File, error) {
	if fileID == 0 {
		return File{}, ErrRequiredFileID
	}
	file := File{
		Model: Model{ID: fileID},
	}

	err := db.Where(&file).Take(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return File{}, ErrFileNotFound
	}
	return file, err
}

// GetFileByPath gets a File by its File.Name and File.FolderID.
func GetFileByPath(db *gorm.DB, file File) (File, error) {
	file.prepare()
	if file.Name == "" {
		return File{}, ErrRequiredFileName
	}
	if file.FolderID == 0 {
		return File{}, ErrRequiredFolderID
	}

	file.ID = 0
	file.LastEditorID = 0
	err := db.Where(&file).Take(&file).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return File{}, ErrFileNotFound
	}
	return file, err
}

// UpdateFile updates a File based on its Model.ID.
func UpdateFile(db *gorm.DB, file File) (File, error) {
	file.prepare()
	if file.LastEditorID == 0 {
		return File{}, ErrRequiredLastEditorID
	}

	oldFile, err := GetFileByID(db, file.ID)
	if err != nil {
		return File{}, err
	}

	if file.Name == "" {
		file.Name = oldFile.Name
	}
	if file.FolderID == 0 {
		file.FolderID = oldFile.FolderID
	} else {
		_, err = getFolderByIDRaw(db, file.FolderID)
		if err != nil {
			return File{}, err
		}
	}

	// Make sure no file exists with the same name in the same folder.
	existingFile, err := GetFileByPath(db, file)
	if err == nil {
		if existingFile.ID != file.ID {
			return File{}, ErrFileAlreadyExists
		}
	} else if err != ErrFileNotFound {
		return File{}, err
	}

	err = db.Model(&file).Select("*").Updates(&file).Take(&file).Error
	return file, err
}

// DeleteFile deletes a file by its Model.ID.
func DeleteFile(db *gorm.DB, fileID, userID uint) error {
	file, err := GetFileByID(db, fileID)
	if err != nil {
		return err
	}

	file.LastEditorID = userID
	file, err = UpdateFile(db, file)
	return db.Where(&file).Delete(&file).Error
}

// GetUserAuthorizationFile gets a User's AccessLevel to a certain File.
func GetUserAuthorizationFile(db *gorm.DB, user User, fileID uint) (File, AccessLevel, error) {
	file, err := GetFileByID(db, fileID)
	if err != nil {
		return File{}, Unset, err
	}

	_, accessLevel, err := GetUserAuthorizationFolder(db, user, file.FolderID)
	return file, accessLevel, err
}
