package migrations

import (
	"fmt"

	PkgDataSource "github.com/Adhiana46/aegis-asesmen/pkg/data_sources"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresMigrator struct {
	db      *PkgDataSource.PostgresDB
	basedir string
}

func NewPostgresMigrator(db *PkgDataSource.PostgresDB, migrationDir string) DatabaseMigrator {
	return &PostgresMigrator{
		db:      db,
		basedir: migrationDir,
	}
}

func (r *PostgresMigrator) Up() (errOut error) {
	driver, err := migratePostgres.WithInstance(r.db.DB.DB, &migratePostgres.Config{})
	if err != nil {
		errOut = err
		return
	}

	sourceDir := fmt.Sprintf("file://%s", r.basedir)
	m, err := migrate.NewWithDatabaseInstance(
		sourceDir,
		"postgres", driver)
	if err != nil {
		errOut = err
		return
	}

	if err := m.Up(); err != nil {
		errOut = err
		return
	}

	return nil
}

func (r *PostgresMigrator) Down() error {
	driver, err := migratePostgres.WithInstance(r.db.DB.DB, &migratePostgres.Config{})
	if err != nil {
		return err
	}

	sourceDir := fmt.Sprintf("file://%s", r.basedir)
	m, err := migrate.NewWithDatabaseInstance(
		sourceDir,
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil {
		return err
	}

	return nil
}
