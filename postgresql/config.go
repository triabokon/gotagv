package postgresql

import (
	"github.com/spf13/pflag"
)

type MigrateConfig struct {
	PgDB Config

	Table  string
	Schema string
}

func (c *MigrateConfig) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("PgDBMigrator", pflag.PanicOnError)

	f.AddFlagSet(c.PgDB.Flags("PgDB", "postgresql"))
	f.StringVar(&c.Table, "table", "migrations", "table name where to store the last applied migration id."+
		" For empty string use the default table name for the migration lib")
	f.StringVar(&c.Schema, "schema", "", "schema name where to look for a `table`. "+
		"For empty string use the default schema for the migration lib")

	return f
}
