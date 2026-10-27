package cmd

import (
	"github.com/spf13/cobra"
	"github.com/urutau-ltd/git-cone/pkg/access"
	"github.com/urutau-ltd/git-cone/pkg/backend"
)

func deleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete REPOSITORY",
		Aliases:           []string{"del", "remove", "rm"},
		Short:             "Delete a repository",
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: requireRepo(access.ReadWriteAccess),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			be := backend.FromContext(ctx)
			return be.DeleteRepository(ctx, repoArg(args))
		},
	}

	return cmd
}
