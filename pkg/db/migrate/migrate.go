package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"charm.land/log/v2"
	"github.com/urutau-ltd/git-cone/pkg/db"
)

// MigrateFunc is a function that executes a migration.
type MigrateFunc func(ctx context.Context, tx *db.Tx) error //nolint:revive

// Migration is a struct that contains the name of the migration and the
// function to execute it.
type Migration struct {
	Version  int64
	Name     string
	Migrate  MigrateFunc
	Rollback MigrateFunc
}

// Migrations is a database model to store migrations.
type Migrations struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Version int64  `db:"version"`
}

func (Migrations) schema(driverName string) string {
	switch driverName {
	case "sqlite3", "sqlite":
		return `CREATE TABLE IF NOT EXISTS migrations (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				version INTEGER NOT NULL UNIQUE
			);
		`
	default:
		panic("unknown driver")
	}
}

// Migrate runs the migrations.
func Migrate(ctx context.Context, dbx *db.DB) error {
	logger := log.FromContext(ctx).WithPrefix("migrate")
	return dbx.TransactionContext(ctx, func(tx *db.Tx) error {
		if !hasTable(tx, "migrations") {
			if _, err := tx.ExecContext(ctx, Migrations{}.schema(tx.DriverName())); err != nil {
				return err
			}
		}

		var migrs Migrations
		if err := tx.Get(&migrs, tx.Rebind("SELECT * FROM migrations ORDER BY version DESC LIMIT 1")); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}

		for _, m := range migrations {
			if m.Version <= migrs.Version {
				continue
			}

			logger.Infof("running migration %d. %s", m.Version, m.Name)
			if err := m.Migrate(ctx, tx); err != nil {
				return err
			}

			if _, err := tx.ExecContext(ctx, tx.Rebind("INSERT INTO migrations (name, version) VALUES (?, ?)"), m.Name, m.Version); err != nil {
				return err
			}
		}

		return nil
	})
}

// Rollback rolls back a migration.
func Rollback(ctx context.Context, dbx *db.DB) error {
	logger := log.FromContext(ctx).WithPrefix("migrate")
	return dbx.TransactionContext(ctx, func(tx *db.Tx) error {
		var migrs Migrations
		if err := tx.Get(&migrs, tx.Rebind("SELECT * FROM migrations ORDER BY version DESC LIMIT 1")); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("there are no migrations to rollback: %w", err)
			}
		}

		if migrs.Version == 0 || len(migrations) < int(migrs.Version) {
			return fmt.Errorf("there are no migrations to rollback")
		}

		m := migrations[migrs.Version-1]
		logger.Infof("rolling back migration %d. %s", m.Version, m.Name)
		if err := m.Rollback(ctx, tx); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, tx.Rebind("DELETE FROM migrations WHERE version = ?"), migrs.Version); err != nil {
			return err
		}

		return nil
	})
}

func hasTable(tx *db.Tx, tableName string) bool {
	query := tx.Rebind("SELECT name FROM sqlite_master WHERE type='table' AND name=?")
	var name string
	err := tx.Get(&name, query, tableName)
	return err == nil
}
