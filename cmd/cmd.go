package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/triabokon/gotagv/server"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "gotagv",
		Aliases:      []string{"gvt"},
		Short:        "gotagv is a simple video and annotations management service.",
		SilenceUsage: true,
	}

	var config Config
	cmd.Flags().AddFlagSet(config.Flags())

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {

		var logger, _ = zap.NewProduction(zap.AddStacktrace(zapcore.FatalLevel))

		srv := server.New(&config.HTTP, logger)

		srv.SetRoutes()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle SIGINT and SIGTERM signals
		go func() {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			<-signals
			cancel()
		}()

		if sErr := srv.ServeWithGracefulShutdown(ctx, logger); sErr != nil {
			logger.Error("Failed to start server", zap.Error(sErr))
			os.Exit(1)
		}
		return nil
	}
	return cmd
}

func Execute() {
	if err := Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
