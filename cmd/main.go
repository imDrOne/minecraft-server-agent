package main

import (
	"flag"
	"github.com/imDrOne/minecraft-server-agent/internal/app"
	"github.com/imDrOne/minecraft-server-agent/internal/pkg/http_server"
	"os"
)

func main() {
	env := flag.String("env", "local", "Application profile")
	flag.Parse()
	application := app.NewApplication(*env)
	config := application.Config()
	logger := application.Logger()

	httpServer := http_server.New(
		http_server.WithPort(config.HTTPServer.Port),
		http_server.WithPrintingRoutes(config.HTTPServer.PrintRoutes),
		http_server.WithAppName(config.Name),
	)

	application.AddServer(httpServer)

	if err := application.Run(); err != nil {
		logger.Error("error during application start", "error", err, "env", *env)
		os.Exit(1)
	}
}
