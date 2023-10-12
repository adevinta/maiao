package cmd

import (
	"os"

	"github.com/adevinta/maiao/pkg/committemplate"
	"github.com/adevinta/maiao/pkg/gerrit"
	lgit "github.com/adevinta/maiao/pkg/git"
	"github.com/spf13/cobra"
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
	if os.Getenv("MAIAO_EXPERIMENTAL_COMMIT_TEMPLATE") == "true" {
		err = committemplate.Install(gitDir)
		if err != nil {
			return err
		}
	}
	return nil
}
