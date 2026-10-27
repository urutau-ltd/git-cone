package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Handler is a database handler.
type Handler interface {
	Rebind(string) string

	Select(any, string, ...any) error
	Get(any, string, ...any) error
	Queryx(string, ...any) (*sqlx.Rows, error)
	QueryRowx(string, ...any) *sqlx.Row
	Exec(string, ...any) (sql.Result, error)

	SelectContext(context.Context, any, string, ...any) error
	GetContext(context.Context, any, string, ...any) error
	QueryxContext(context.Context, string, ...any) (*sqlx.Rows, error)
	QueryRowxContext(context.Context, string, ...any) *sqlx.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}
