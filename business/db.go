package business

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

const migrationsPath = "db/migrations"

var ErrNoMigrationsPending = errors.New("no migrations pending")

func GetDb(dbpath string) (*sqlx.DB, error) {
	err := doMigrations(dbpath)
	if err != nil {
		if errors.Is(err, ErrNoMigrationsPending) {
			log.Println("no migrations to run")
		} else {
			return nil, fmt.Errorf("error running migrations: %w", err)
		}
	}

	db, err := sqlx.Connect("sqlite3", dbpath)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	db.SetMaxOpenConns(1)

	return db, nil
}

func doMigrations(dbpath string) error {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return fmt.Errorf("error opening sqlite database: %w", err)
	}

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("error initializing sqlite database: %w", err)
	}

	fSrc, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		return fmt.Errorf("error opening migrations directory: %w", err)
	}

	m, err := migrate.NewWithInstance(dbpath, fSrc, "sqlite3", instance)
	if err != nil {
		return fmt.Errorf("error initializing migrate: %w", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoMigrationsPending
		}

		return fmt.Errorf("error running migrations up: %w", err)
	}

	return nil
}
