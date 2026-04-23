package cmd

import (
	"fmt"

	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/proto"
	"github.com/urutau-ltd/git-cone/pkg/sshutils"
	"github.com/urutau-ltd/git-cone/pkg/version"
	"github.com/spf13/cobra"
	gossh "golang.org/x/crypto/ssh"
)

// AuditCommand returns a command that prints session and server audit info.
func AuditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "audit",
		Short: "Show session and server audit info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			be := backend.FromContext(ctx)
			user := proto.UserFromContext(ctx)
			pk := sshutils.PublicKeyFromContext(ctx)

			ver := version.Version
			if ver == "" {
				ver = "unknown"
			}
			cmd.Printf("server:      git-cone %s\n", ver)

			if user == nil && pk == nil {
				// Unauthenticated/keyless session: print only the server line.
				return nil
			}

			if user != nil {
				line := fmt.Sprintf("user:        %s", user.Username())
				if user.IsAdmin() {
					line += " (admin)"
				}
				cmd.Println(line)
			}

			if pk != nil {
				cmd.Printf("pubkey:      %s %s\n", pk.Type(), gossh.FingerprintSHA256(pk))
			}
			cmd.Printf("pubkey-age:  unknown\n")
			cmd.Printf("cipher:      (not available in session context)\n")
			cmd.Printf("kex:         (not available in session context)\n")

			// Count repos owned and repos where user is a collaborator.
			owned, collab := 0, 0
			if user != nil && be != nil {
				repos, _ := be.Repositories(ctx)
				for _, r := range repos {
					if r.UserID() == user.ID() {
						owned++
					} else {
						_, isCollab, _ := be.IsCollaborator(ctx, r.Name(), user.Username())
						if isCollab {
							collab++
						}
					}
				}
			}
			cmd.Printf("repos-owned: %d\n", owned)
			cmd.Printf("repos-collab:%d\n", collab)

			return nil
		},
	}
}
