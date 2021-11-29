package gerrit

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/adevinta/maiao/pkg/log"
	"github.com/adevinta/maiao/pkg/system"
)

const (
	gitCommitMsgHookPath = "hooks/commit-msg"
	gitHubURL            = "https://raw.githubusercontent.com"
	repo                 = "GerritCodeReview/gerrit"
	commitHash           = "43d985a2a15a7d59d42e19ffd60d41c0de6c3e59"
	commitMsgHookPath    = "gerrit-server/src/main/resources/com/google/gerrit/server/tools/root/hooks/commit-msg"
)

var commitMsgHookURL = fmt.Sprintf("%s/%s/%s/%s", gitHubURL, repo, commitHash, commitMsgHookPath)

type Interface interface {
	Installed() bool
	Install() error
}

type Gerrit struct {
	gitDir string
}

func HookURL() string {
	return commitMsgHookURL
}

// Installed returned wether the gerrit hook message is installed
func (g *Gerrit) Installed() bool {
	path := hookPath(g.gitDir)
	_, err := system.DefaultFileSystem.Stat(path)
	log.Logger.WithFields(logrus.Fields{
		"gitDir":               g.gitDir,
		"commit-hook path":     path,
		"commit-msg installed": err == nil,
	}).Debugf("err: %v", err)
	return err == nil
}

// Install installs the gerrit commit message hook in a repository
func (g *Gerrit) Install() error {
	path := hookPath(g.gitDir)

	l := log.Logger.WithFields(logrus.Fields{
		"gitDir":           g.gitDir,
		"commit-hook path": path,
		"download-url":     commitMsgHookURL,
	})
	l.Debug("downloading commit message hook")
	r, err := http.Get(commitMsgHookURL)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to download commit message hook from %s", commitMsgHookURL))
	}
	defer r.Body.Close()
	l.Debugf("downloaded commit message hook")
	d := filepath.Dir(path)
	s, err := system.DefaultFileSystem.Stat(d)
	if err != nil {
		if os.IsNotExist(err) {
			l.Debugf("created hooks directory %s", d)
			err = system.DefaultFileSystem.MkdirAll(d, 0777)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to create hooks directory %s", d))
			}
		} else {
			return errors.Wrap(err, fmt.Sprintf("failed to create hooks directory %s", d))
		}
	} else {
		if !s.IsDir() {
			return fmt.Errorf("could not create commit message hook, %s is not a directory", d)
		}
	}
	fd, err := system.DefaultFileSystem.Create(path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create commit message hook file %s", path))
	}
	defer fd.Close()
	_, err = io.Copy(fd, r.Body)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write commit message hook file %s", path))
	}

	err = system.DefaultFileSystem.Chmod(path, 0755)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to set execution rights to message hook file %s", path))
	}
	return nil
}

// Installed returned wether the gerrit hook message is installed
func Installed(gitDir string) bool {
	g := &Gerrit{gitDir}
	return g.Installed()
}

// Install installs the gerrit commit message hook in a repository
func Install(gitDir string) error {
	g := &Gerrit{gitDir}
	return g.Install()
}

func hookPath(gitDir string) string {
	commonDirPath := filepath.Join(gitDir, "commondir")
	_, err := system.DefaultFileSystem.Stat(commonDirPath)
	if err == nil {
		fd, err := system.DefaultFileSystem.Open(commonDirPath)
		if err == nil {
			defer fd.Close()
			bytes, err := ioutil.ReadAll(fd)
			if err == nil {
				fmt.Println(gitDir, strings.TrimSpace(string(bytes)), filepath.Join(gitDir, strings.TrimSpace(string(bytes))))
				gitDir = filepath.Join(gitDir, strings.TrimSpace(string(bytes)))
			}
		}
	}
	return filepath.Join(gitDir, gitCommitMsgHookPath)
}
