package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	router *http.ServeMux
	logger *zap.Logger
	config *Config
}

func New(config *Config, logger *zap.Logger) *Server {
	srv := &Server{
		router: http.NewServeMux(),
		logger: logger,
		config: config,
	}
	return srv
}

func (s *Server) newHTTPSrv() *http.Server {
	return &http.Server{
		Handler:      s.router,
		Addr:         s.config.Bind,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}
}

func (s *Server) ServeWithGracefulShutdown(ctx context.Context, logger *zap.Logger) error {
	srv := s.newHTTPSrv()
	// start listening to SIGINT and SIGTERM syscalls
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		close(c)
	}()

	shutdown := int64(0)
	done := make(chan struct{})
	gracefulShutdown := func() {
		defer close(done)
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if atomic.CompareAndSwapInt64(&shutdown, 0, 1) {
			if err := srv.Shutdown(tctx); err != nil {
				logger.Error("srv shutdown failed", zap.Error(err))
			}
		}
	}

	go func() {
		select {
		case <-ctx.Done():
			logger.Info(
				"context done received",
				zap.Error(ctx.Err()),
			)
		case sig := <-c:
			logger.Info(
				"signal received",
				zap.String("signal", sig.String()),
			)
		}

		gracefulShutdown()
	}()

	logger.Info(
		"service started",
		zap.Int("pid", syscall.Getpid()),
		zap.String("bind", srv.Addr),
	)

	err := srv.ListenAndServe()
	if atomic.LoadInt64(&shutdown) > 0 && err != nil {
		logger.Info(
			"http server ListenAndServe failed",
			zap.Error(err),
		)
		// hiding bugs and races in server shutdown proc
		err = nil
	}

	// Ensure graceful shutdown
	if atomic.LoadInt64(&shutdown) != 1 {
		// if there is no shutdown signal, execute shutdown anyway
		gracefulShutdown()
	}

	<-done

	return err
}

func (s *Server) JSONResponse(w http.ResponseWriter, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON(body))
}

func (s *Server) ErrorResponse(w http.ResponseWriter, error string, code int) {
	data := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: error,
	}

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON(body))
}

func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	return out.Bytes()
}
