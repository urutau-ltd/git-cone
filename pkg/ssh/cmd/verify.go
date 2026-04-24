package cmd

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/urutau-ltd/git-cone/pkg/access"
	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/proto"
	"github.com/urutau-ltd/git-cone/pkg/sshutils"
	"github.com/urutau-ltd/git-cone/pkg/utils"
)

func verifyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "verify REPOSITORY",
		Short: "Verify repository integrity via git fsck",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			be := backend.FromContext(ctx)
			user := proto.UserFromContext(ctx)
			session := sshutils.SessionFromContext(ctx)

			repoName := utils.SanitizeRepo(args[0])

			// Check access: server admin or write-level collaborator.
			auth := be.AccessLevelForUser(ctx, repoName, user)
			if auth < access.ReadWriteAccess {
				fmt.Fprintf(cmd.ErrOrStderr(), "permission denied\n")
				if session != nil {
					_ = session.Exit(1)
				}
				return nil
			}

			// Resolve repo path.
			rr, err := be.Repository(ctx, repoName)
			if err != nil {
				return err
			}
			repo, err := rr.Open()
			if err != nil {
				return err
			}

			gitCmd := exec.CommandContext(ctx, "git", "-C", repo.Path, "fsck", "--full")
			if session != nil {
				gitCmd.Stdout = session
				gitCmd.Stderr = session.Stderr()
			} else {
				gitCmd.Stdout = cmd.OutOrStdout()
				gitCmd.Stderr = cmd.ErrOrStderr()
			}
			if err := gitCmd.Run(); err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					if session != nil {
						_ = session.Exit(exitErr.ExitCode()) //nolint:errcheck
					}
					return nil
				}
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", err)
				if session != nil {
					_ = session.Exit(1) //nolint:errcheck
				}
				return nil
			}

			if session != nil {
				_ = session.Exit(0) //nolint:errcheck
			}
			return nil
		},
	}
}
