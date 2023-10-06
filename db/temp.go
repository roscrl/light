package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/roscrl/light/db/sqlc"
)

func NewTempDBFileBenchmarkWithCleanup(b *testing.B) string {
	temp, err := os.CreateTemp("", "app-*")
	if err != nil {
		b.Fatalf("creating temp file: %s", err)
	}

	b.Cleanup(func() {
		os.Remove(temp.Name())
	})

	return temp.Name()
}

func NewTempMigratedDBAndQueriesTestingWithCleanup(t *testing.T) (*sql.DB, *sqlc.Queries) {
	temp, err := os.CreateTemp("", "app-*")
	if err != nil {
		t.Errorf("creating temp file: %s", err)
	}

	tempDB := New(temp.Name())

	t.Cleanup(func() {
		if err := tempDB.Close(); err != nil {
			t.Errorf("closing temp db: %s", err)
		}

		os.Remove(temp.Name())
	})

	RunMigrations(tempDB)
	qry := sqlc.New(tempDB)

	return tempDB, qry
}
