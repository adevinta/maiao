package git

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/adevinta/maiao/pkg/system"
)

const (
	CommitMsgHook Hook = "commit-msg"
)

type Hook string

const (
	ApplyPatchMsgHook    Hook = "applypatch-msg"
	PostApplyPatchHook   Hook = "post-applypatch"
	PostCheckoutHook     Hook = "post-checkout"
	PostCommitHook       Hook = "post-commit"
	PostMergeHook        Hook = "post-merge"
	PostRewriteHook      Hook = "post-rewrite"
	PrepareCommitMsgHook Hook = "prepare-commit-msg"
	PreApplyPatchHook    Hook = "pre-applypatch"
	PreAutoGCHook        Hook = "pre-auto-gc"
	PrePushHook          Hook = "pre-push"
	PreRebaseHook        Hook = "pre-rebase"
)

func HookPath(gitDir string, hook Hook) string {
	commonDirPath := filepath.Join(gitDir, "commondir")
	_, err := system.DefaultFileSystem.Stat(commonDirPath)
	if err == nil {
		fd, err := system.DefaultFileSystem.Open(commonDirPath)
		if err == nil {
			defer fd.Close()
			bytes, err := ioutil.ReadAll(fd)
			if err == nil {
				gitDir = filepath.Join(gitDir, strings.TrimSpace(string(bytes)))
			}
		}
	}
	return filepath.Join(gitDir, "hooks", string(hook))
}
