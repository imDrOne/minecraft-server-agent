package app

import (
	"context"
	"fmt"
	"github.com/imDrOne/minecraft-server-agent/internal/config"
	"github.com/imDrOne/minecraft-server-agent/internal/pkg/http_server"
	"github.com/imDrOne/minecraft-server-agent/internal/pkg/logger"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server interface {
	Start() error
	Name() string
	Shutdown() error
	Notify() <-chan error
}

var _ Server = (*http_server.HttpServer)(nil)

type Application struct {
	env     string
	config  *config.Config
	logger  *slog.Logger
	servers map[string]Server
}

func (app *Application) Logger() *slog.Logger {
	return app.logger
}

func NewApplication(env string) *Application {
	cfg := config.MustLoad(env)
	l := logger.MustCreateLogger(cfg.Logging.Level, cfg.Logging.Format)
	return &Application{
		env:     env,
		config:  cfg,
		logger:  l,
		servers: make(map[string]Server),
	}
}

func (app *Application) Config() *config.Config {
	return app.config
}

func (app *Application) AddServer(server Server, consumer func(*Application)) {
	if _, ok := app.servers[server.Name()]; !ok {
		app.servers[server.Name()] = server
		consumer(app)
		app.logger.Debug("server registered", "name", server.Name())
	}
	app.logger.Debug("server already registered", "name", server.Name())
}

func (app *Application) GetServer(name string) Server {
	return app.servers[name]
}

func (app *Application) Run() error {
	if len(app.servers) == 0 {
		return fmt.Errorf("no servers registered")
	}

	app.logger.Info("starting all servers", "count", len(app.servers))

	for _, server := range app.servers {
		if err := server.Start(); err != nil {
			return fmt.Errorf("failed to start %s: %w", server.Name(), err)
		}
	}

	return app.waitForShutdown()
}

func (app *Application) waitForShutdown() error {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	serverErrors := make(chan error, len(app.servers))
	for _, server := range app.servers {
		go func() {
			err := <-server.Notify()
			if err != nil {
				app.logger.Error("server error", "name", server.Name(), "error", err)
			} else {
				app.logger.Info("server stopped normally", "name", server.Name())
			}
			serverErrors <- err
		}()
	}

	select {
	case <-ctx.Done():
		app.logger.Info("shutdown signal received")
		stop()

	case err := <-serverErrors:
		if err != nil {
			app.logger.Error("server error triggered shutdown", "error", err)
		} else {
			app.logger.Info("server stopped, triggering application shutdown")
		}
	}

	return app.performShutdown()
}

func (app *Application) performShutdown() error {
	app.logger.Info("performing graceful shutdown", "servers", len(app.servers))

	var wg sync.WaitGroup
	errorsChan := make(chan error, len(app.servers))

	for _, server := range app.servers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			app.logger.Info("shutting down server", "name", server.Name())
			startTime := time.Now()

			if err := server.Shutdown(); err != nil {
				app.logger.Error("server shutdown error",
					"name", server.Name(),
					"error", err,
					"duration", time.Since(startTime))
				errorsChan <- fmt.Errorf("%s: %w", server.Name(), err)
				return
			}

			app.logger.Info("server shutdown completed",
				"name", server.Name(),
				"duration", time.Since(startTime))
		}()
	}

	wg.Wait()
	close(errorsChan)

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		app.logger.Error("shutdown completed with errors", "errors_count", len(errors))
		return errors[0]
	}

	app.logger.Info("graceful shutdown completed successfully")
	return nil
}
