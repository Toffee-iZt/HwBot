package execdir

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

// Open opens file in the exec directory.
func Open(name string) (fs.File, error) {
	return GetFS("").Open(name)
}

// OpenWriter opens file with write flag in the exec directory.
func OpenWriter(name string) (*os.File, error) {
	return os.OpenFile(filepath.Join(d, name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
}

// LoadJSON opens and parses json file in the exec directory.
func LoadJSON(name string, dst interface{}) error {
	f, err := Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(dst)
}

// SaveJSON opens file in the exec directory and writes the json encoding.
func SaveJSON(name string, v interface{}) error {
	f, err := OpenWriter(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}
