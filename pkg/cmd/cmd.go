package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/99designs/keyring"
	"github.com/adevinta/maiao/pkg/git"
	"github.com/adevinta/maiao/pkg/log"
	"github.com/adevinta/maiao/pkg/version"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	changeIDEditorHelp = `Reproduces git interactions a human should do to add Change-Ids to commits
	Used as 'env GIT_EDITOR="git review add-change-id-editor" git rebase -i origin/master'
	This command selects changes all pickups to rewords and keeps the message intact for the
	Change-Id hook to be ran
	`
)

// NewCommand implements a new cobra command to run git-review
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "git-review [<targetBranch>]",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch cmd.Flag("verbose").Value.String() {
			case "0":
				log.Logger.SetLevel(logrus.FatalLevel)
			case "1":
				log.Logger.SetLevel(logrus.ErrorLevel)
			case "2":
				log.Logger.SetLevel(logrus.WarnLevel)
			case "3":
				log.Logger.SetLevel(logrus.InfoLevel)
			case "4":
				keyring.Debug = true
				log.Logger.SetLevel(logrus.DebugLevel)
			case "5":
				keyring.Debug = true
				log.Logger.SetLevel(logrus.TraceLevel)
			default:
				return fmt.Errorf("unexpected log level %s expecting 0-5", cmd.Flag("verbose").Value.String())
			}
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("too many arguments provided")
			}
			return nil
		},
		Version: version.Version,
	}

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		maiaoArgs := os.Getenv(git.RebaseArgsEnvVar)
		if maiaoArgs != "" {
			os.Unsetenv(git.RebaseArgsEnvVar)
			originalArgs := []string{}
			err := json.Unmarshal([]byte(maiaoArgs), &originalArgs)
			if err != nil {
				return err
			}
			err = rootCmd.ParseFlags(originalArgs)
			if err != nil {
				return err
			}
			err = rootCmd.Execute()
			if err != nil {
				return err
			}
			return nil
		}
		return review(cmd, args)
	}

	rootCmd.PersistentFlags().IntP("verbose", "v", 0, "the logging verbosity (0-5)")
	rootCmd.PersistentFlags().StringP("path", "C", ".", "Path of the repository to push reviews")
	rootCmd.PersistentFlags().BoolP("no-rebase", "R", false, "Don't rebase changes before submitting")
	rootCmd.PersistentFlags().StringP("topic", "t", "", "Topic to submit branch to")
	rootCmd.PersistentFlags().Bool("debug", false, "Run the command in debug mode")
	rootCmd.PersistentFlags().String("remote", "", "Specifies the remote the review should be done on. By default the tracking remote of the target branch is used")
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "install",
			Short: "Installs commit message hook to the repository",
			Long:  `Installs commit message hook to the repository`,
			RunE:  install,
		},
		&cobra.Command{
			Use:   "version",
			Short: "Installs commit message hook to the repository",
			Long:  `Installs commit message hook to the repository`,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(version.Version)
			},
		},
		&cobra.Command{
			Use:    "add-change-id-editor",
			Short:  "Handles rebase interactive file edition",
			Long:   changeIDEditorHelp,
			RunE:   rebaseEditor,
			Hidden: true,
		},
		&cobra.Command{
			Use:    "prepare-commit-message",
			Short:  "Initialises a commit message with Pull Request templates",
			RunE:   prepare_commit_message,
			Hidden: true,
		},
	)
	rootCmd.AddCommand()
	return rootCmd
}
