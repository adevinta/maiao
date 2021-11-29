package system

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertFileContents(t *testing.T, fs afero.Fs, path string, content string) bool {
	t.Helper()
	fd, err := fs.Open(path)
	require.NoError(t, err)
	bytes, err := ioutil.ReadAll(fd)
	require.NoError(t, err)
	return assert.Equal(t, content, string(bytes))
}

func AssertPathExists(t *testing.T, fs afero.Fs, path string, msgAndArgs ...interface{}) bool {
	t.Helper()
	_, err := fs.Stat(path)
	if os.IsNotExist(err) {
		return assert.Fail(t, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
	}
	return true
}

func AssertModePerm(t testing.TB, fs afero.Fs, path, mode string) bool {
	t.Helper()
	s, err := fs.Stat(path)
	require.NoError(t, err)
	return assert.Equal(t, mode, s.Mode().Perm().String())
}

func EnsureTestFileContent(t *testing.T, fs afero.Fs, path string, content string) {
	require.NoError(t, EnsureFileContent(fs, path, strings.NewReader(content)))
}
