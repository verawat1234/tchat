package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"tchat.dev/auth/internal/config"
	"tchat.dev/auth/internal/handlers"
)

// Server wraps the HTTP stack for the auth service.
type Server struct {
	cfg        config.Config
	httpServer *http.Server
}

// New constructs a server with the default middleware/handlers.
func New(cfg config.Config) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/ready", handlers.Ready)

	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
		},
	}
}

// Run starts the HTTP server and blocks until shutdown.
func (s *Server) Run() error {
	log.Printf("authsvc listening on %s", s.cfg.HTTPAddr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http listen: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
