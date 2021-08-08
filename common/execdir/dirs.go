package execdir

import (
	"io/fs"
	"os"
	"path/filepath"
)

// GetFS returns the file system rooted in the exec directory.
func GetFS(name string) fs.FS {
	return os.DirFS(filepath.Join(d, name))
}

// ReadDir reads the related named directory and returns a list of directory entries.
func ReadDir(name string) ([]fs.DirEntry, error) {
	f, err := os.Open(filepath.Join(d, name))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.ReadDir(-1)
}
