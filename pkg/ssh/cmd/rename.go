package cmd

import (
	"github.com/spf13/cobra"
	"github.com/urutau-ltd/git-cone/pkg/access"
	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/proto"
	"github.com/urutau-ltd/git-cone/pkg/utils"
)

func renameCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "rename REPOSITORY NEW_NAME",
		Aliases:           []string{"mv", "move"},
		Short:             "Rename an existing repository",
		Args:              cobra.ExactArgs(2),
		PersistentPreRunE: requireRepo(access.ReadWriteAccess),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			be := backend.FromContext(ctx)
			oldName := repoArg(args)
			newName := utils.SanitizeRepo(args[1])

			// The destination is gated like create: it may already be
			// somebody else's repository.
			if repoAccessLevel(ctx, newName) < access.ReadWriteAccess {
				return proto.ErrUnauthorized
			}

			return be.RenameRepository(ctx, oldName, newName)
		},
	}

	return cmd
}
