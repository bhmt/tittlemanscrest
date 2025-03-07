package repository

import (
	"context"
	"database/sql"
	"fmt"
)

var ErrNoTransaction = fmt.Errorf("no transaction defined")

// Worker interface abstracts the shared method between [*]sql.DB and [*]sql.Tx.
// These are intended to be used in a unit of work
type Worker interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type UnitOfWork struct {
	tx     *sql.Tx
	db     *sql.DB
	Worker Worker
}

func NewUnitOfWork(session *Session, opts ...func(*UnitOfWork) error) (*UnitOfWork, error) {
	uow := UnitOfWork{db: session.DB}
	for _, o := range opts {
		err := o(&uow)
		if err != nil {
			return nil, err
		}
	}

	if uow.Worker == nil {
		uow.Worker = uow.db
	}

	return &uow, nil
}

func WithTransaction(session *Session) func(*UnitOfWork) error {
	return func(uow *UnitOfWork) error {
		tx, err := session.Begin()
		if err != nil {
			return err
		}

		uow.tx = tx
		uow.Worker = tx
		return nil
	}
}

func (uow UnitOfWork) GetTx() *sql.Tx {
	return uow.tx
}

func (uow UnitOfWork) GetDb() *sql.DB {
	return uow.db
}

func (uow UnitOfWork) Commit() error {
	if uow.tx == nil {
		return ErrNoTransaction
	}

	return uow.tx.Commit()
}

func (uow UnitOfWork) Rollback() error {
	if uow.tx == nil {
		return ErrNoTransaction
	}
	return uow.tx.Rollback()
}

// Stmt is added for completness.
// The method is wrapped to ensure a transaction is defined.
func (uow UnitOfWork) Stmt(stmt *sql.Stmt) (*sql.Stmt, error) {
	if uow.tx == nil {
		return nil, ErrNoTransaction
	}
	return uow.tx.Stmt(stmt), nil
}

// StmtContext is added for completness.
// The method is wrapped to ensure a transaction is defined.
func (uow UnitOfWork) StmtContext(ctx context.Context, stmt *sql.Stmt) (*sql.Stmt, error) {
	if uow.tx == nil {
		return nil, ErrNoTransaction
	}
	return uow.tx.StmtContext(ctx, stmt), nil
}
