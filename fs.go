package vfs

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// A File is returned by a FileSystem's Open method.
type File interface {
	io.Closer
	io.Reader
	io.Writer
	Readdir(count int) ([]os.FileInfo, error)
	Seek(offset int64, whence int) (int64, error)
	Stat() (os.FileInfo, error)
}

// A fileSystem implements access to files below a root directory.
type fileSystem struct {
	root string
}

// FileSystem returns a new file system that is rooted in dir.
// A file in a file system always uses the slash as path separator,
// regardless of the operating system.
func FileSystem(dir string) *fileSystem {
	return &fileSystem{root: dir}
}

// Resolve returns an absolute file name inside or below the
// root of the FileSystem.
func (fs *fileSystem) Resolve(name string) string {
	fname := path.Clean("/" + filepath.ToSlash(name))
	return filepath.Join(fs.root, filepath.FromSlash(fname))
}

// Join joins any number of path elements into a single file path, adding
// slashes if necessary. The result is cleaned, empty strings are ignored.
// The result is always an absolute file name inside or below the
// root of the FileSystem.
func (fs *fileSystem) Join(elem ...string) string {
	return fs.Resolve(path.Join(elem...))
}

// ServeFile serves the file for HTTP.
func (fs *fileSystem) ServeFile(w http.ResponseWriter, r *http.Request, name string) {
	http.ServeFile(w, r, fs.Resolve(name))
}

// Mkdir creates a directory below the FileSystem root.
func (fs *fileSystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(fs.Resolve(name), perm)
}

// MkdirAll creates all the directories below the FileSystem root.
func (fs *fileSystem) MkdirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(fs.Resolve(name), perm)
}

// Open a file in the file system.
func (fs *fileSystem) Open(name string) (File, error) {
	f, err := os.Open(fs.Resolve(name))
	if err != nil {
		return nil, err
	}
	return f, nil
}

// OpenFile opens a file with the specified permissions.
func (fs *fileSystem) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(fs.Resolve(name), flag, perm)
}

// Create creates the file.
func (fs *fileSystem) Create(name string) (File, error) {
	f, err := os.Create(fs.Resolve(name))
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Stat stats the file.
func (fs *fileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(fs.Resolve(name))
}

// Remove the file.
func (fs *fileSystem) Remove(name string) error {
	return os.Remove(fs.Resolve(name))
}
