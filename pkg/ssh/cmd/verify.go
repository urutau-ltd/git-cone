package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/urutau-ltd/git-cone/pkg/access"
	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/config"
	"github.com/urutau-ltd/git-cone/pkg/proto"
	"github.com/urutau-ltd/git-cone/pkg/sshutils"
	"github.com/urutau-ltd/git-cone/pkg/utils"
	"github.com/spf13/cobra"
)

func verifyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "verify REPOSITORY",
		Short: "Verify repository integrity via git fsck",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			be := backend.FromContext(ctx)
			cfg := config.FromContext(ctx)
			user := proto.UserFromContext(ctx)

			repoName := utils.SanitizeRepo(args[0])

			// Check access: server admin or write-level collaborator.
			auth := be.AccessLevelForUser(ctx, repoName, user)
			if auth < access.ReadWriteAccess {
				fmt.Fprintf(cmd.ErrOrStderr(), "permission denied\n")
				if sess := sshutils.SessionFromContext(ctx); sess != nil {
					_ = sess.Exit(1)
				}
				return nil
			}

			// Resolve repo path.
			if _, err := be.Repository(ctx, repoName); err != nil {
				return err
			}
			repoPath := filepath.Join(cfg.DataPath, "repos", repoName+".git")

			gitCmd := exec.CommandContext(ctx, "git", "-C", repoPath, "fsck", "--full")
			gitCmd.Stdout = cmd.OutOrStdout()
			gitCmd.Stderr = cmd.ErrOrStderr()
			if err := gitCmd.Run(); err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					if sess := sshutils.SessionFromContext(ctx); sess != nil {
						_ = sess.Exit(exitErr.ExitCode()) //nolint:errcheck
					}
					return nil
				}
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", err)
				if sess := sshutils.SessionFromContext(ctx); sess != nil {
					_ = sess.Exit(1) //nolint:errcheck
				}
				return nil
			}

			if sess := sshutils.SessionFromContext(ctx); sess != nil {
				_ = sess.Exit(0) //nolint:errcheck
			}
			return nil
		},
	}
}
