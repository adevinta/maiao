package cmd

import (
	"github.com/spf13/cobra"
	"github.com/adevinta/maiao/pkg/gerrit"
	lgit "github.com/adevinta/maiao/pkg/git"
)

func install(cmd *cobra.Command, args []string) error {
	gitDir, err := lgit.FindGitDir(cmd.Flag("path").Value.String())
	if err != nil {
		return err
	}
	err = gerrit.Install(gitDir)
	if err != nil {
		return err
	}
	return nil
}
