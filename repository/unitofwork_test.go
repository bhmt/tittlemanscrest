package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bhmt/tittlemanscrest/repository"
	_ "modernc.org/sqlite"
)

var s, _ = repository.NewSession("sqlite", "file::memory:?cache=shared")

var migration = `
create table if not exists testcomponent (
    id integer primary key autoincrement not null,
    name text not null
)
`

var ErrTestComponent = fmt.Errorf("test error")

type TestComponent struct {
	repository.NamedEntity
}

func NewTestComponent(name string) TestComponent {
	entity := TestComponent{}
	entity.Name = name
	return entity
}

func (entity *TestComponent) Create(w repository.Worker, ctx context.Context) error {
	query := `insert into testcomponent (name) values (?)`
	stmt, err := w.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(entity.Name)
	return err
}

func (entity *TestComponent) CreateError(w repository.Worker, ctx context.Context) error {
	return ErrTestComponent
}

func (entity *TestComponent) FilterByName(w repository.Worker, ctx context.Context) ([]TestComponent, error) {
	query := `select id, name from testcomponent where name like ?`
	stmt, err := w.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(entity.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := []TestComponent{}
	for rows.Next() {
		var e TestComponent
		if err := rows.Scan(&e.Id, &e.Name); err != nil {
			return entities, err
		}
		entities = append(entities, e)
	}
	if err = rows.Err(); err != nil {
		return entities, err
	}
	return entities, nil
}

func TestNewUnitOfWork(t *testing.T) {
	uow, err := repository.NewUnitOfWork(s)
	if err != nil {
		t.Error(err)
		return
	}

	if uow.GetDb() == nil {
		t.Error("database should be defined")
		return
	}

	if uow.GetTx() != nil {
		t.Error("transaction should not be defined")
		return
	}

	uowTx, err := repository.NewUnitOfWork(
		s,
		repository.WithTransaction(s),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if uow.GetDb() == nil {
		t.Error("database should be defined")
		return
	}

	if uowTx.GetTx() == nil {
		t.Error("transaction should be defined")
	}
}

func TestUnitOfWorkTransactionCommit(t *testing.T) {
	uow, _ := repository.NewUnitOfWork(s, repository.WithTransaction(s))
	s.DB.Exec(migration)

	tc := NewTestComponent("TestUnitOfWorkTransactionCommit")
	ctx := context.Background()

	func() {
		if err := tc.Create(uow.Worker, ctx); err != nil {
			t.Error(err)
			uow.Rollback()
			return
		}
		uow.Commit()
	}()

	uow, _ = repository.NewUnitOfWork(s)
	tc = NewTestComponent("%Commit%")

	obtained, err := tc.FilterByName(uow.Worker, ctx)
	if err != nil {
		t.Error(err)
		return
	}

	if len(obtained) != 1 {
		t.Errorf("commit testcomponents want=%v, got=%v", 1, len(obtained))
	}
}

func TestUnitOfWorkTransactionRollback(t *testing.T) {
	uow, _ := repository.NewUnitOfWork(s, repository.WithTransaction(s))
	s.DB.Exec(migration)

	tc := NewTestComponent("TestUnitOfWorkTransactionRollback")
	ctx := context.Background()

	func() {
		if err := tc.CreateError(uow.Worker, ctx); err != nil {
			uow.Rollback()
			return
		}
		uow.Commit()
	}()

	uow, _ = repository.NewUnitOfWork(s)
	tc = NewTestComponent("%Rollback%")

	obtained, err := tc.FilterByName(uow.Worker, ctx)
	if err != nil {
		t.Error(err)
		return
	}

	if len(obtained) != 0 {
		t.Errorf("rollback testcomponents want=%v, got=%v", 0, len(obtained))
	}
}

func TestUnitOfWorkErrNoTransaction(t *testing.T) {
	uow, _ := repository.NewUnitOfWork(s)
	ctx := context.Background()

	errs := []error{
		uow.Commit(),
		uow.Rollback(),
	}

	stmt, _ := uow.GetDb().Prepare("SELECT 1")
	_, err := uow.Stmt(stmt)
	errs = append(errs, err)

	_, err = uow.StmtContext(ctx, stmt)
	errs = append(errs, err)

	for _, e := range errs {
		if e != repository.ErrNoTransaction {
			t.Error(e)
		}
	}
}

func TestUnitOfWorkStmt(t *testing.T) {
	uow, _ := repository.NewUnitOfWork(s, repository.WithTransaction(s))
	ctx := context.Background()

	errs := []error{}

	stmt, _ := uow.GetDb().Prepare("SELECT 1")
	stmt, err := uow.Stmt(stmt)
	errs = append(errs, err)

	_, err = uow.StmtContext(ctx, stmt)
	errs = append(errs, err)

	for _, e := range errs {
		if e != nil {
			t.Error(e)
		}
	}
}
