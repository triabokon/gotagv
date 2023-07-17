package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/mux"

	"github.com/triabokon/gotagv/internal/auth"
	"github.com/triabokon/gotagv/internal/controller"
	"github.com/triabokon/gotagv/internal/model"
)

const entityIDKey = "id"

type Auth interface {
	CreateToken(userID string) (string, error)
	ValidateToken(tknStr string) (*auth.Claims, error)

	HandleAuth(next http.HandlerFunc) http.HandlerFunc
}

type Controller interface {
	GetUser(ctx context.Context, id string) error
	CreateUser(ctx context.Context, id string) error

	ListVideos(ctx context.Context) ([]*model.Video, error)
	CreateVideo(ctx context.Context, p *controller.CreateVideoParams) (string, error)
	DeleteVideo(ctx context.Context, id string) error

	ListAnnotations(ctx context.Context, p *controller.ListAnnotationsParams) ([]*model.Annotation, error)
	CreateAnnotation(ctx context.Context, p *model.CreateAnnotationParams) (string, error)
	UpdateAnnotation(ctx context.Context, id string, p *model.UpdateAnnotationParams) error
	DeleteAnnotation(ctx context.Context, id string) error
}

type Server struct {
	router *mux.Router
	logger *zap.Logger
	config *Config

	auth       Auth
	controller Controller
}

func New(logger *zap.Logger, config *Config, a Auth, ctrl Controller) *Server {
	srv := &Server{
		router:     mux.NewRouter(),
		logger:     logger,
		config:     config,
		auth:       a,
		controller: ctrl,
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
		err = nil
	}

	// ensure graceful shutdown
	if atomic.LoadInt64(&shutdown) != 1 {
		// if there is no shutdown signal, execute shutdown anyway
		gracefulShutdown()
	}

	<-done

	return err
}

func (s *Server) SuccessResponse(w http.ResponseWriter, result interface{}) {
	body, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, wErr := w.Write(body); wErr != nil {
		s.logger.Error("failed to write response body", zap.Error(wErr))
	}
}

type Response struct {
	Message string `json:"message"`
}

func (s *Server) ErrorResponse(w http.ResponseWriter, err error, code int) {
	body, err := json.MarshalIndent(&Response{Message: err.Error()}, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("JSON marshal failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if _, wErr := w.Write(body); wErr != nil {
		s.logger.Error("failed to write response body", zap.Error(wErr))
	}
}

func parseDuration(d string) (time.Duration, error) {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return 0, err
	}
	return duration, nil
}
