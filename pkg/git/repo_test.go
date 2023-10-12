package git

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adevinta/maiao/pkg/log"
	"github.com/adevinta/maiao/pkg/system"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindGitDir(t *testing.T) {
	t.Cleanup(system.Reset)
	system.DefaultFileSystem = afero.NewMemMapFs()

	require.NoError(t, system.DefaultFileSystem.MkdirAll("/some/path/to/repository/.git", 0700))
	fd, err := system.DefaultFileSystem.Create("/some/path/to/something/else/.git")
	require.NoError(t, err)
	fd.Close()

	dir, err := FindGitDir("/some/path/to/repository/.git")
	assert.NoError(t, err)
	assert.Equal(t, "/some/path/to/repository/.git", dir)

	dir, err = FindGitDir("/some/path/to/repository/.git/something/else")
	assert.NoError(t, err)
	assert.Equal(t, "/some/path/to/repository/.git", dir)

	dir, err = FindGitDir("/some/path/to/repository/something/else")
	assert.NoError(t, err)
	assert.Equal(t, "/some/path/to/repository/.git", dir)

	dir, err = FindGitDir("/some/path/to/something/else/.git")
	assert.Error(t, err)
	assert.Equal(t, "", dir)

}

func TestFindGitDirWithWorkDir(t *testing.T) {
	repo, err := os.MkdirTemp("", "maiao-worktree-test-git-dir-repo")
	require.NoError(t, err)
	worktree, err := os.MkdirTemp("", "maiao-worktree-test-git-dir-worktree")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(repo)
		os.RemoveAll(worktree)
	})
	cmd(t, "git", "-C", repo, "init")
	cmd(t, "git", "-C", repo, "commit", "--allow-empty", "-m", "initial commit")
	// Use show toplevel as tempfiles on mac may be in /var/folders/... while mounted volumes are in /private/var/folders/...
	// Internally, when creating worktrees, git uses the toplevel to point to the worktree dir
	// Hint it to do the same
	repoGitDir, err := FindGitDir(cmdOutput(t, "git", "-C", repo, "rev-parse", "--show-toplevel"))
	assert.NoError(t, err)
	cmd(t, "git", "-C", repo, "worktree", "add", worktree)
	worktreeGitDir, err := FindGitDir(worktree)
	assert.NoError(t, err)

	if !strings.HasPrefix(worktreeGitDir, filepath.Join(repoGitDir, "worktrees")) {
		assert.Failf(t, "Unexpected prefix", "worktree git dir '%s' should be included in the repo git dir '%s'", worktreeGitDir, filepath.Join(repoGitDir, ".git", "worktrees"))
	}
}

func cmd(t testing.TB, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		t.Errorf("failed to run command %s: %s", strings.Join(append([]string{cmd}, args...), " "), err.Error())
		t.FailNow()
	}
}

func cmdOutput(t testing.TB, cmd string, args ...string) string {
	c := exec.Command(cmd, args...)
	b := bytes.NewBuffer(nil)
	c.Stdout = b
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		t.Errorf("failed to run command %s: %s", strings.Join(append([]string{cmd}, args...), " "), err.Error())
		t.FailNow()
	}
	return strings.Trim(b.String(), " \n")
}

func commitFile(t *testing.T, dir, path, content, message string) string {
	fd, err := os.Create(filepath.Join(dir, path))
	assert.NoError(t, err)
	fd.Write([]byte(content))
	cmd(t, "git", "-C", dir, "add", path)
	cmd(t, "git", "-C", dir, "commit", "-m", message)
	return cmdOutput(t, "git", "-C", dir, "rev-parse", "HEAD")
}

// func TestRepo(t *testing.T) {
// 	t.Run("in a standard directory", func(t *testing.T) {
// 		r, err := PlainOpen("../..")
// 		assert.NoError(t, err)
// 		assert.NotNil(t, r)
// 		assert.Equal(t, "../../.git", r.CommonGitDir())
// 		t.Run("head is retrievable", func(t *testing.T) {
// 			head, err := r.Head()
// 			assert.NoError(t, err)
// 			assert.NotNil(t, head)
// 			assert.Equal(t, cmdOutput(t, "git", "rev-parse", "--symbolic-full-name", "HEAD"), head.Name().String())
// 		})
// 	})
// 	t.Run("inside a worktree", func(t *testing.T) {
// 		branch := uuid.New().String()
// 		wt := "tests/worktree-tests/" + branch
// 		cmd(t, "git", "-C", "..", "worktree", "add", wt, "-B", branch)
// 		defer func() {
// 			cmd(t, "rm", "-rf", "../"+wt)
// 			cmd(t, "git", "-C", "../..", "worktree", "prune")
// 			cmd(t, "git", "branch", "-D", branch)
// 		}()

// 		r, err := PlainOpen("../" + wt + "/")

// 		assert.NoError(t, err)
// 		assert.NotNil(t, r)
// 		abs, _ := filepath.Abs("../../.git")
// 		assert.Equal(t, abs, r.CommonGitDir())
// 		t.Run("head is retrievable", func(t *testing.T) {
// 			t.Skip("gopkg.in/src-d/go-git.v4 does not support working in git worktrees")
// 			head, err := r.Head()
// 			assert.NoError(t, err)
// 			assert.NotNil(t, head)
// 			assert.Equal(t, branch, head.Name())
// 		})
// 	})
// }

// func TestMergeBase(t *testing.T) {
// 	dir, err := ioutil.TempDir("", t.Name())
// 	assert.NoError(t, err)
// 	fmt.Println(dir)
// 	defer os.RemoveAll(dir)
// 	cmd(t, "git", "init", dir)
// 	cmd(t, "git", "-C", dir, "config", "commit.gpgsign", "false")
// 	commitFile(t, dir, "README.md", strings.Join(faker.Hacker().Phrases(), "\n\n"), faker.Hacker().SaySomethingSmart())
// 	c1 := commitFile(t, dir, "file1", strings.Join(faker.Hacker().Phrases(), "\n\n"), faker.Hacker().SaySomethingSmart())
// 	cmd(t, "git", "-C", dir, "checkout", "-b", "old-branch")
// 	commitFile(t, dir, "file1", strings.Join(faker.Hacker().Phrases(), "\n\n"), faker.Hacker().SaySomethingSmart())
// 	cmd(t, "git", "-C", dir, "checkout", "-b", "new-branch", c1)
// 	commitFile(t, dir, "file1", strings.Join(faker.Hacker().Phrases(), "\n\n"), faker.Hacker().SaySomethingSmart())
// 	repo, err := PlainOpen(dir)
// 	assert.NoError(t, err)
// 	assert.Equal(t, c1, repo.MergeBase("old-branch", "new-branch"))
// 	assert.Equal(t, "", repo.MergeBase("old-branch", "some-branch"))
// }

// get all logs when running tests
func init() {
	log.Logger.SetLevel(logrus.DebugLevel)
}
