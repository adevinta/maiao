package cmd

import (
	"context"
	"fmt"

	"github.com/adevinta/maiao/pkg/committemplate"
	"github.com/adevinta/maiao/pkg/log"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func prepare_commit_message(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing required commit message file path")
	}
	if len(args) >= 3 {
		// This is an ammend
		// See official git documentation: https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks
		// The prepare-commit-msg hook is run before the commit message editor is fired up but after
		// the default message is created. It lets you edit the default message before the commit author sees it.
		// This hook takes a few parameters: the path to the file that holds the commit message so far,
		// the type of commit, and the commit SHA-1 if this is an amended commit
		log.Logger.Debug("this is an ammend commit, skipping")
		return nil
	}
	repo, err := git.PlainOpenWithOptions(cmd.Flag("path").Value.String(), &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return err
	}
	err = committemplate.Prepare(context.Background(), repo, args[0])
	if err != nil {
		return err
	}
	return nil
}
