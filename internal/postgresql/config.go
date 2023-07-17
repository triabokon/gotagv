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

	f.AddFlagSet(c.PgDB.Flags("postgresql"))
	f.StringVar(&c.Table, "table", "migrations", "table name where to store the last applied migration id")
	f.StringVar(&c.Schema, "schema", "", "schema name where to look for a `table`")

	return f
}
