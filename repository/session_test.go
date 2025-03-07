package repository_test

import (
	"testing"

	"github.com/bhmt/tittlemanscrest/repository"
	_ "modernc.org/sqlite"
)

func TestNewSession(t *testing.T) {
	s, err := repository.NewSession("sqlite", ":memory:")
	if err != nil {
		t.Error(err)
	}

	if err := s.DB.Ping(); err != nil {
		t.Error(err)
	}
}
