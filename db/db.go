package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

var PathMigrations = ""

func init() {
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled
	PathMigrations = filepath.Join(filepath.Dir(filename), "migrations")
}

func New(dataSource string) *sql.DB {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(1)

	err = setPragmas(db)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func RunMigrations(db *sql.DB) {
	migrationsDir, err := os.ReadDir(PathMigrations)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range migrationsDir {
		if strings.HasSuffix(file.Name(), ".sql") {
			pathMigrationSQL := filepath.Join(PathMigrations, file.Name())

			migration, err := os.ReadFile(pathMigrationSQL)
			if err != nil {
				log.Fatalf("error reading migration %s: %s", pathMigrationSQL, err)
			}

			_, err = db.Exec(string(migration))
			if err != nil {
				log.Fatalf("error running migration %s: %s", pathMigrationSQL, err)
			}
		}
	}
}

func setPragmas(db *sql.DB) error {
	_, err := db.Exec("PRAGMA synchronous = OFF")
	if err != nil {
		return fmt.Errorf("setting PRAGMA synchronous: %w", err)
	}

	_, err = db.Exec("PRAGMA cache_size = 50000")
	if err != nil {
		return fmt.Errorf("setting PRAGMA cache_size: %w", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("setting PRAGMA foreign_keys: %w", err)
	}

	_, err = db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		return fmt.Errorf("setting PRAGMA busy_timeout: %w", err)
	}

	_, err = db.Exec("PRAGMA temp_store = MEMORY")
	if err != nil {
		return fmt.Errorf("setting PRAGMA temp_store: %w", err)
	}

	_, err = db.Exec("PRAGMA mmap_size = 300000000")
	if err != nil {
		return fmt.Errorf("setting PRAGMA mmap_size: %w", err)
	}

	return nil
}
