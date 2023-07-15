package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/triabokon/gotagv/cmd/server"
	fm "github.com/triabokon/gotagv/internal/migrations"
	"github.com/triabokon/gotagv/postgresql"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:     "gotagv",
		Aliases: []string{"gtv"},
		Short:   "gotagv is a simple video and annotations management service.",
	}

	var logger, _ = zap.NewProduction(zap.AddStacktrace(zapcore.InfoLevel))

	migrations, err := fm.FindMigrations()
	if err != nil {
		return err
	}

	rootCmd.AddCommand(server.Cmd())
	rootCmd.AddCommand(postgresql.MigrationCommand(migrations))

	if err := rootCmd.Execute(); err != nil {
		logger.Error("failed to execute root cmd", zap.Error(err))
		return err
	}

	return nil
}
