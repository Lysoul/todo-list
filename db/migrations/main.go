package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

// nolint: gochecknoglobals //later
var Migration = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	if err := Migration.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
