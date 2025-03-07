package repository

import (
	"context"
)

type Creater interface {
	Create(w Worker, ctx context.Context) error
}

type Reader interface {
	Read(w Worker, ctx context.Context) error
}

type Updater interface {
	Update(w Worker, ctx context.Context) (int64, error)
}

type Deleter interface {
	Delete(w Worker, ctx context.Context) (int64, error)
}

type Repository interface {
	Creater
	Reader
	Updater
	Deleter
}

type Entity struct {
	Id int64
}

type NamedEntity struct {
	Entity
	Name string
}
