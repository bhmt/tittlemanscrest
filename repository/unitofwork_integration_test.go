//go:build integration

package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bhmt/tittlemanscrest/repository"
	_ "github.com/lib/pq"
)

// Container is a test entity.
// It represents the following postgres table
//
// create table if not exists container (
//
//	id bigserial primary key,
//	name text unique not null,
//	classification text not null default 'goodware'
//
// );
type Container struct {
	repository.NamedEntity
	Classification string
}

func NewContainer(name string) Container {
	entity := Container{Classification: "goodware"}
	entity.Name = name
	return entity
}

func (entity *Container) Create(w repository.Worker, ctx context.Context) error {
	query := `insert into container (name, classification) values ($1 , $2) returning id`
	stmt, err := w.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRow(entity.Name, entity.Classification).Scan(&entity.Id)
}

func (entity *Container) Read(w repository.Worker, ctx context.Context) error {
	query := `select id, name, classification from container where id = $1`
	stmt, err := w.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRow(entity.Id).Scan(&entity.Id, &entity.Name, &entity.Classification)
}

func (entity Container) Update(w repository.Worker, ctx context.Context) (int64, error) {
	query := `update container set (name, classification) = ($1, $2) where id = $3`
	stmt, err := w.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(entity.Name, entity.Classification, entity.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (entity Container) Delete(w repository.Worker, ctx context.Context) (int64, error) {
	query := `delete from container where id = $1`
	stmt, err := w.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(entity.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (entity Container) FilterByName(w repository.Worker, ctx context.Context) ([]Container, error) {
	query := `select id, name, classification from container where name like $1`
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

	entities := []Container{}
	for rows.Next() {
		var e Container
		if err := rows.Scan(&e.Id, &e.Name, &e.Classification); err != nil {
			return entities, err
		}
		entities = append(entities, e)
	}
	if err = rows.Err(); err != nil {
		return entities, err
	}
	return entities, nil
}

// Component is a test entity.
// It represents the following postgres table
//
// create table if not exists component (
//
//	id bigserial primary key,
//	score int not null default 5
//
// );
type Component struct {
	repository.Entity
	Score int64
}

func NewComponent(score int64) Component {
	entity := Component{Score: score}
	return entity
}

func (entity *Component) Create(w repository.Worker, ctx context.Context) error {
	query := `insert into component (score) values ($1) returning id`
	stmt, err := w.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRow(entity.Score).Scan(&entity.Id)
}

func (entity *Component) Read(w repository.Worker, ctx context.Context) error {
	query := `select id, score from component where id = $1`
	stmt, err := w.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRow(entity.Id).Scan(&entity.Id, &entity.Score)
}

func (entity Component) Update(w repository.Worker, ctx context.Context) (int64, error) {
	query := `update component set score = $1 where id = $2`
	stmt, err := w.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(entity.Score, entity.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (entity Component) Delete(w repository.Worker, ctx context.Context) (int64, error) {
	query := `delete from component where id = $1`
	stmt, err := w.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(entity.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Association is a test entity.
// It represents the following postgres table
//
// create table  if not exists association (
//
//	id bigserial primary key,
//	container_id  bigint references container(id) on delete cascade,
//	component_id  bigint references component(id) on delete cascade,
//	constraint uq_mapping unique (container_id, component_id)
//
// );
type Association struct {
	repository.Entity
	ContainerId int64
	ComponentId int64
}

func NewAssociation(containerId, componentId int64) Association {
	return Association{ContainerId: containerId, ComponentId: componentId}
}

func (entity Association) Create(w repository.Worker, ctx context.Context) error {
	query := `insert into association (container_id, component_id) values ($1, $2) on conflict do nothing returning id`

	stmt, err := w.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(entity.ContainerId, entity.ComponentId).Scan(&entity.Id)
}

func (entity Association) Delete(w repository.Worker, ctx context.Context) (int64, error) {
	query := `delete from association where id = $1`
	stmt, err := w.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(entity.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (entity Association) GetContainerComponents(w repository.Worker, ctx context.Context) ([]Component, error) {
	query := `select c.id, c.score from component c join association a on c.id = a.component_id where a.container_id = $1`

	stmt, err := w.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(entity.ContainerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := []Component{}
	for rows.Next() {
		var e Component
		if err := rows.Scan(&e.Id, &e.Score); err != nil {
			return entities, err
		}
		entities = append(entities, e)
	}
	if err = rows.Err(); err != nil {
		return entities, err
	}
	return entities, nil
}

type Action func(s *repository.Session) error

// Integration test is created to test the repository functionality using postgres.
// The test example represents a many-to-many relationship between Container and Component.
// Migration is required to use the test entities.
// The migrations are desined in the migrations directory.
func TestUnitOfWorkIntegration(t *testing.T) {
	postgres_db, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		postgres_db = "localhost"
	}

	dns := fmt.Sprintf("postgresql://postgres:postgres@%s:5432/gotest?sslmode=disable", postgres_db)
	session, err := repository.NewSession("postgres", dns)
	if err != nil {
		t.Fatal(err)
		return
	}

	actions := []Action{
		create,
		update,
		delete,
		find,
	}

	for _, action := range actions {
		if err := action(session); err != nil {
			t.Fatal(err)
		}
	}
}

func create(s *repository.Session) error {
	log.SetPrefix("[create]")
	uow, err := repository.NewUnitOfWork(
		s,
		repository.WithTransaction(s),
	)
	if err != nil {
		return err
	}

	defer uow.Rollback()
	ctx := context.Background()

	container := NewContainer("container")
	if err := container.Create(uow.Worker, ctx); err != nil {
		return err
	}
	log.Println("container created")

	components := []Component{
		NewComponent(10),
		NewComponent(7),
	}

	for _, component := range components {
		if err := component.Create(uow.Worker, ctx); err != nil {
			return err
		}
		log.Println("component created")

		a := NewAssociation(container.Id, component.Id)
		if err := a.Create(uow.Worker, ctx); err != nil {
			return err
		}
		log.Println("association created")
	}

	return uow.Commit()
}

func update(s *repository.Session) error {
	log.SetPrefix("[update]")
	uow, err := repository.NewUnitOfWork(
		s,
		repository.WithTransaction(s),
	)
	if err != nil {
		return err
	}

	defer uow.Rollback()
	ctx := context.Background()

	container := NewContainer("container")
	containers, err := container.FilterByName(uow.Worker, ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no rows found")
			return nil
		}
		return err
	}

	toUpdate := containers[0]

	toUpdate.Classification = "malicious"
	affected, err := toUpdate.Update(uow.Worker, ctx)
	if err != nil {
		return err
	}
	log.Printf("%d rows affectedn\n", affected)

	return uow.Commit()
}

func delete(s *repository.Session) error {
	log.SetPrefix("[delete]")
	uow, err := repository.NewUnitOfWork(
		s,
		repository.WithTransaction(s),
	)
	if err != nil {
		return err
	}

	defer uow.Rollback()
	ctx := context.Background()

	component := NewComponent(7)
	component.Id = 1

	if err := component.Read(uow.Worker, ctx); err != nil {
		return err
	}

	affected, err := component.Delete(uow.Worker, ctx)
	if err != nil {
		return err
	}
	log.Printf("%d rows affectedn\n", affected)

	return uow.Commit()
}

func find(s *repository.Session) error {
	log.SetPrefix("[find]")
	uow, err := repository.NewUnitOfWork(s)
	if err != nil {
		return err
	}

	ctx := context.Background()

	container := NewContainer("container")
	containers, err := container.FilterByName(uow.Worker, ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no rows found")
			return nil
		}
		return err
	}

	a := NewAssociation(containers[0].Id, 0)
	components, err := a.GetContainerComponents(uow.Worker, ctx)
	if err != nil {
		return err
	}

	log.Printf("%d components found", len(components))
	return nil
}
