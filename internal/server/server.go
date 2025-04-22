package server

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/middleware"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	mux        *mux.Router
	cfg        *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
		mux: mux.NewRouter(),
		httpServer: &http.Server{
			Addr:           cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderBytes,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) SetupRoutes(routeConfig func(*mux.Router)) {
	apiRouter := s.mux.PathPrefix("/api").Subrouter()
	routeConfig(apiRouter)

	apiRouter.Use(middleware.RecoveryMiddleware)
	s.httpServer.Handler = apiRouter
}

func (s *Server) GracefulStop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
