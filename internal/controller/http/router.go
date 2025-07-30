package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imDrOne/minecraft-server-agent/internal/app"
	"github.com/imDrOne/minecraft-server-agent/internal/pkg/http_server"
	"net/http"
)

func RegisterRoutes(application *app.Application, server *http_server.HttpServer) {
	server.App.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })
}
