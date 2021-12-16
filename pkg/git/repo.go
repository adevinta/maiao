package git

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adevinta/maiao/pkg/system"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

const (
	RebaseArgsEnvVar = "MAIAO_REBASE_ARGUMENTS"
)

type Repository interface {
	Head() (*plumbing.Reference, error)
	Remote(name string) (*git.Remote, error)
	Push(o *git.PushOptions) error
	Branches() (storer.ReferenceIter, error)
	Config() (*config.Config, error)
	Fetch(o *git.FetchOptions) error
	Log(o *git.LogOptions) (object.CommitIter, error)
	ResolveRevision(rev plumbing.Revision) (*plumbing.Hash, error)
	Worktree() (*git.Worktree, error)
}

func MergeBase(ctx context.Context, repo Repository, base, head plumbing.Revision) (plumbing.Hash, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return plumbing.Hash{}, err
	}
	b, err := repo.ResolveRevision(base)
	if err != nil {
		return plumbing.Hash{}, err
	}
	h, err := repo.ResolveRevision(head)
	if err != nil {
		return plumbing.Hash{}, err
	}
	c := exec.Command("git", "-C", wt.Filesystem.Root(), "merge-base", b.String(), h.String())
	out := bytes.Buffer{}
	c.Stdout = &out
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		return plumbing.Hash{}, err
	}
	return plumbing.NewHash(strings.Trim(out.String(), " \n")), nil
}

func RebaseCommits(ctx context.Context, repo Repository, base, onto plumbing.Hash, todo string) error {
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	fd, err := ioutil.TempFile("", "rebase-todo-")
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(fd.Name())
	}()
	todo = todo + "\n" + "exec " + os.Args[0]
	_, err = fd.Write([]byte(todo + "\n"))
	if err != nil {
		fd.Close()
		return err
	}
	err = fd.Close()
	if err != nil {
		return err
	}

	c := exec.Command("git", "-C", wt.Filesystem.Root(), "rebase", "-i", "--onto", onto.String(), base.String())
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	args, err := json.Marshal(os.Args[1:])
	if err != nil {
		return err
	}
	c.Env = append(os.Environ(), "GIT_EDITOR="+os.Args[0]+" add-change-id-editor "+fd.Name(), RebaseArgsEnvVar+"="+string(args))
	err = c.Run()
	if err != nil {
		return err
	}
	return nil
}

func FindGitDir(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = abs
	for {
		gitDir := filepath.Join(path, ".git")
		stat, err := system.DefaultFileSystem.Stat(gitDir)
		if os.IsNotExist(err) {
			path = filepath.Dir(path)
			if path == "/" || path == "" {
				return "", errors.New("unable to find git directory")
			}
			continue
		}
		if stat.IsDir() {
			return gitDir, nil
		}
		return "", errors.New("gitdir is not a directory")

	}
}
