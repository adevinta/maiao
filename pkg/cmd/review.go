package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/adevinta/maiao/pkg/committemplate"
	"github.com/adevinta/maiao/pkg/gerrit"
	lgit "github.com/adevinta/maiao/pkg/git"
	"github.com/adevinta/maiao/pkg/maiao"
	"github.com/adevinta/maiao/pkg/prompt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

const (
	hookMissing          = "commit message hook is missing, do you want to install it automatically?"
	prHookMissing        = "PR template injection hook is missing, do you want to install it automatically?"
	noAutoInstallHookFmt = "You are missing change ids in your commits. \nPlease install the commit hook by running\n`curl -o .git/hooks/commit-msg %s && chmod +x .git/hooks/commit-msg`"
)

func review(cmd *cobra.Command, args []string) error {
	path := cmd.Flag("path").Value.String()
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return err
	}
	branch := ""
	if len(args) > 0 {
		branch = args[0]
	}
	gitDir, err := lgit.FindGitDir(path)
	if err != nil {
		return err
	}
	if !gerrit.Installed(gitDir) {
		if prompt.YesNo(hookMissing) {
			err = gerrit.Install(gitDir)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(fmt.Sprintf(noAutoInstallHookFmt, gerrit.HookURL()))
			return nil
		}
	}
	if os.Getenv("MAIAO_EXPERIMENTAL_COMMIT_TEMPLATE") == "true" {
		if !committemplate.Installed(gitDir) {
			if prompt.YesNo(prHookMissing) {
				err = committemplate.Install(gitDir)
				if err != nil {
					return err
				}
			}
		}
	}
	return maiao.Review(context.Background(), repo, maiao.ReviewOptions{
		Remote:     cmd.Flag("remote").Value.String(),
		SkipRebase: cmd.Flag("no-rebase").Value.String() != "false",
		Topic:      cmd.Flag("topic").Value.String(),
		Branch:     branch,
	})
}
