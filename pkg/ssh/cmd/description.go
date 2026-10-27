package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/urutau-ltd/git-cone/pkg/access"
	"github.com/urutau-ltd/git-cone/pkg/backend"
)

func descriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "description REPOSITORY [DESCRIPTION]",
		Aliases:           []string{"desc"},
		Short:             "Set or get the description for a repository",
		Args:              cobra.MinimumNArgs(1),
		PersistentPreRunE: requireRepo(access.ReadOnlyAccess),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			be := backend.FromContext(ctx)
			rr, _, err := repoFromContext(ctx)
			if err != nil {
				return err
			}

			switch len(args) {
			case 1:
				cmd.Println(rr.Description())
			default:
				if err := requireLevel(ctx, access.ReadWriteAccess); err != nil {
					return err
				}
				if err := be.SetDescription(ctx, rr.Name(), strings.Join(args[1:], " ")); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}
