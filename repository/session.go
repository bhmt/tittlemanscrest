package repository

import (
	"database/sql"
)

type Session struct {
	*sql.DB
}

func NewSession(driverName string, dns string) (*Session, error) {
	db, err := sql.Open(driverName, dns)
	if err != nil {
		return nil, err
	}

	s := Session{DB: db}
	return &s, nil
}
