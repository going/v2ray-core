package connector

import (
	. "context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var _ Q = &sqlx.DB{}
var _ Q = &sqlx.Tx{}

type Connector interface {
	Invoke(ctx Context, handle interface{}) error
}

type handleDB func(*sqlx.DB) error
type handleTX func(*sqlx.Tx) error
type handleQ func(Q) error

type dbGetter func() *sqlx.DB

type Q interface {
	BindNamed(string, interface{}) (string, []interface{}, error)
	Exec(string, ...interface{}) (sql.Result, error)
	ExecContext(Context, string, ...interface{}) (sql.Result, error)
	Get(interface{}, string, ...interface{}) error
	GetContext(Context, interface{}, string, ...interface{}) error
	MustExec(string, ...interface{}) sql.Result
	MustExecContext(Context, string, ...interface{}) sql.Result
	NamedExec(string, interface{}) (sql.Result, error)
	NamedExecContext(Context, string, interface{}) (sql.Result, error)
	NamedQuery(string, interface{}) (*sqlx.Rows, error)
	// NamedQueryContext(Context, string, interface{}) (*sqlx.Rows, error)
	Prepare(string) (*sql.Stmt, error)
	PrepareContext(Context, string) (*sql.Stmt, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	// PrepareNamedContext(Context, string) (*sqlx.NamedStmt, error)
	Preparex(string) (*sqlx.Stmt, error)
	// PreparexContext(Context, string) (*sqlx.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryContext(Context, string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	QueryRowContext(Context, string, ...interface{}) *sql.Row
	QueryRowx(string, ...interface{}) *sqlx.Row
	QueryRowxContext(Context, string, ...interface{}) *sqlx.Row
	Queryx(string, ...interface{}) (*sqlx.Rows, error)
	QueryxContext(Context, string, ...interface{}) (*sqlx.Rows, error)
	Rebind(string) string
	Select(interface{}, string, ...interface{}) error
	SelectContext(Context, interface{}, string, ...interface{}) error
}
