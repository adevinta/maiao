package committemplate

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	lgit "github.com/adevinta/maiao/pkg/git"
	"github.com/adevinta/maiao/pkg/system"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fsStub struct {
	billy.Filesystem
	rootFunc func() string
}

func (f *fsStub) Root() string {
	if f.rootFunc == nil {
		return ""
	}
	return f.rootFunc()
}

type repoStub struct {
	lgit.Repository
	workteeFunc func() (*git.Worktree, error)
	configFunc  func() (*config.Config, error)
}

func (r *repoStub) Worktree() (*git.Worktree, error) {
	if r.workteeFunc == nil {
		return nil, errors.New("Worktree not implemented")
	}
	return r.workteeFunc()
}

func (r *repoStub) Config() (*config.Config, error) {
	if r.configFunc == nil {
		return nil, errors.New("Config not implemented")
	}
	return r.configFunc()
}

func TestInstalled(t *testing.T) {
	fs := afero.NewMemMapFs()
	system.DefaultFileSystem = fs
	t.Cleanup(system.Reset)

	hookPath := "/path/to/git/dir/hooks/prepare-commit-msg"
	gitDir := "/path/to/git/dir"
	g := &PrepareCommitMessageHook{
		gitDir: gitDir,
	}

	// Hook not installed
	fs.Remove(hookPath)
	assert.False(t, g.Installed())

	// Hook installed
	fs.MkdirAll(filepath.Dir(hookPath), 0755)
	f, err := fs.Create(hookPath)
	assert.NoError(t, err)
	f.Close()
	assert.True(t, g.Installed())
}

func TestInstall(t *testing.T) {
	fs := afero.NewMemMapFs()
	system.DefaultFileSystem = fs
	t.Cleanup(system.Reset)

	gitDir := "/path/to/git/dir"

	g := &PrepareCommitMessageHook{
		gitDir: gitDir,
	}

	err := g.Install()
	assert.NoError(t, err)

	hookPath := filepath.Join(gitDir, "hooks", "prepare-commit-msg")
	exists, err := afero.Exists(fs, hookPath)
	assert.NoError(t, err)
	assert.True(t, exists)

	hookContents, err := afero.ReadFile(fs, hookPath)
	assert.NoError(t, err)
	expectedContents := `#!/bin/sh
# This hook has been installed by maiao
# see https://github.com/adevinta/maiao
# This hook is used to add a pull request template to the commit message
exec git-review prepare-commit-message "$@"
`
	assert.Equal(t, expectedContents, string(hookContents))
}

func TestPrepare(t *testing.T) {
	fs := afero.NewMemMapFs()
	system.DefaultFileSystem = fs
	t.Cleanup(system.Reset)

	repo := &repoStub{
		configFunc: func() (*config.Config, error) {
			cfg := config.NewConfig()
			cfg.Core.CommentChar = "#"
			return cfg, nil
		},
		workteeFunc: func() (*git.Worktree, error) {
			return &git.Worktree{
				Filesystem: &fsStub{
					rootFunc: func() string {
						return "/path/to/git/dir"
					},
				},
			}, nil
		},
	}

	t.Run("When the PR template does not exist", func(t *testing.T) {
		t.Run("When the commit message file does not exist", func(t *testing.T) {
			assert.NoError(t, Prepare(context.Background(), repo, "/path/to/git/dir/.git/COMMIT_MSG"))
		})
		t.Run("When the commit message file does not exist", func(t *testing.T) {
			fd, err := fs.Create("/path/to/git/dir/.git/COMMIT_MSG")
			require.NoError(t, err)
			_, err = fd.Write([]byte("commit message"))
			require.NoError(t, err)
			require.NoError(t, fd.Close())
			assert.NoError(t, Prepare(context.Background(), repo, "/path/to/git/dir/.git/COMMIT_MSG"))
		})
	})

	t.Run("When the PR template exists", func(t *testing.T) {
		fs = afero.NewMemMapFs()
		system.DefaultFileSystem = fs
		fd, err := fs.Create("/path/to/git/dir/.github/pull_request_template.md")
		require.NoError(t, err)
		_, err = fd.Write([]byte("PR template"))
		require.NoError(t, err)
		require.NoError(t, fd.Close())

		t.Run("When the commit message file does not exist", func(t *testing.T) {
			assert.Error(t, Prepare(context.Background(), repo, "/path/to/git/dir/.git/COMMIT_MSG"))
		})
		t.Run("When the commit message file does not exist", func(t *testing.T) {
			fd, err := fs.Create("/path/to/git/dir/.git/COMMIT_MSG")
			require.NoError(t, err)
			_, err = fd.Write([]byte("\n# instructions for commit message"))
			require.NoError(t, err)
			require.NoError(t, fd.Close())
			assert.NoError(t, Prepare(context.Background(), repo, "/path/to/git/dir/.git/COMMIT_MSG"))
		})
	})

	testPrepareCommitMessageContent(t, "#", "\n# instructions for commit message", "PR template", "\n\nPR template\n# instructions for commit message")
	testPrepareCommitMessageContent(t, "#", " pre-defined commit message\n# instructions for commit message", "PR template", " pre-defined commit message\n# instructions for commit message")
	testPrepareCommitMessageContent(t, "^", " pre-defined commit message\n^ instructions for commit message", "PR template", " pre-defined commit message\n^ instructions for commit message")
}

func testPrepareCommitMessageContent(t *testing.T, commentChar, messageContent, templateContent, expectedContent string) {
	t.Helper()
	fs := afero.NewMemMapFs()
	system.DefaultFileSystem = fs
	t.Cleanup(system.Reset)

	repo := &repoStub{
		configFunc: func() (*config.Config, error) {
			cfg := config.NewConfig()
			cfg.Core.CommentChar = commentChar
			return cfg, nil
		},
		workteeFunc: func() (*git.Worktree, error) {
			return &git.Worktree{
				Filesystem: &fsStub{
					rootFunc: func() string {
						return "/path/to/git/dir"
					},
				},
			}, nil
		},
	}

	fd, err := fs.Create("/path/to/git/dir/.git/COMMIT_MSG")
	require.NoError(t, err)
	_, err = fd.Write([]byte(messageContent))
	require.NoError(t, err)
	require.NoError(t, fd.Close())

	fd, err = fs.Create("/path/to/git/dir/.github/pull_request_template.md")
	require.NoError(t, err)
	_, err = fd.Write([]byte(templateContent))
	require.NoError(t, err)
	require.NoError(t, fd.Close())

	assert.NoError(t, Prepare(context.Background(), repo, "/path/to/git/dir/.git/COMMIT_MSG"))

	fd, err = fs.Open("/path/to/git/dir/.git/COMMIT_MSG")
	require.NoError(t, err)
	content, err := afero.ReadAll(fd)
	require.NoError(t, err)
	require.NoError(t, fd.Close())

	assert.Equal(t, expectedContent, string(content))
}
