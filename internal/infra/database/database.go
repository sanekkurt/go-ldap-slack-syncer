package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go-ldap-slack-syncer/internal/config/configstructs"
)

const (
	dbDriverName = "mysql"
)

var (
	connTTL = 3 * time.Minute
)

type Database struct {
	dbConn *sqlx.DB
}

type databaseCtxKey struct{}

func InitDatabase(ctx context.Context, cfg configstructs.MySQL) (*Database, error) {

	dbConn, err := sqlx.Connect(dbDriverName, fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true", cfg.Username, cfg.Password, fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), cfg.Database))
	if err != nil {
		return nil, err
	}

	dbConn.SetConnMaxLifetime(connTTL)

	db := &Database{
		dbConn: dbConn,
	}

	_, err = dbConn.Exec("SET autocommit = 0")
	if err != nil {
		return nil, fmt.Errorf("error disable autocommit: %w", err)
	}

	if err := db.migrateDatabase(ctx); err != nil {
		return nil, fmt.Errorf("error performing schema migration: %w", err)
	}

	return db, nil
}

func (db *Database) GetDBConn() *sqlx.DB {
	return db.dbConn
}

func ContextWithDatabase(ctx context.Context, db *Database) context.Context {
	return context.WithValue(ctx, databaseCtxKey{}, db)
}

func GetDatabaseFromContext(ctx context.Context) (*Database, error) {
	db, ok := ctx.Value(databaseCtxKey{}).(*Database)
	if !ok {
		return nil, fmt.Errorf("database not found in context")
	}

	return db, nil
}
