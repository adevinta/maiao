package system

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	// DefaultFileSystem allows to change the implementation of the default file system, for example for tests
	DefaultFileSystem  afero.Fs = afero.NewOsFs()
	originalFileSystem          = DefaultFileSystem
)

func EnsureFileContent(fs afero.Fs, path string, reader io.Reader) error {
	dir := filepath.Dir(path)
	stat, err := fs.Stat(dir)
	if os.IsNotExist(err) {
		err = fs.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	if stat != nil && !stat.IsDir() {
		return fmt.Errorf("path %s is not a directory", dir)
	}
	fd, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = io.Copy(fd, reader)
	return err
}
