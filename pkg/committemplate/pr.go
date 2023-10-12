package committemplate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adevinta/maiao/pkg/git"
	lgit "github.com/adevinta/maiao/pkg/git"
	"github.com/adevinta/maiao/pkg/log"
	"github.com/adevinta/maiao/pkg/system"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Interface interface {
	Installed() bool
	Install() error
}

type PrepareCommitMessageHook struct {
	gitDir string
}

// Installed returned wether the gerrit hook message is installed
func (g *PrepareCommitMessageHook) Installed() bool {
	path := git.HookPath(g.gitDir, git.PrepareCommitMsgHook)
	_, err := system.DefaultFileSystem.Stat(path)
	log.Logger.WithFields(logrus.Fields{
		"gitDir":                       g.gitDir,
		"prepare-commit-msg path":      path,
		"prepare-commit-msg installed": err == nil,
	}).Debugf("err: %v", err)
	return err == nil
}

// Install installs the gerrit prepare commit message hook in a repository
func (g *PrepareCommitMessageHook) Install() error {
	path := git.HookPath(g.gitDir, git.PrepareCommitMsgHook)

	l := log.Logger.WithFields(logrus.Fields{
		"gitDir":                  g.gitDir,
		"prepare-commit-msg path": path,
	})
	l.Debug("creating 'prepare-commit-msg' hook")
	fd, err := system.DefaultFileSystem.Create(path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create prepare commit message hook file %s", path))
	}
	defer fd.Close()
	_, err = fd.Write([]byte(`#!/bin/sh
# This hook has been installed by maiao
# see https://github.com/adevinta/maiao
# This hook is used to add a pull request template to the commit message
exec git-review prepare-commit-message "$@"
`))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create prepare commit message hook file %s", path))
	}

	err = system.DefaultFileSystem.Chmod(path, 0755)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to set execution rights to message hook file %s", path))
	}
	return nil
}

// Installed returned wether the gerrit hook message is installed
func Installed(gitDir string) bool {
	g := &PrepareCommitMessageHook{gitDir}
	return g.Installed()
}

// Install installs the gerrit prepare commit message hook in a repository
func Install(gitDir string) error {
	g := &PrepareCommitMessageHook{gitDir}
	return g.Install()
}

func Prepare(ctx context.Context, repo lgit.Repository, path string) error {
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	root := wt.Filesystem.Root()
	templatePath := filepath.Join(root, ".github", "pull_request_template.md")
	ctx = log.WithContextFields(ctx, logrus.Fields{"template-path": templatePath})
	_, err = system.DefaultFileSystem.Stat(templatePath)
	if os.IsNotExist(err) {
		log.Logger.WithContext(ctx).Debugf("no pull request template found")
		return nil
	}
	if err != nil {
		return err
	}

	instructions, err := afero.ReadFile(system.DefaultFileSystem, path)
	if err != nil {
		return err
	}
	cfg, err := repo.Config()
	if err != nil {
		return err
	}
	for i := 0; i < len(instructions); i++ {
		if instructions[i] == '\n' {
			continue
		}
		if len(cfg.Core.CommentChar) == 0 {
			cfg.Core.CommentChar = "#"
		}
		if instructions[i] == cfg.Core.CommentChar[0] {
			break
		}
		log.Logger.WithContext(ctx).Debugf("commit message already contains a message, skipping")
		return nil
	}

	prTemplate, err := afero.ReadFile(system.DefaultFileSystem, templatePath)
	if err != nil {
		return err
	}
	fd, err := system.DefaultFileSystem.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	// Create a space for the commit title and split commit title and message by one empty line
	fd.Write([]byte("\n\n"))
	_, err = fd.Write(append(prTemplate, instructions...))
	if err != nil {
		return err
	}
	return nil
}
