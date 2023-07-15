package postgresql

import (
	"embed"

	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed *.sql
var fs embed.FS

func FindMigrations() ([]*migrate.Migration, error) {
	var err error
	migrations, err := migrate.EmbedFileSystemMigrationSource{FileSystem: fs, Root: "."}.FindMigrations()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find migrations")
	}

	return migrations, nil
}
