package migrations

import (
	"embed"

	"github.com/Lysoul/gocommon/monitoring"
	"github.com/uptrace/bun/migrate"
	"go.uber.org/zap"
)

// nolint: gochecknoglobals //later
var Migration = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	log := monitoring.Logger()
	log.Info("Discovering SQL migrations...", zap.Any("sqlMigrations", sqlMigrations))
	if err := Migration.Discover(sqlMigrations); err != nil {
		panic(err)
	} else {
		log.Info("SQL migrations discovered successfully")
	}
}
