package migrations

import "embed"

// DBMigrations content holds database migrations scripts.
//go:embed *.sql
var DBMigrations embed.FS
