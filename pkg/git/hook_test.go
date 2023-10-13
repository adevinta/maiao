package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHookPath(t *testing.T) {
	repo, err := os.MkdirTemp("", "maiao-test-repo")
	require.NoError(t, err)
	worktree, err := os.MkdirTemp("", "maiao-test-worktree")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(repo)
		os.RemoveAll(worktree)
	})

	cmd(t, "git", "-C", repo, "init")
	cmd(t, "git", "-C", repo, "commit", "--allow-empty", "-m", "initial commit")
	cmd(t, "git", "-C", repo, "worktree", "add", worktree)

	// Use show toplevel as tempfiles on mac may be in /var/folders/... while mounted volumes are in /private/var/folders/...
	// Internally, when creating worktrees, git uses the toplevel to point to the worktree dir
	// Hint it to do the same
	repoGitDir, err := FindGitDir(cmdOutput(t, "git", "-C", repo, "rev-parse", "--show-toplevel"))
	require.NoError(t, err)
	worktreeGitDir, err := FindGitDir(cmdOutput(t, "git", "-C", worktree, "rev-parse", "--show-toplevel"))
	require.NoError(t, err)

	assert.Equal(
		t,
		filepath.Join(repoGitDir, "hooks", "commit-msg"),
		HookPath(
			filepath.Join(cmdOutput(t, "git", "-C", repo, "rev-parse", "--show-toplevel"), ".git"),
			CommitMsgHook,
		),
	)

	assert.Equal(
		t,
		filepath.Join(repoGitDir, "hooks", "commit-msg"),
		HookPath(
			worktreeGitDir,
			CommitMsgHook,
		),
	)
}
