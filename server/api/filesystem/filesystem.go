// Package filesystem provides a wrapper around the afero filesystem.
package filesystem

import (
	"io"
	"strconv"

	"github.com/spf13/afero"
)

// FileSystem struct provides a wrapper around the afero filesystem.
type FileSystem struct {
	afero.Fs
	FilePath string
}

// idToFilePath converts a file's id to its local path.
func (fs *FileSystem) idToFilePath(id uint) string {
	return fs.FilePath + "/" + strconv.Itoa(int(id))
}

// UpsertFileRaw upserts a file's data in the local filesystem by its id.
func (fs *FileSystem) UpsertFileRaw(id uint, reader io.Reader) error {
	return afero.WriteReader(fs, fs.idToFilePath(id), reader)
}

// GetFileRaw returns a file's data by its id.
func (fs *FileSystem) GetFileRaw(id uint) (io.ReadSeeker, error) {
	return fs.Open(fs.idToFilePath(id))
}
