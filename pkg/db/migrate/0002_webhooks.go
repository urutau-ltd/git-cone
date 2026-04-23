package migrate

import (
	"context"

	"github.com/urutau-ltd/git-cone/pkg/db"
)

const (
	webhooksName    = "webhooks"
	webhooksVersion = 2
)

var webhooks = Migration{
	Name:    webhooksName,
	Version: webhooksVersion,
	Migrate: func(ctx context.Context, tx *db.Tx) error {
		return migrateUp(ctx, tx, webhooksVersion, webhooksName)
	},
	Rollback: func(ctx context.Context, tx *db.Tx) error {
		return migrateDown(ctx, tx, webhooksVersion, webhooksName)
	},
}
