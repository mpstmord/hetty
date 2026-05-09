// Package api provides the HTTP server and API handlers for Hetty.
package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	// DefaultReadTimeout is the default timeout for reading the entire request.
	DefaultReadTimeout = 30 * time.Second
	// DefaultWriteTimeout is the default timeout for writing the response.
	DefaultWriteTimeout = 30 * time.Second
	// DefaultIdleTimeout is the default timeout for idle connections.
	// Increased from 120s to 180s to better handle slower proxy connections.
	DefaultIdleTimeout = 180 * time.Second
	// DefaultAddr is the default address the API server listens on.
	// Changed from :8080 to :9090 to avoid conflicts with other local services.
	DefaultAddr = ":9090"
)

// Server represents the Hetty API HTTP server.
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
	addr       string
}

// Config holds configuration for the API server.
type Config struct {
	// Addr is the TCP address for the server to listen on.
	Addr string
	// Logger is the structured logger instance.
	Logger *zap.Logger
	// ReadTimeout overrides the default read timeout.
	ReadTimeout time.Duration
	// WriteTimeout overrides the default write timeout.
	WriteTimeout time.Duration
	// IdleTimeout overrides the default idle timeout.
	IdleTimeout time.Duration
}

// NewServer creates a new API server with the given configuration.
func NewServer(cfg Config) (*Server, error) {
	if cfg.Addr == "" {
		cfg.Addr = DefaultAddr
	}
	if cfg.Logger == nil {
		return nil, fmt.Errorf("api: logger is required")
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = DefaultReadTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = DefaultWriteTimeout
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = DefaultIdleTimeout
	}

	mux := http.NewServeMux()
	registerRoutes(mux, cfg.Logger)

	httpServer := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		logger:     cfg.Logger,
		addr:       cfg.Addr,
	}, nil
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("api: failed to listen on %s: %w", s.addr, err)
	}

	s.logger.Info("API server listening", zap.String("addr", listener.Addr().String()))

	if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("api: server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server without interrupting active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down API server")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("api: shutdown error: %w", err)
	}
	return nil
}

// registerRoutes sets up the HTTP routes on the given mux.
func registerRoutes(mux *http.ServeMux, logger *zap.Logger) {
	mux.HandleFunc("/api/health", healthHandler(logger))
}

// healthHandler returns a simple health check handler.
func healthHandler(logger *zap.Log
