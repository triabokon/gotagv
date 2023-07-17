package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/triabokon/gotagv/internal/auth"
	"github.com/triabokon/gotagv/internal/controller"
	"github.com/triabokon/gotagv/internal/flags"
	"github.com/triabokon/gotagv/internal/postgresql"
	"github.com/triabokon/gotagv/internal/server"
	"github.com/triabokon/gotagv/internal/storage"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "server",
		Aliases:      []string{"s"},
		Short:        "starts a server.",
		SilenceUsage: true,
	}

	var config Config
	cmd.Flags().AddFlagSet(config.Flags())

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		flags.MustBindEnvToFlagSet(cmd.Flags())
		var logger, _ = zap.NewProduction(zap.AddStacktrace(zapcore.InfoLevel))
		pgClient, pgClientCl, err := postgresql.New(cmd.Context(), config.Postgres)
		if err != nil {
			return fmt.Errorf("failed to init postgresql client: %w", err)
		}
		defer func() {
			pgErr := pgClientCl()
			if pgErr != nil {
				logger.Error("failed to close pg client", zap.Error(pgErr))
			}
		}()

		srv := server.New(logger, &config.HTTP, auth.New(&config.Auth), controller.New(storage.New(pgClient)))
		srv.SetRoutes()

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		// Handle SIGINT and SIGTERM signals
		go func() {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			<-signals
			cancel()
		}()

		if sErr := srv.ServeWithGracefulShutdown(ctx, logger); sErr != nil {
			logger.Error("failed to start server", zap.Error(sErr))
			os.Exit(1)
		}
		return nil
	}
	return cmd
}
