//go:build integration

package repository_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bhmt/tittlemanscrest/repository"
	_ "github.com/lib/pq"
)

func TestIntegrationNewSession(t *testing.T) {
	postgres_db, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		postgres_db = "localhost"
	}

	dns := fmt.Sprintf("postgresql://postgres:postgres@%s:5432/gotest?sslmode=disable", postgres_db)
	session, err := repository.NewSession("postgres", dns)
	if err != nil {
		t.Error(err)
		return
	}

	if err := session.DB.Ping(); err != nil {
		t.Error(err)
	}
}
