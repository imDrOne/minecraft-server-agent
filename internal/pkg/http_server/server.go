package http_server

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type HttpServer struct {
	*fiber.App
	name                 string
	notify               chan error
	address              string
	shutdownTimeout      time.Duration
	readTimeout          time.Duration
	writeTimeout         time.Duration
	enablePrintingRoutes bool
}

func New(opts ...HttpOption) *HttpServer {

	server := &HttpServer{
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
	}

	for _, opt := range opts {
		opt(server)
	}

	server.App = fiber.New(fiber.Config{
		ReadTimeout:       server.readTimeout,
		WriteTimeout:      server.writeTimeout,
		AppName:           server.name,
		EnablePrintRoutes: server.enablePrintingRoutes,
	})
	return server
}

func (s *HttpServer) Start() error {
	go func() {
		slog.Info("starting HTTP server", "address", s.address, "name", s.name)
		s.notify <- s.App.Listen(s.address)
		close(s.notify)
	}()
	return nil
}

func (s *HttpServer) Notify() <-chan error {
	return s.notify
}

func (s *HttpServer) Shutdown() error {
	slog.Info("shutting down HTTP server", "name", s.name, "timeout", s.shutdownTimeout)
	return s.App.ShutdownWithTimeout(s.shutdownTimeout)
}

func (s *HttpServer) Name() string {
	tmp := s.name
	return tmp
}
