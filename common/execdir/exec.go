package execdir

import (
	"os"
	"path/filepath"
)

var ex, d, name string

func init() {
	var err error
	ex, err = os.Executable()
	if err != nil {
		panic(err)
	}
	d, name = filepath.Split(ex)
}

// GetExec returns the path to the executable file.
func GetExec() string {
	return ex
}

// GetExecDir returns the directory of the excutable file.
func GetExecDir() string {
	return d
}

// GetExecName returns the name of the excutable file.
func GetExecName() string {
	return name
}
