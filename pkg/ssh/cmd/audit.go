package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/urutau-ltd/git-cone/pkg/backend"
	"github.com/urutau-ltd/git-cone/pkg/db"
	"github.com/urutau-ltd/git-cone/pkg/proto"
	"github.com/urutau-ltd/git-cone/pkg/sshutils"
	"github.com/urutau-ltd/git-cone/pkg/version"
	gossh "golang.org/x/crypto/ssh"
)

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
			cmd.Printf("pubkey-age:  %s\n", publicKeyAge(ctx, pk))
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

func publicKeyAge(ctx context.Context, pk gossh.PublicKey) string {
	if pk == nil {
		return "unknown"
	}

	dbx, _ := ctx.Value(db.ContextKey).(*db.DB)
	if dbx == nil {
		return "unknown"
	}

	var createdAt string
	err := dbx.TransactionContext(ctx, func(tx *db.Tx) error {
		query := tx.Rebind(`SELECT created_at FROM public_keys WHERE public_key = ? LIMIT 1`)
		return tx.GetContext(ctx, &createdAt, query, sshutils.MarshalAuthorizedKey(pk))
	})
	if err != nil || createdAt == "" {
		return "unknown"
	}

	ts, ok := parseAuditTime(createdAt)
	if !ok {
		return "unknown"
	}

	days := int(time.Since(ts).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return fmt.Sprintf("%d days", days)
}

func parseAuditTime(value string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, layout := range layouts {
		ts, err := time.Parse(layout, value)
		if err == nil {
			return ts, true
		}
	}
	return time.Time{}, false
}
