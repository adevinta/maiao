package gerrit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/adevinta/maiao/pkg/log"
	"github.com/adevinta/maiao/pkg/system"
)

const (
	testHookScript = `#!/bin/bash
echo you downloaded me`
)

func setHookURL(url string) {
	commitMsgHookURL = url
}

func TestIsInstalled(t *testing.T) {
	fs := afero.NewMemMapFs()
	system.DefaultFileSystem = fs
	t.Cleanup(system.Reset)

	system.EnsureTestFileContent(t, fs, "/src/.git/hooks/commit-msg", "some-content")
	system.EnsureTestFileContent(t, fs, "/src/.git/worktrees/some-name/commondir", "../..")
	t.Run("when the hook is installed, installed returns True", func(t *testing.T) {
		assert.True(t, Installed("/src/.git/"))
	})
	t.Run("when running in a worktree, installed returns True", func(t *testing.T) {
		assert.True(t, Installed("/src/.git/worktrees/some-name"))
	})
	t.Run("when the hook is not installed, installed returns False", func(t *testing.T) {
		assert.False(t, Installed("/src/"))
	})
}

func TestInstall(t *testing.T) {
	// restore the original value at the end of the test
	t.Cleanup(system.Reset)
	defer setHookURL(commitMsgHookURL)
	s := httptest.NewServer(http.HandlerFunc(replyTestHookScript))
	defer s.Close()

	t.Run("with an invalid hook url, the installation fails", func(t *testing.T) {
		setHookURL("https://localhost:32")
		fs := afero.NewMemMapFs()
		system.DefaultFileSystem = fs

		system.EnsureTestFileContent(t, fs, hookPath("/src/some-repo/.git"), "#!/bin/bash\necho hello world")
		require.NoError(t, fs.Chmod(hookPath("/src/some-repo/.git"), 0644))

		assert.Error(t, Install("/src/some-repo/.git"))
		system.AssertPathExists(t, fs, hookPath("/src/some-repo/.git"))
		system.AssertFileContents(t, fs, hookPath("/src/some-repo/.git"), "#!/bin/bash\necho hello world")
		system.AssertModePerm(t, fs, hookPath("/src/some-repo/.git"), "-rw-r--r--")
	})
	t.Run("with a valid hook URL, no error is returned", func(t *testing.T) {
		setHookURL(s.URL + "/commit-msg-hook")
		t.Run("when the hooks directory does not exist", func(t *testing.T) {
			fs := afero.NewMemMapFs()
			system.DefaultFileSystem = fs
			testHookInstalled(t, fs, "/src/some-repo/.git")
		})
		t.Run("when the hooks directory already exists", func(t *testing.T) {
			fs := afero.NewMemMapFs()
			system.DefaultFileSystem = fs
			system.EnsureTestFileContent(t, fs, hookPath("/src/some-repo/.git"), "#!/bin/bash\necho hello world")

			testHookInstalled(t, fs, "/src/some-repo/.git")
		})
	})
}

func replyTestHookScript(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(testHookScript))
}

func testHookInstalled(t *testing.T, fs afero.Fs, path string) {
	t.Helper()
	t.Run("hook installation succeed", func(t *testing.T) {
		assert.NoError(t, Install(path))
		system.AssertPathExists(t, fs, hookPath(path))
		system.AssertFileContents(t, fs, hookPath(path), testHookScript)
		system.AssertModePerm(t, fs, hookPath(path), "-rwxr-xr-x")
	})
}

// get all logs when running tests
func init() {
	log.Logger.SetLevel(logrus.DebugLevel)
}
