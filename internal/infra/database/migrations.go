//go:build go1.16
// +build go1.16

package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"go-ldap-slack-syncer/migrations"
)

var (
	dbMigrations = migrations.DBMigrations
)

func (db *Database) migrateDatabase(ctx context.Context) error {
	source, err := httpfs.New(http.FS(dbMigrations), ".")
	if err != nil {
		return fmt.Errorf("cannot create migration database source: %w", err)
	}

	driver, err := mysql.WithInstance(db.dbConn.DB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("cannot create migration database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("httpfs", source, "mysql", driver)
	if err != nil {
		return fmt.Errorf("cannot create migration instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
