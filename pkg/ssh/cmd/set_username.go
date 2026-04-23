package cmd

import (
	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/sshutils"
	"github.com/spf13/cobra"
)

// SetUsernameCommand returns a command that sets the user's username.
func SetUsernameCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-username USERNAME",
		Short: "Set your username",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			be := backend.FromContext(ctx)
			pk := sshutils.PublicKeyFromContext(ctx)
			user, err := be.UserByPublicKey(ctx, pk)
			if err != nil {
				return err
			}

			return be.SetUsername(ctx, user.Username(), args[0])
		},
	}

	return cmd
}
