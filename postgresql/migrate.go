package postgresql

import (
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/triabokon/gotagv/flags"
)

func MigrationCommand(migrations []*migrate.Migration) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "psql",
		Short: "psqlCmd to hold up and down migrations for postgresql",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flags.MustBindEnvToFlagSet(cmd.Flags())
		},
	}

	var config MigrateConfig
	cobraCmd.PersistentFlags().AddFlagSet(config.Flags())

	upCmd := &cobra.Command{
		Use:   "up",
		Short: "apply Up migrations to the postgresql specified under config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ApplyMigrations(migrate.Up, migrations, config)
		},
	}

	downCmd := &cobra.Command{
		Use:   "down",
		Short: "apply Down migrations to the postgresql specified under config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ApplyMigrations(migrate.Down, migrations, config)
		},
	}

	cobraCmd.AddCommand(upCmd)
	cobraCmd.AddCommand(downCmd)

	return cobraCmd
}

func ApplyMigrations(
	direction migrate.MigrationDirection, migrations []*migrate.Migration, config MigrateConfig,
) error {
	if config.Table == "" {
		return errors.New("empty table to store migrations state")
	}

	migrate.SetTable(config.Table)
	migrate.SetSchema(config.Schema)

	var logger, _ = zap.NewProduction(zap.AddStacktrace(zapcore.InfoLevel))
	source := migrate.MemoryMigrationSource{Migrations: migrations}
	db, closeDB, pgDBErr := NewStdSQL(config.PgDB)
	if pgDBErr != nil {
		return errors.Wrap(pgDBErr, "failed to init pg std client")
	}
	defer func() {
		if pgClErr := closeDB(); pgClErr != nil {
			logger.Error("failed to close pg client", zap.Error(pgClErr))
		}
	}()

	var nApplied int
	var err error
	if direction == migrate.Up {
		nApplied, err = migrate.Exec(db, "postgres", source, direction)
	} else {
		nApplied, err = migrate.Exec(db, "postgres", source, direction)
	}
	if err != nil {
		return errors.Wrap(err, "failed to exec migrations")
	}

	if nApplied > 0 {
		logger.Info(
			"migrations applied",
			zap.Int("migrationsAppliedCount", nApplied),
		)
	} else {
		logger.Info("no migrations applied")
	}

	return nil
}
